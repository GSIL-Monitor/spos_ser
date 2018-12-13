#!/usr/bin/env bash	 

echo "---"
echo "build greeter-client..."

CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o greeter-client .

echo "---"
echo "delete greeter-client docker image..."
docker rmi --force greeter-client:v1

echo "---"
echo "build greeter-client docker image..."
docker build -t greeter-client:v1 .

echo "---"
echo "docker run greeter-client..."
#docker run -it -e WORDPRESS_DB_HOST=192.168.99.100 -e WORDPRESS_DB_PORT=31135 -e WORDPRESS_DB_PASSWORD=Unitone -e SPOS_AUTH=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyaWQiOiJVMTU0MjE3NjkyNiIsImV4cCI6MTU0MjE4ODk2MSwiaWF0IjoxNTQyMTg1MzYxfQ.BigMH6VCnvU09KsseUoqgpSrZGOSEXQvlehSbBxTrUI -e ADDRESS=192.168.99.100:30372 greeter-client:v1
#docker run -it -e ADDRESS=192.168.20.104:32001 greeter-client:v1
docker run -it -e ADDRESS=192.168.99.100:32001 greeter-client:v1

