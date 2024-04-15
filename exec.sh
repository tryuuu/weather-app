#!/bin/bash
eval $(minikube docker-env)

# v1バージョンとv2バージョンのDockerイメージを並列でビルド
docker build -f Dockerfile.v1 -t weather-app:v1 . &
docker build -f Dockerfile.v2 -t weather-app:v2 . &

# すべてのバックグラウンドプロセスが終了するのを待つ
wait

kubectl apply -f k8s/deployment.yml
kubectl apply -f k8s/service.yml

# Podの状態を表示
kubectl get pods

# MinikubeのサービスURLを取得
minikube service weather-service --url
