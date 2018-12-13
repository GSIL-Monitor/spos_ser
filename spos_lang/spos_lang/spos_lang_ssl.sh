#!/usr/bin/env bash	 

echo "---"
echo "kubectl delete spos-lang service..."
kubectl delete svc spos-lang
#kubectl delete svc greeter-server

echo "---"
echo "kubectl delete spos-lang pod..."
#kubectl delete po --force --grace-period=0 -l app=spos-auth 
kubectl delete po -l app=spos-lang
#kubectl delete po -l app=greeter-server

echo "---"
echo "kubectl delete spos-lang deploy..."
kubectl delete deploy spos-lang
#kubectl delete deploy greeter-server

echo "---"
echo "build spos_lang..."

CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o spos-lang-ssl .

echo "---"
echo "delete spos-lang docker image..."
docker rmi  --force patlihh/spos-lang:v1.0

echo "---"
echo "build spos-lang docker image..."
docker build -t patlihh/spos-lang:v1.0 .

echo "---"
echo "push spos-lang docker image..."
docker push patlihh/spos-lang:v1.0

#echo "docker run greeter-client..."
#docker run -it -e WORDPRESS_DB_HOST=192.168.99.100 -e WORDPRESS_DB_PORT=31135 -e WORDPRESS_DB_PASSWORD=Unitone -e #SPOS_AUTH=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MzU5NDcyNTUsImlhdCI6MTUzNTk0MzY1NX0.C5axQvnEdpqiE3P2ADr4gj3i3Mvc1zLGaR2TPSelYa8 -e ADDRESS=192.168.99.100:30372 greeter-client:v1
echo "---"
echo "kubectl deploy spos-lang..."
kubectl apply -f ./spos-lang.yaml
