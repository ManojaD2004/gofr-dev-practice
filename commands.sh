# Redis
docker run --name gofr-redis -p 6379:6379 -d redis     

docker exec -it gofr-redis bash -c 'redis-cli SET greeting "Hello from Redis."'

# MySQL
docker run --name gofr-mysql -e MYSQL_ROOT_PASSWORD=root123 -e MYSQL_DATABASE=test_db -p 3306:3306 -d mysql:8.0.30

docker exec -it gofr-mysql mysql -uroot -proot123 test_db -e "CREATE TABLE customers (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255) NOT NULL);"
