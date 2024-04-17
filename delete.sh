#!/bin/bash
kubectl delete -f k8s/deployment.yml -n development
kubectl delete -f k8s/service.yml -n development
kubectl delete -f istio/gateway.yml -n istio-system
kubectl delete -f istio/destinationrule.yml -n development
kubectl delete -f istio/virtualservice.yml -n development