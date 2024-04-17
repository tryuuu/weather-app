#!/bin/bash
eval $(minikube docker-env)

# v1バージョンとv2バージョンのDockerイメージを並列でビルド
docker build -f Dockerfile.v1 -t weather-app:v1 . &
docker build -f Dockerfile.v2 -t weather-app:v2 . &

# すべてのバックグラウンドプロセスが終了するのを待つ
wait

kubectl create namespace development
kubectl label namespace development istio-injection=enabled # サイドカーを自動で注入
kubectl apply -f k8s/deployment.yml -n development
kubectl apply -f k8s/service.yml -n development
kubectl apply -f istio/gateway.yml -n istio-system
kubectl apply -f istio/destinationrule.yml -n development
kubectl apply -f istio/virtualservice.yml -n development

# MinikubeのサービスURLを取得
minikube service weather-service --url -n development
