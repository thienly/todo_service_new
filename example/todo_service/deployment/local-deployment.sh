#!/bin/bash
#start min
#minikube start --driver=virtual
# add helm repo
#helm repo add bitnami https://charts.bitnami.com/bitnami

if [[ "$(helm status my-mongodb 2>/dev/null)" == "" ]]; then
  echo "mongodb helm chart is not existed"
  helm upgrade --install  my-mongodb bitnami/mongodb
fi
# get the mongodb username and password.
eval $(minikube docker-env) && docker build -t todo_service:0.0.1 -f ../build/Dockerfile ../
