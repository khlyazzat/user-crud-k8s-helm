#!/bin/bash

helm install users-db oci://registry-1.docker.io/bitnamicharts/postgresql -f kube/postgres/values.yaml

kubectl apply -f kube/configs/
kubectl apply -f kube/secrets/
kubectl apply -f kube/migrations/job.yaml

kubectl wait --for=condition=complete job/user-db-migrate --timeout=60s

kubectl apply -f kube/deployment.yaml
kubectl apply -f kube/service.yaml
kubectl apply -f kube/ingress.yaml
