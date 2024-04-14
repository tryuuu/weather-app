# ステージ1: ビルド環境
FROM golang:1.18 as builder

# 作業ディレクトリを設定
WORKDIR /app

# モジュールファイルをコピー
COPY go.mod .
COPY go.sum .

# モジュールをダウンロード
RUN go mod download

# ソースコードをコピー
COPY . .

# アプリケーションをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o weather-app .

# ステージ2: 実行環境
FROM alpine:latest  

# アルパインには標準でbashが含まれていないため、シェルとしてashが利用される
WORKDIR /root/

# ビルダーステージから実行ファイルをコピー
COPY --from=builder /app/weather-app .

# ポート8080を公開
EXPOSE 8080

# 実行コマンド
CMD ["./weather-app"]