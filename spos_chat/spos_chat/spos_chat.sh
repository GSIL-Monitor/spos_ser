#!/usr/bin/env bash	 

echo "---"
echo "kubectl delete spos-chat service deploy pod..."
kubectl delete -f ./spos-chat.yaml
#kubectl delete svc spos-chat

#echo "kubectl delete spos-chat pod..."
#kubectl delete po --force --grace-period=0 -l app=spos-auth 
#kubectl delete po -l app=spos-chat

#echo "kubectl delete spos-chat deploy..."
#kubectl delete deploy spos-chat

echo "---"
echo "build spos_chat to bin..."
CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o spos-chat .

echo "---"
echo "delete spos-chat docker image..."
docker rmi --force patlihh/spos-chat:v1.0

echo "---"
echo "build spos-chat docker image..."
docker build -t patlihh/spos-chat:v1.0 .

echo "---"
echo "push spos-chat docker image..."
docker push patlihh/spos-chat:v1.0

#echo "docker run greeter-client..."
#docker run -it -e WORDPRESS_DB_HOST=192.168.99.100 -e WORDPRESS_DB_PORT=31135 -e WORDPRESS_DB_PASSWORD=Unitone -e #SPOS_AUTH=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MzU5NDcyNTUsImlhdCI6MTUzNTk0MzY1NX0.C5axQvnEdpqiE3P2ADr4gj3i3Mvc1zLGaR2TPSelYa8 -e ADDRESS=192.168.99.100:30372 greeter-client:v1

echo "---"
echo "kubectl deploy spos-chat..."
kubectl apply -f ./spos-chat.yaml
