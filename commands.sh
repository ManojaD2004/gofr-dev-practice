# Redis
docker run --name gofr-redis -p 6379:6379 -d redis     

docker exec -it gofr-redis bash -c 'redis-cli SET greeting "Hello from Redis."'

# MySQL
docker run --name gofr-mysql -e MYSQL_ROOT_PASSWORD=root123 -e MYSQL_DATABASE=test_db -p 3306:3306 -d mysql:8.0.30

docker exec -it gofr-mysql mysql -uroot -proot123 test_db -e "CREATE TABLE customers (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255) NOT NULL);"
# Tracing
docker run --name gofr-zipkin -p 2005:9411 -d openzipkin/zipkin:latest
docker run -d --name jaeger -e COLLECTOR_OTLP_ENABLED=true \
  -p 16686:16686 \
  -p 14317:4317 \
  -p 14318:4318 \
  jaegertracing/all-in-one:1.41
