# Stage 1: Build Vue
FROM node:22-alpine AS builder


WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

# Stage 2: Serve with Nginx
FROM nginx:alpine

# 設定工作目錄
WORKDIR /usr/share/nginx/html

# 1. 清空 Nginx 預設的 html
RUN rm -rf ./*

# 2. 把 Vue 編譯好的檔案 (dist) 複製進來
COPY --from=builder /app/dist .

# 3. 複製 Nginx 設定檔 (注意路徑)
# 假設你的 nginx.conf 現在放在 frontend 資料夾下
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
