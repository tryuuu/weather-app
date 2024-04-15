#!/bin/bash
eval $(minikube docker-env)
docker build -t weather-app:latest .
# 既存のデプロイメントがあれば削除
if kubectl get deployment | grep -q 'weather-deployment'; then
  kubectl delete -f k8s/deployment.yml
fi

# 既存のサービスがあれば削除
if kubectl get service | grep -q 'weather-service'; then
  kubectl delete -f k8s/service.yml
fi
kubectl apply -f k8s/deployment.yml
kubectl apply -f k8s/service.yml
kubectl get pods
minikube service weather-service --url