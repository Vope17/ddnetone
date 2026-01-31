package utils

import (
	"strings"
)

// ParseRunnerNames 處理分隔符號：支援 "A, B" 或 "A & B" 或 "A,B"

func ParseRunnerNames(raw string) []string {
	if raw == "" {
		return []string{}
	}

	// 先把 '&' 替換成 ','，再統一用 ',' 切割
	normalized := strings.ReplaceAll(raw, "&", ",")
	names := strings.Split(normalized, ",")

	var result []string
	for _, name := range names {
		trimmed := strings.TrimSpace(name)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
