#!/bin/bash
echo ""
echo "VirtualSservice"
kubectl get virtualservices -n development
echo ""
echo "DestinationRule"
kubectl get destinationrules -n development
echo ""
echo "Gateway"
kubectl get gateways --all-namespaces
echo ""
echo "Deployment"
kubectl get deployments -n development
echo ""
echo "Service"
kubectl get services -n development
echo ""
echo "Pods"
kubectl get pods -n development