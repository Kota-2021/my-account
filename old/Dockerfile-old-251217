# --------------------
# ステージ 1: ビルドステージ (軽量化のために最終イメージにGoツールを含めない)
# --------------------
  FROM golang:1.25-alpine AS builder

  # 依存関係のキャッシュとビルドを分離
  WORKDIR /app
  COPY go.mod go.sum ./
  RUN go mod download
  
  # ソースコードをコンテナにコピー
  COPY . .
  
  # アプリケーションをビルド
  # CGO_ENABLED=0 は静的リンクを行い、実行環境の依存性をなくす（alpineベースイメージと相性が良い）
  RUN CGO_ENABLED=0 go build -o /main
  
  # --------------------
  # ステージ 2: 実行ステージ (最小限の実行環境)
  # --------------------
  FROM alpine:latest
  
  # ビルドステージで作成したバイナリをコピー
  COPY --from=builder /main /main
    
  # アプリケーションの実行
  CMD ["/main"]