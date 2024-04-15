#!/bin/bash
eval $(minikube docker-env)
# v1バージョンのDockerイメージをビルド
docker build -f Dockerfile.v1 -t weather-app:v1 .
# v2バージョンのDockerイメージをビルド
docker build -f Dockerfile.v2 -t weather-app:v2 .
kubectl apply -f k8s/deployment.yml
kubectl apply -f k8s/service.yml
kubectl get pods
minikube service weather-service --url