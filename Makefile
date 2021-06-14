DOCKER_COMPOSE := docker compose -p gqlgentodos

help:             ## このヘルプを表示します
	@awk -F: '/^[A-Za-z0-9_-]+:.*## / { sub(/.*## /, "", $$2); printf "make %-15s %s\n", $$1, $$2 }' Makefile

up:               ## dockerコンテナの作成と起動
	$(DOCKER_COMPOSE) up -d

down:             ## dockerコンテナの停止と削除
	$(DOCKER_COMPOSE) down

ps:               ## dockerコンテナの一覧を確認
	$(DOCKER_COMPOSE) ps

generate:         ## GraphQLのコードを自動生成
	go generate ./...

build:            ## サーバーをビルドする
	go build -o server

run: build        ## サーバーを起動する
	go run server.go
