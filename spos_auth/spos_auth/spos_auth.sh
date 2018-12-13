#!/usr/bin/env bash	 

#echo "kubectl delete spos-auth service..."
#kubectl delete svc spos-auth

#echo "kubectl delete spos-auth pod..."
#kubectl delete po --force --grace-period=0 -l app=spos-auth 
#kubectl delete po -l app=spos-auth

#echo "kubectl delete spos-auth deploy..."
#kubectl delete deploy spos-auth

echo "---"
echo "kubectl delete spos-auth service deploy pod..."
kubectl delete -f ./spos-auth.yaml

echo "---"
echo "build spos_auth to bin..."
CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o spos-auth .


echo "---"
echo "delete spos-auth docker image..."
docker rmi --force patlihh/spos-auth:v1.0

echo "---"
echo "build spos-auth docker image..."
docker build -t patlihh/spos-auth:v1.0 .

echo "---"
echo "push spos-auth docker image..."
docker push patlihh/spos-auth:v1.0

echo "---"
#echo "docker run greeter-client..."
#docker run -it -e WORDPRESS_DB_HOST=192.168.99.100 -e WORDPRESS_DB_PORT=31135 -e WORDPRESS_DB_PASSWORD=Unitone -e #SPOS_AUTH=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MzU5NDcyNTUsImlhdCI6MTUzNTk0MzY1NX0.C5axQvnEdpqiE3P2ADr4gj3i3Mvc1zLGaR2TPSelYa8 -e ADDRESS=192.168.99.100:30372 greeter-client:v1

echo "---"
echo "kubectl deploy spos-auth..."
kubectl apply -f ./spos-auth.yaml
