# ASTO LMS - Deployment Guide

This guide provides step-by-step instructions for deploying and running the ASTO LMS system. This guide focuses on deployment and running the system, not development setup.

## Table of Contents

- [Prerequisites](#prerequisites)
- [System Requirements](#system-requirements)
- [Installation Steps](#installation-steps)
- [Building Docker Images](#building-docker-images)
- [Deploying the System](#deploying-the-system)
- [Verifying Deployment](#verifying-deployment)
- [Accessing the System](#accessing-the-system)
- [Troubleshooting](#troubleshooting)
- [Uninstalling](#uninstalling)

---

## Prerequisites

### Required Software

1. **Kubernetes Cluster**
   - **Minikube** (for local development/testing)
   - OR **Production Kubernetes cluster** (for production deployment)
   - Kubernetes version: 1.24 or higher

2. **kubectl**
   - Kubernetes command-line tool
   - Version: 1.24 or higher
   - Must be configured to connect to your cluster
   official link <https://kubernetes.io/docs/tasks/tools/>

3. **Docker**
   - For building container images
   - Version: 20.10 or higher
   - Must be running and accessible
4. **Helm**
   - For kubernetes manifests configurations
   - Official documentation
   - <https://helm.sh/docs/intro/install/>

4. **Bash Shell**
   - Required for running deployment scripts
   - Available on Linux, macOS, and Windows (WSL/Git Bash)
   - It is recommended to deploy on linux environment. This guide is focus on deployment inside linux environment.
     To install wsl inside windows use the following commend in Power shell.

  ```bash
  # On PowerShell
  wsl --install

### Minimum System Resources

For **Minikube**:
- **CPU**: 4 cores minimum
- **Memory**: 8GB RAM minimum (16GB recommended)
- **Disk Space**: 20GB free space

---

## System Requirements

### Kubernetes Operators and APIs Required

The system requires the following Kubernetes operators and APIs to be installed:

1. **Gateway API CRDs**
   - Required for Gateway API functionality
   - CRDs: `gateways.gateway.networking.k8s.io`, `httproutes.gateway.networking.k8s.io`, etc.
   - Must be installed before NGINX Gateway Controller

2. **Redis Operator** (Opstree Labs)
   - Manages Redis instances
   - CRD: `redis.redis.opstreelabs.in/v1beta2`

3. **MinIO Operator**
   - Manages MinIO object storage
   - CRD: `minio.min.io/v2`

4. **RabbitMQ Operator**
   - Manages RabbitMQ clusters
   - CRD: `rabbitmq.com/v1beta1`

5. **NGINX Gateway Controller**
   - Provides Gateway API implementation
   - Required for API Gateway functionality

---

## Installation Steps
### Step 1: Install Docker 
Docker can be installed by following the official docummentation
- For docker desktop
https://docs.docker.com/get-started/get-docker/
- For docker engine 
https://docs.docker.com/engine/



### Step 2: Set Up Kubernetes Cluster 

#### Using Minikube (Local Development)

```bash
# Install CURL
sudo apt update 
sudo apt install curl

# Install Minikube (if not already installed)
# On Linux:
curl -LO https://github.com/kubernetes/minikube/releases/latest/download/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube && rm minikube-linux-amd64


# Start Minikube with sufficient resources
minikube start --cpus=4 --memory=8192 --disk-size=20g

# Verify Minikube is running
minikube status

# Enable Minikube addons for PersistentVolumeClaims (PostgreSQL databases, Redis, MinIO)
minikube addons enable storage-provisioner

```

```bash
# Verify cluster access
kubectl cluster-info

# Verify nodes are ready
kubectl get nodes

# Ensure you have cluster admin permissions
kubectl auth can-i '*' '*' --all-namespaces
```

### Step 3: Install Redis Operator

The Redis Operator is required to manage Redis instances.
Offiial Documentaion
<https://redis-operator.opstree.dev/docs/installation/installation/>

```bash
# Install Helm 
sudo apt-get install curl gpg apt-transport-https --yes
curl -fsSL https://packages.buildkite.com/helm-linux/helm-debian/gpgkey | gpg --dearmor | sudo tee /usr/share/keyrings/helm.gpg > /dev/null
echo "deb [signed-by=/usr/share/keyrings/helm.gpg] https://packages.buildkite.com/helm-linux/helm-debian/any/ any main" | sudo tee /etc/apt/sources.list.d/helm-stable-debian.list
sudo apt-get update
sudo apt-get install helm


# Install Redis Operator using Helm 
helm repo add ot-helm https://ot-container-kit.github.io/helm-charts/
helm repo update
helm install redis-operator ot-helm/redis-operator --namespace ot-operators 


# Verify installation
$ kubectl describe --namespace ot-operators pods
```

### Step 4: Install MinIO Operator

The MinIO Operator is required to manage MinIO object storage.

Official Documentation
<https://github.com/minio/operator>

```bash
# Install MinIO Operator using kubectl kustomize
kubectl kustomize github.com/minio/operator\?ref=v7.1.1 | kubectl apply -f -


# Verify installation
kubectl get pods -n minio-operator
```

**Note:** The MinIO Operator creates a namespace `minio-operator` automatically.

### Step 5: Install RabbitMQ Operator

The RabbitMQ Operator is required to manage RabbitMQ clusters.

```bash
# Install RabbitMQ Operator using kubectl
kubectl apply -f "https://github.com/rabbitmq/cluster-operator/releases/latest/download/cluster-operator.yml"


# Verify installation
kubectl get pods -n rabbitmq-system
```

**Add namespace**

```bash
# run this commend to add asto-lms namespace to rabbitmq cluster operator 
kubectl -n rabbitmq-system edit deployment rabbitmq-cluster-operator

# After running the commend add the asto-lms namespace 
# the following is an example 
-------------
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app.kubernetes.io/component: rabbitmq-operator
    app.kubernetes.io/name: rabbitmq-cluster-operator
    app.kubernetes.io/part-of: rabbitmq
  name: rabbitmq-cluster-operator
  namespace: rabbitmq-system
spec:
  template:
    spec:
      containers:
      - command:
        - /manager
        env:
        - name: OPERATOR_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        - name: OPERATOR_SCOPE_NAMESPACE
          value: "sample-namespace,asto-lms"
# ...

# ensure availity 
kubectl get customresourcedefinitions.apiextensions.k8s.io

```

**Note:** The RabbitMQ Operator creates a namespace `rabbitmq-system` automatically.

### Step 6: Install Gateway API CRDs

Gateway API CRDs are required before installing the NGINX Gateway Controller. Minikube doesn't include these by default.

Official Documentation
<https://gateway-api.sigs.k8s.io/guides/#installing-gateway-api>

```bash
# Install Gateway API CRDs
kubectl apply --server-side -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.4.0/standard-install.yaml

# Verify Gateway API CRDs are installed
kubectl get crd gateways.gateway.networking.k8s.io
kubectl get crd httproutes.gateway.networking.k8s.io

# Should show both CRDs
```

### Step 7: Install NGINX Gateway Controller

The NGINX Gateway Controller provides Gateway API implementation for the API Gateway.

```bash
# Install NGINX Gateway Controller
kubectl apply -f https://raw.githubusercontent.com/nginxinc/nginx-kubernetes-gateway/main/deploy/manifests/install.yaml

# Wait for controller to be ready
kubectl wait --for=condition=ready pod \
  -l app=nginx-kubernetes-gateway \
  -n nginx-gateway \
  --timeout=300s

# Verify installation
kubectl get pods -n nginx-gateway
```

**Verification:**

```bash
# Check if Gateway CRD is available
kubectl get crd gateways.gateway.networking.k8s.io

# Should show: gateways.gateway.networking.k8s.io
```

**Alternative: Install via Helm**

```bash
helm repo add nginx-stable https://helm.nginx.com/stable
helm repo update
helm install nginx-gateway nginx-stable/nginx-kubernetes-gateway \
  --namespace nginx-gateway \
  --create-namespace
```

---

## Building Docker Images

Before deploying, you need to build Docker images for all services and the frontend.

### Prerequisites for Building

1. **Docker must be running**
2. **For Minikube**: Docker must be configured to use Minikube's Docker daemon
3. **For Production**: Docker images should be pushed to a container registry

### Building Images

#### For Minikube (Local Development)

```bash
# Navigate to backend directory
cd /home/paingkhant/projects/asto-lms/backend

# Set up Minikube Docker environment
eval $(minikube docker-env)

# Run the build script
./build.sh
```

The build script will:

- Build all microservice images (user-service, auth-service, course-service, etc.)
- Build all migration job images
- Build the frontend image
- Tag all images as `asto-lms/<service-name>:latest`

**Build Time**: Approximately 10-15 minutes depending on system performance.

## Deploying the System

### Step 1: Configure Secrets

Before deployment, you must configure secrets in the Kubernetes cluster.

```bash
# Navigate to backend directory
cd /home/paingkhant/projects/asto-lms/backend

# Edit the secrets file
nano manifests/base/secrets.yaml
```

**Required Secrets:**

- Database passwords (POSTGRES_USER, POSTGRES_PASSWORD)
- Redis password (REDIS_PASSWORD)
- JWT secret key (JWT_SECRET_KEY)
- MinIO access keys (MINIO_ACCESS_KEY, MINIO_SECRET_KEY)
- Zoom API credentials (ZOOM_ACCOUNT_ID, ZOOM_CLIENT_ID, ZOOM_CLIENT_SECRET, ZOOM_SECRET_TOKEN)

**Example secrets.yaml structure:**

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: asto-lms-secrets
  namespace: asto-lms
type: Opaque
stringData:
  POSTGRES_USER: "postgres"
  POSTGRES_PASSWORD: "your-secure-password"
  REDIS_PASSWORD: "your-redis-password"
  JWT_SECRET_KEY: "your-jwt-secret-key-min-32-chars"
  MINIO_ACCESS_KEY: "your-minio-access-key"
  MINIO_SECRET_KEY: "your-minio-secret-key"
  ZOOM_ACCOUNT_ID: "your-zoom-account-id"
  ZOOM_CLIENT_ID: "your-zoom-client-id"
  ZOOM_CLIENT_SECRET: "your-zoom-client-secret"
  ZOOM_SECRET_TOKEN: "your-zoom-secret-token"
```

### Step 2: Configure ConfigMap (Optional)

Review and update the ConfigMap if needed:

```bash
# Edit configmap if needed
nano manifests/base/configmap.yaml
```

The ConfigMap contains non-sensitive configuration like:

- Service ports
- Database hosts
- Environment settings
- API Gateway URLs

### Step 3: Run Deployment Script

```bash
# Navigate to backend directory
cd /home/paingkhant/projects/asto-lms/backend

# Make deploy script executable (if not already)
chmod +x deploy.sh

# Run deployment
./deploy.sh
```

**What the deployment script does:**

1. Creates namespace `asto-lms`
2. Creates ConfigMap and Secrets
3. Deploys PostgreSQL databases (5 databases)
4. Waits for databases to be ready
5. Deploys RabbitMQ cluster
6. Deploys Redis instance
7. Deploys MinIO tenant
8. Runs database migration jobs
9. Deploys all microservices
10. Deploys frontend service
11. Deploys API Gateway and routes

**Deployment Time**: Approximately 5-10 minutes depending on cluster performance.

### Step 4: Manual Deployment (Alternative)

If you prefer to deploy manually:

```bash
cd /home/paingkhant/projects/asto-lms/backend/manifests/base

# Create namespace
kubectl apply -f namespace.yaml

# Create ConfigMap and Secrets
kubectl apply -f configmap.yaml
kubectl apply -f secrets.yaml

# Deploy databases
kubectl apply -f databases/

# Wait for databases
kubectl wait --for=condition=ready pod -l app=user-db -n asto-lms --timeout=120s
# Repeat for other databases...

# Deploy infrastructure
kubectl apply -f infrastructure/

# Wait for infrastructure
kubectl wait --for=condition=ready pod -l app=rabbitmq -n asto-lms --timeout=300s
# Wait for Redis and MinIO...

# Run migrations
kubectl apply -f migrate/

# Wait for migrations to complete
kubectl wait --for=condition=complete job/user-service-migrate -n asto-lms --timeout=300s
# Repeat for other migration jobs...

# Deploy services
kubectl apply -f services/

# Deploy gateway
kubectl apply -f gateway/
```

---

## Verifying Deployment

### Check Operator Status

```bash
# Check Redis Operator
kubectl get pods -n redis-operator

# Check MinIO Operator
kubectl get pods -n minio-operator

# Check RabbitMQ Operator
kubectl get pods -n rabbitmq-system

# Check NGINX Gateway Controller
kubectl get pods -n nginx-gateway
```

All operators should show `Running` status.

### Check Application Status

```bash
# Check all pods in asto-lms namespace
kubectl get pods -n asto-lms

# Check services
kubectl get services -n asto-lms

# Check deployments
kubectl get deployments -n asto-lms

# Check gateway
kubectl get gateway -n asto-lms

# Check HTTPRoutes
kubectl get httproute -n asto-lms

# Check infrastructure resources
kubectl get redis -n asto-lms
kubectl get tenant -n asto-lms
kubectl get rabbitmqcluster -n asto-lms
```

### Verify Pods Are Running

```bash
# All pods should be in Running state
kubectl get pods -n asto-lms

# Expected pods:
# - user-db-*
# - file-db-*
# - course-db-*
# - zoom-db-*
# - enrollment-db-*
# - rabbitmq-*
# - asto-redis-*
# - asto-minio-*
# - user-service-*
# - auth-service-*
# - course-service-*
# - enrollment-service-*
# - file-service-*
# - notification-service-*
# - zoom-service-*
# - frontend-service-*
```

### Check Pod Logs

```bash
# Check service logs
kubectl logs -n asto-lms deployment/user-service --tail=50
kubectl logs -n asto-lms deployment/auth-service --tail=50
kubectl logs -n asto-lms deployment/frontend-service --tail=50

# Check database logs
kubectl logs -n asto-lms deployment/user-db --tail=50

# Check infrastructure logs
kubectl logs -n asto-lms deployment/rabbitmq --tail=50
```

### Verify Database Migrations

```bash
# Check migration job status
kubectl get jobs -n asto-lms

# All migration jobs should show COMPLETIONS: 1/1

# Check migration logs
kubectl logs -n asto-lms job/user-service-migrate
kubectl logs -n asto-lms job/course-service-migrate
# etc.
```

---

## Accessing the System

### For Minikube

```bash
# Check if Apache is running
sudo systemctl status apache2
# or
sudo service apache2 status

# Stop Apache
sudo systemctl stop apache2
# or
sudo service apache2 stop

# Disable it from starting on boot (optional)
sudo systemctl disable apache2

#Summary: Commands to Enable SnippetsFilters
#1. Enable SnippetsFilters Flag in Gateway Controller
# Add --snippets-filters flag to nginx-gateway container
kubectl patch deployment nginx-gateway -n nginx-gateway --type='json' -p='[{"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value": "--snippets-filters"}]'

#2. Update RBAC Permissions - Add snippetsfilters/status
# Add snippetsfilters/status resource with update verb
kubectl patch clusterrole nginx-gateway --type='json' -p='[{"op": "add", "path": "/rules/-", "value": {"apiGroups": ["gateway.nginx.org"], "resources": ["snippetsfilters/status"], "verbs": ["update"]}}]'

# Get Minikube IP
minikube ip

# Add to /etc/hosts (Linux/macOS) or C:\Windows\System32\drivers\etc\hosts (Windows)
# Add line: <minikube-ip> asto-lms.local

# Example:
# 192.168.49.2 asto-lms.local

# start the tunnel to expose external ip
minikube tunnel 

# Note: You may need to enter password when sudo permission is asked

# After starting the tunnel, restart the NGINX Gateway Controller
kubectl rollout restart deployment/nginx-kubernetes-gateway -n nginx-gateway

# Wait for the gateway to restart
kubectl wait --for=condition=available deployment/nginx-kubernetes-gateway -n nginx-gateway --timeout=120s

# Access the application
# Open browser: http://asto-lms.local

# Troubleshooting: If accessing the website returns 404, see "Minikube Tunnel Issues and Gateway Restart" 
# in the Troubleshooting section
```

---

## Troubleshooting

### Common Issues

#### 1. Operators Not Installed

**Symptoms:**

```
error: unable to recognize "redis.yaml": no matches for kind "Redis"
```

**Solution:**

```bash
# Verify operator is installed
kubectl get crd redis.redis.opstreelabs.in

# If not found, install Redis Operator (see Step 2)
```

#### 2. Pods Stuck in Pending

**Symptoms:**

```
kubectl get pods shows STATUS: Pending
```

**Solution:**

```bash
# Check pod events
kubectl describe pod <pod-name> -n asto-lms

# Common causes:
# - Insufficient resources (CPU/memory)
# - PersistentVolumeClaim not bound
# - Image pull errors

# Check resource availability
kubectl top nodes

# Check PVC status
kubectl get pvc -n asto-lms
```

#### 3. Database Connection Errors

**Symptoms:**

```
Service logs show: "failed to connect to database"
```

**Solution:**

```bash
# Verify database pods are running
kubectl get pods -l app=user-db -n asto-lms

# Check database logs
kubectl logs -n asto-lms deployment/user-db

# Verify database service
kubectl get service user-db-service -n asto-lms

# Test connection from service pod
kubectl exec -n asto-lms deployment/user-service -- \
  nc -zv user-db-service 5432
```

#### 4. Migration Jobs Failing

**Symptoms:**

```
Migration job shows STATUS: Failed
```

**Solution:**

```bash
# Check migration job logs
kubectl logs -n asto-lms job/user-service-migrate

# Common issues:
# - Database not ready (wait longer)
# - Migration files missing
# - Database credentials incorrect

# Re-run migration
kubectl delete job user-service-migrate -n asto-lms
kubectl apply -f manifests/base/migrate/user-service-migrate-job.yaml
```

#### 5. Gateway Not Routing

**Symptoms:**

```
404 errors when accessing routes
```

**Solution:**

```bash
# Check gateway status
kubectl get gateway -n asto-lms
kubectl describe gateway asto-lms-gateway -n asto-lms

# Check HTTPRoutes
kubectl get httproute -n asto-lms
kubectl describe httproute <route-name> -n asto-lms

# Check NGINX Gateway Controller logs
kubectl logs -n nginx-gateway deployment/nginx-kubernetes-gateway
```

#### 6. Minikube Tunnel Issues and Gateway Restart

**Symptoms:**

```
- Website returns 404 errors after starting minikube tunnel
- Gateway not accessible after tunnel is established
- Need to restart gateway after opening tunnel
```

**Solution:**

```bash
# 1. Start minikube tunnel (in a separate terminal)
minikube tunnel

# Note: You may need to enter password when sudo permission is asked

# 2. If accessing the website returns 404, restart the NGINX Gateway Controller
kubectl rollout restart deployment/nginx-kubernetes-gateway -n nginx-gateway

# Wait for the deployment to restart
kubectl wait --for=condition=available deployment/nginx-kubernetes-gateway -n nginx-gateway --timeout=120s

# 3. If issues persist, restart the Gateway resource
kubectl delete gateway asto-lms-gateway -n asto-lms
kubectl apply -f manifests/base/gateway/lms-gateway.yaml

# 4. If still having issues, close the tunnel and restart:
# - Press Ctrl+C to stop minikube tunnel
# - Wait a few seconds
# - Start tunnel again: minikube tunnel
# - Restart gateway controller: kubectl rollout restart deployment/nginx-kubernetes-gateway -n nginx-gateway
# - Sometimes you may need to repeat this process 2-3 times
```

#### 7. Images Not Found

**Symptoms:**

```
ImagePullBackOff or ErrImagePull errors
```

**Solution:**

```bash
# For Minikube: Ensure Docker is using Minikube's daemon
eval $(minikube docker-env)
docker images | grep asto-lms

# For Production: Ensure images are pushed to registry
# Update image pull secrets if using private registry
```

#### 8. RabbitMQ Not Ready

**Symptoms:**

```
RabbitMQ pod not starting
```

**Solution:**

```bash
# Check RabbitMQ operator
kubectl get pods -n rabbitmq-system

# Check RabbitMQ cluster status
kubectl get rabbitmqcluster -n asto-lms
kubectl describe rabbitmqcluster rabbitmq -n asto-lms

# Check RabbitMQ pod logs
kubectl logs -n asto-lms -l app=rabbitmq
```

#### 9. Redis Not Ready

**Symptoms:**

```
Redis pod not starting
```

**Solution:**

```bash
# Check Redis operator
kubectl get pods -n redis-operator

# Check Redis resource status
kubectl get redis -n asto-lms
kubectl describe redis asto-redis -n asto-lms

# Check Redis pod logs
kubectl logs -n asto-lms -l app=redis
```

#### 10. MinIO Not Ready

**Symptoms:**

```
MinIO pod not starting
```

**Solution:**

```bash
# Check MinIO operator
kubectl get pods -n minio-operator

# Check MinIO tenant status
kubectl get tenant -n asto-lms
kubectl describe tenant asto-minio -n asto-lms

# Check MinIO pod logs
kubectl logs -n asto-lms -l app=minio
```

### Getting Help

**Check Resource Status:**

```bash
# Overall status
kubectl get all -n asto-lms

# Resource usage
kubectl top pods -n asto-lms
kubectl top nodes

# Events
kubectl get events -n asto-lms --sort-by='.lastTimestamp'
```

**View Detailed Information:**

```bash
# Describe resources for detailed status
kubectl describe pod <pod-name> -n asto-lms
kubectl describe deployment <deployment-name> -n asto-lms
kubectl describe service <service-name> -n asto-lms
```

---

## Uninstalling

### Remove Application

```bash
# Delete all application resources
kubectl delete namespace asto-lms

# Or delete individually
cd /home/paingkhant/projects/asto-lms/backend/manifests/base
kubectl delete -f gateway/
kubectl delete -f services/
kubectl delete -f migrate/
kubectl delete -f infrastructure/
kubectl delete -f databases/
kubectl delete -f secrets.yaml
kubectl delete -f configmap.yaml
kubectl delete -f namespace.yaml
```

### Remove Operators (Optional)

```bash
# Remove NGINX Gateway Controller
kubectl delete -f https://raw.githubusercontent.com/nginxinc/nginx-kubernetes-gateway/main/deploy/manifests/install.yaml

# Remove RabbitMQ Operator
kubectl delete -f https://github.com/rabbitmq/cluster-operator/releases/latest/download/cluster-operator.yml

# Remove MinIO Operator
kubectl delete -f https://github.com/minio/operator/releases/latest/download/minio-operator.yaml

# Remove Redis Operator (if installed via Helm)
helm uninstall redis-operator -n redis-operator

# Or if installed via YAML
kubectl delete -f <path-to-redis-operator-manifests>
```

### Clean Up Minikube (If Using)

```bash
# Stop Minikube
minikube stop

# Delete Minikube cluster
minikube delete
```

---

## Quick Start Checklist

Use this checklist to ensure all steps are completed:

- [ ] Kubernetes cluster is running and accessible
- [ ] kubectl is installed and configured
- [ ] Docker is installed and running
- [ ] Redis Operator is installed and running
- [ ] MinIO Operator is installed and running
- [ ] RabbitMQ Operator is installed and running
- [ ] NGINX Gateway Controller is installed and running
- [ ] Docker images are built (or available in registry)
- [ ] Secrets are configured in `manifests/base/secrets.yaml`
- [ ] ConfigMap is reviewed (if needed)
- [ ] Deployment script is executed (`./deploy.sh`)
- [ ] All pods are in Running state
- [ ] Database migrations completed successfully
- [ ] Gateway is accessible
- [ ] Application is accessible via browser

---

---

## Additional Resources

### Operator Documentation

- **Redis Operator**: <https://github.com/OT-CONTAINER-KIT/redis-operator>
- **MinIO Operator**: <https://github.com/minio/operator>
- **RabbitMQ Operator**: <https://github.com/rabbitmq/cluster-operator>
- **NGINX Gateway Controller**: <https://github.com/nginxinc/nginx-kubernetes-gateway>

### Kubernetes Documentation

- **Kubernetes Basics**: <https://kubernetes.io/docs/tutorials/>
- **kubectl Reference**: <https://kubernetes.io/docs/reference/kubectl/>
- **Minikube Guide**: <https://minikube.sigs.k8s.io/docs/>

### Troubleshooting Resources

- **Kubernetes Troubleshooting**: <https://kubernetes.io/docs/tasks/debug/>
- **Operator Issues**: Check operator GitHub repositories for known issues

---

---

**Version**: 1.0
