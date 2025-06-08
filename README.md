# user-crud-k8s-helm-prometheus

### 1: Install the kube-prometheus-stack via Helm
helm install prometheus prometheus-community/kube-prometheus-stack \
  --namespace monitoring \
  --create-namespace

helm upgrade prometheus prometheus-community/kube-prometheus-stack \
  -n monitoring \
  -f kube/values/prometheus-values.yaml

### 2: Connect custom services to Prometheus
kubectl apply -f kube/service.yaml

kubectl apply -f kube/monitoring/servicemonitor-ingress.yaml

kubectl apply -f kube/monitoring/servicemonitor-user.yaml  

### 3: Install PostgreSQL via Helm with custom values
helm install users-db oci://registry-1.docker.io/bitnamicharts/postgresql -f kube/values/postgres-values.yaml

### 4: Apply ConfigMaps and Secrets required by the app
kubectl apply -f kube/configs/

kubectl apply -f kube/secrets/

### 5: Run the initial database migration Job
kubectl apply -f kube/migrations/job.yaml

### 6: Wait for the migration Job to complete
kubectl wait --for=condition=complete job/user-db-migrate --timeout=60s

### 7: Deploy the application components
kubectl apply -f kube/deployment.yaml

kubectl apply -f kube/service.yaml

kubectl apply -f kube/ingress.yaml

### 8: Accessing Prometheus, Grafana, and Service Metrics

### Verify Prometheus targets
kubectl get svc -n monitoring

kubectl port-forward svc/prometheus-kube-prometheus-prometheus -n monitoring 9090:9090

http://localhost:9090/targets

### View application metrics directly
kubectl port-forward svc/user-service 8080

http://localhost:8080/metrics

### Access Grafana UI
kubectl port-forward svc/prometheus-grafana -n monitoring 3000:80

http://localhost:3000

kubectl get secret -n monitoring prometheus-grafana -o jsonpath="{.data.admin-user}" | base64 -d && \
echo

kubectl get secret -n monitoring prometheus-grafana -o jsonpath="{.data.admin-password}" | base64 -d && \
echo

### 9: PostgreSQL Monitoring via Prometheus Exporter
pg_stat_activity_count

pg_stat_database_numbackends

pg_stat_database_blks_hit

pg_stat_database_blks_read

pg_stat_database_tup_returned

http://localhost:9090/query

### 10: Enable Metrics for NGINX Ingress Controller. Install the NGINX Ingress Controller via Helm
kubectl get all -A | grep ingress-nginx   # should return nothing

helm list -A | grep ingress               # should return nothing

kubectl get ns | grep ingress             # should return nothing

#### if not 

minikube addons disable ingress  # disable built-in ingress if used

helm uninstall ingress-nginx -n ingress-nginx || true

kubectl delete ns ingress-nginx --grace-period=0 --force

#### then

helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx

helm repo update

helm install ingress-nginx ingress-nginx/ingress-nginx \
--namespace ingress-nginx \
--create-namespace \
-f kube/values/ingress-values.yaml

### 11: Accessing NGINX Ingress Controller Metrics
kubectl port-forward -n ingress-nginx svc/ingress-nginx-controller-metrics 10254:10254

localhost:10254/metrics

Minikube Note: LoadBalancer workaround
Minikube doesn’t have a real cloud LoadBalancer like AWS or GCP.
So when you create a Kubernetes service of type LoadBalancer, it won’t receive a public IP by default.
Instead of relying on minikube tunnel (which requires root privileges, can hang, and may break after reboots),
this setup uses a manual externalIPs patch to simulate external access:

kubectl patch svc ingress-nginx-controller -n ingress-nginx \
-p '{"spec": {"externalIPs": ["$(minikube ip)"]}}'

kubectl get svc -n ingress-nginx ingress-nginx-controller

This allows you to access your Ingress via domain (e.g., arch.homework) without using minikube tunnel,
while still keeping the service type as LoadBalancer.

curl -i http://.../health 

http://localhost:9090/query

in table

nginx_ingress_controller_requests

in graph

{__name__=~"nginx_ingress_controller_.*"}

