FROM golang:1.23-alpine

WORKDIR /app

# 必要なパッケージをインストール
RUN apk add --no-cache git

# 依存関係をコピーとダウンロード
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# ビルドキャッシュの実行
RUN go build -o /dev/null ./...

# ホットリロード用のAir（オプション）
RUN go install github.com/cosmtrek/air@latest

# 8080ポートを公開
EXPOSE 8080

# デフォルトのコマンド（ホットリロードがある場合はair）
CMD ["go", "run", "cmd/api/main.go"] 