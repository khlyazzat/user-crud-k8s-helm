# user-crud-k8s-helm

1: Install PostgreSQL via Helm with custom values
helm install users-db oci://registry-1.docker.io/bitnamicharts/postgresql -f kube/postgres/values.yaml

2: Apply ConfigMaps and Secrets required by the app and migration job
kubectl apply -f kube/configs/
kubectl apply -f kube/secrets/

3: Run the initial database migration Job
kubectl apply -f kube/migrations/job.yaml

4: Wait for the migration Job to complete
kubectl wait --for=condition=complete job/user-db-migrate --timeout=60s

5: Deploy the application components
kubectl apply -f kube/deployment.yaml
kubectl apply -f kube/service.yaml
kubectl apply -f kube/ingress.yaml
