FROM node:18-alpine

WORKDIR /app

# パッケージインストール前の依存関係ファイルをコピー
COPY package.json package-lock.json* ./

# 依存関係のインストール
RUN npm ci

# 開発用サーバー起動
CMD ["npm", "run", "dev"]

# 3000ポートを公開
EXPOSE 3000 