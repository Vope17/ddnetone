# --- Stage 1: 建置前端 (Vue) ---
# Node 版本維持剛剛修好的 22
FROM node:22-alpine AS frontend-builder

WORKDIR /app/frontend

COPY frontend/package*.json ./

RUN npm install


COPY frontend/ .
RUN npm run build


# --- Stage 2: 建置後端 (Go) ---
# ★★★ 修改這裡：把 1.21 改成 1.25 (或最新版) ★★★
FROM golang:1.25-alpine AS backend-builder
WORKDIR /app/backend


# 安裝必要套件
RUN apk add --no-cache git

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o server .


# --- Stage 3: 最終執行環境 (Runtime) ---
FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

# 複製編譯好的檔案
COPY --from=backend-builder /app/backend/server .
COPY --from=frontend-builder /app/frontend/dist ./dist

EXPOSE 8080

CMD ["./server"]
