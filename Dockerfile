FROM debian:stable-slim

WORKDIR /app

# 在目前的資料夾 go build，然後把執行檔複製到 container 裡
COPY blog-aggregator /bin/blog-aggregator

# 因為連線到 postgres container 在 docker-compose.yml 裡是叫 postgres
# 連線的地址會是 @postgres:5432 跟本機開發的 @localhost:5432 不一樣
# 但是 godotenv.Load() 要找 .env 檔，所以才在本地建一個 .env.docker 
# 然後複製到 container 裡還是叫 .env
COPY .env.docker /app/.env

# 執行程式
CMD ["/bin/blog-aggregator"]