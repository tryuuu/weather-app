## 実行
`minikube start --driver=docker`

`eval $(minikube docker-env)`

`docker build -t weather-app:latest .`

`kubectl apply -f k8s/deployment.yml`

`kubectl apply -f k8s/service.yml`

`kubectl delete -f k8s/deployment.yml`

`kubectl delete -f k8s/service.yml`

`minikube service list`

`kubectl get pods`

`minikube service weather-service --url`