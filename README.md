# user-crud-k8s-helm-prometheus

### 1: Install PostgreSQL via Helm with custom values
helm install users-db oci://registry-1.docker.io/bitnamicharts/postgresql -f kube/postgres/values.yaml

### 2: Apply ConfigMaps and Secrets required by the app and migration job
kubectl apply -f kube/configs/
kubectl apply -f kube/secrets/

### 3: Run the initial database migration Job
kubectl apply -f kube/migrations/job.yaml

### 4: Wait for the migration Job to complete
kubectl wait --for=condition=complete job/user-db-migrate --timeout=60s

### 5: Deploy the application components
kubectl apply -f kube/deployment.yaml
kubectl apply -f kube/service.yaml
kubectl apply -f kube/ingress.yaml

### 6: Install the kube-prometheus-stack via Helm
helm install prometheus prometheus-community/kube-prometheus-stack \
  --namespace monitoring \
  --create-namespace \
  --set prometheus.prometheusSpec.maximumStartupDurationSeconds=300

### 7: Connect custom services to Prometheus
kubectl apply -f kube/service.yaml
kubectl apply -f kube/monitoring/servicemonitor-ingress.yaml
kubectl apply -f kube/monitoring/servicemonitor-user.yaml

### 8: Accessing Prometheus, Grafana, and Service Metrics

### Verify Prometheus targets
kubectl get svc -n monitoring
kubectl port-forward svc/prometheus-kube-prometheus-prometheus -n monitoring 9090:9090
Open http://localhost:9090/targets

### View application metrics directly
kubectl port-forward svc/user-service 8080
http://localhost:8080/metrics

### Access Grafana UI
kubectl port-forward svc/prometheus-grafana -n monitoring 3000:80
http://localhost:3000

kubectl get secret -n monitoring prometheus-grafana -o jsonpath="{.data.admin-user}" | base64 -d
echo

kubectl get secret -n monitoring prometheus-grafana -o jsonpath="{.data.admin-password}" | base64 -d
echo
