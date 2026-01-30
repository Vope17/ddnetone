# Stage 1: Build
FROM golang:1.25-alpine AS builder

WORKDIR /app

# 安裝 git (如果依賴庫需要)
RUN apk add --no-cache git

# 利用 Docker Cache 機制，先複製 mod 檔
COPY go.mod go.sum ./
RUN go mod download


# 複製程式碼並編譯
COPY . .

# -o server 指定輸出檔名
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# Stage 2: Runtime

FROM alpine:latest
WORKDIR /app


RUN apk add --no-cache ca-certificates tzdata

# 從 builder 階段複製編譯好的執行檔
COPY --from=builder /app/server .

# 後端在容器內監聽 8080
EXPOSE 8080

CMD ["./server"]
