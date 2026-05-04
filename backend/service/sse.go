package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"DDNETONE/db"
	"DDNETONE/model"
	"github.com/gin-gonic/gin"
)

// SSE Hub — manages all connected clients

type sseClient struct {
	ch chan []byte
}

var (
	sseClients   = make(map[*sseClient]struct{})
	sseMu        sync.RWMutex
	sseHeartbeat = 25 * time.Second
)

// collectAllData gathers the same data that the frontend used to poll.
func collectAllData() (map[string]interface{}, error) {
	database := db.GetDB()

	// summary
	UpdateGlobalSummary()
	var summary model.Summary
	database.Last(&summary)

	// leaderboard (reuse logic from player.go)
	leaderboard := buildLeaderboard()

	// maps
	var maps []model.MapRecord
	database.Order("score desc").Find(&maps)

	// growth (last 7 days)
	sevenDaysAgo := time.Now().AddDate(0, 0, -7).Format(time.RFC3339)
	var growth []model.GrowthData
	database.Where("timestamp >= ?", sevenDaysAgo).Order("id asc").Find(&growth)

	// milestones
	milestones := buildMilestones()

	// score milestones
	scoreMilestones := buildScoreMilestones()

	return map[string]interface{}{
		"summary":          summary,
		"leaderboard":      leaderboard,
		"maps":             maps,
		"growth":           growth,
		"milestones":       milestones,
		"score_milestones": scoreMilestones,
	}, nil
}

// BroadcastUpdate collects fresh data and pushes it to every connected SSE client.
func BroadcastUpdate() {
	sseMu.RLock()
	count := len(sseClients)
	sseMu.RUnlock()

	if count == 0 {
		return
	}

	data, err := collectAllData()
	if err != nil {
		log.Println("SSE: failed to collect data:", err)
		return
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Println("SSE: failed to marshal data:", err)
		return
	}

	msg := formatSSE("update", payload)

	sseMu.RLock()
	defer sseMu.RUnlock()
	for c := range sseClients {
		select {
		case c.ch <- msg:
		default:
			// client too slow, skip
		}
	}
}

func formatSSE(event string, data []byte) []byte {
	return []byte(fmt.Sprintf("event: %s\ndata: %s\n\n", event, data))
}

// HandleSSE is the Gin handler for GET /api/sse
func HandleSSE(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no") // nginx

	// Flush headers immediately
	c.Writer.Flush()

	client := &sseClient{ch: make(chan []byte, 8)}
	sseMu.Lock()
	sseClients[client] = struct{}{}
	sseMu.Unlock()

	defer func() {
		sseMu.Lock()
		delete(sseClients, client)
		sseMu.Unlock()
		close(client.ch)
	}()

	// Send initial full payload on connect
	data, err := collectAllData()
	if err == nil {
		payload, _ := json.Marshal(data)
		fmt.Fprintf(c.Writer, "event: update\ndata: %s\n\n", payload)
		c.Writer.Flush()
	}

	// Stream loop
	ctx := c.Request.Context()
	heartbeatTicker := time.NewTicker(sseHeartbeat)
	defer heartbeatTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-client.ch:
			if !ok {
				return
			}
			_, err := c.Writer.Write(msg)
			if err != nil {
				return
			}
			c.Writer.Flush()
		case <-heartbeatTicker.C:
			_, err := io.WriteString(c.Writer, ": heartbeat\n\n")
			if err != nil {
				return
			}
			c.Writer.Flush()
		}
	}
}
