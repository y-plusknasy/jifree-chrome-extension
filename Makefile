include .env
export

.PHONY: deploy
deploy:
	gcloud functions deploy $(FUNCTION_NAME) \
		--project $(PROJECT_ID) \
		--region $(REGION) \
		--runtime $(RUNTIME) \
		--source ./backend \
		--entry-point $(FUNCTION_NAME) \
		--trigger-http \
		--allow-unauthenticated \
		--set-env-vars SHARED_SECRET=$(SHARED_SECRET),ALLOWED_ORIGIN=$(ALLOWED_ORIGIN),RATE_LIMIT_WINDOW_SEC=$(RATE_LIMIT_WINDOW_SEC)

.PHONY: run-backend
run-backend:
	cd backend && go run cmd/main.go

# フロントエンド開発用（コンテナ内で実行）
.PHONY: dev-extension
dev-extension:
	docker compose exec extension sh -c "npm install && npm run dev"

# フロントエンド本番ビルド用（コンテナ内で実行）
.PHONY: build-extension
build-extension:
	docker compose exec extension sh -c "npm install && npm run build"

# 拡張機能のパッケージング（Chromeウェブストア提出用）
.PHONY: package-extension
package-extension: build-extension
	rm -rf extension/package extension/jifree.zip
	mkdir -p extension/package
	cp extension/manifest.json extension/package/
	cp -r extension/images extension/package/
	cp -r extension/dist extension/package/
	cd extension/package && zip -r ../jifree.zip .
	rm -rf extension/package
	@echo "Package created at extension/jifree.zip"
