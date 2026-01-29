all:
	docker compose up -d --build

push:
	# 1. 重新編譯 (加上 --build 確保它不會用舊的快取)
	docker compose build
	# 2. 貼上標籤 (建議把 v1 改成 v2，或是用 latest)
	docker tag ddnetone-app:latest asia-east1-docker.pkg.dev/optical-net-485503-g6/ddnetone/ddnetone-app:latest
	# 3. 推送到 GCP
	docker push asia-east1-docker.pkg.dev/optical-net-485503-g6/ddnetone/ddnetone-app:latest
