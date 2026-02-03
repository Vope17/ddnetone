all:
	docker compose up -d --build

# --- 變數設定 (方便維護) ---
GCP_REGION := asia-east1
PROJECT_ID := optical-net-485503-g6
REPO_NAME  := ddnetone
APP_NAME   := ddnetone-app
BASE_URL   := $(GCP_REGION)-docker.pkg.dev/$(PROJECT_ID)/$(REPO_NAME)/$(APP_NAME)

# --- 版本控制邏輯 ---
# 預設使用 Git 的短 Hash (例如: a1b2c3d)
# 如果你有手動輸入 VERSION=v1.0，則使用你輸入的
GIT_HASH := $(shell git rev-parse --short HEAD)
VERSION  ?= $(GIT_HASH)

# --- Frontend ---
push_f:
	@echo "正在構建並推送 Frontend 版本: $(VERSION)"
	docker compose build
	
	# 1. 標記特定版本 (如 v1 或 git-hash)
	docker tag ddnetone-frontend $(BASE_URL)/frontend:$(VERSION)
	# 2. 標記 latest (保持最新)
	docker tag ddnetone-frontend $(BASE_URL)/frontend:latest
	
	# 3. 推送兩者
	docker push $(BASE_URL)/frontend:$(VERSION)
	docker push $(BASE_URL)/frontend:latest
	@echo "Frontend 推送完成！版本: $(VERSION)"

# --- Backend ---
push_b:
	@echo "正在構建並推送 Backend 版本: $(VERSION)"
	docker compose build
	
	# 1. 標記特定版本
	docker tag ddnetone-backend $(BASE_URL)/backend:$(VERSION)
	# 2. 標記 latest
	docker tag ddnetone-backend $(BASE_URL)/backend:latest
	
	# 3. 推送兩者
	docker push $(BASE_URL)/backend:$(VERSION)
	docker push $(BASE_URL)/backend:latest
	@echo "Backend 推送完成！版本: $(VERSION)"

# --- 一次推全部 (選用) ---
push_all: push_f push_b
