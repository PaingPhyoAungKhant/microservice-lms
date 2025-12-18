#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MANIFESTS_DIR="${SCRIPT_DIR}/manifests"

NAMESPACE="asto-lms"
KUBECTL="${KUBECTL:-kubectl}"

echo "Deployment process started"

echo "Creating namespace..."
"${KUBECTL}" apply -f "${MANIFESTS_DIR}/base/namespace.yaml"

echo "Creating configmap..."
"${KUBECTL}" apply -f "${MANIFESTS_DIR}/base/configmap.yaml"

echo "Creating secrets..."
"${KUBECTL}" apply -f "${MANIFESTS_DIR}/base/secrets.yaml"

echo "Deploying databases..."
"${KUBECTL}" apply -f "${MANIFESTS_DIR}/base/databases/"

echo "Waiting for databases to be ready..."
for db in user-db file-db course-db zoom-db enrollment-db; do
  if "${KUBECTL}" wait --for=condition=ready pod -l app="${db}" \
    -n "${NAMESPACE}" \
    --timeout=120s; then
    echo "${db} is ready!"
  else
    echo "${db} failed to become ready!"
    "${KUBECTL}" get pods -l app="${db}" -n "${NAMESPACE}"
    exit 1
  fi
done
echo ""

echo "Deploying RabbitMQ..."
"${KUBECTL}" apply -f "${MANIFESTS_DIR}/base/infrastructure/rabbitmqcluster.yaml"
echo "Waiting for RabbitMQ to be ready..."
sleep 20
echo ""

echo "Deploying Redis..."
"${KUBECTL}" apply -f "${MANIFESTS_DIR}/base/infrastructure/redis.yaml"
echo "Waiting for Redis to be ready..."
sleep 10
echo ""

echo "Deploying MinIO..."
"${KUBECTL}" apply -f "${MANIFESTS_DIR}/base/infrastructure/minio-tenant.yaml"
echo "Waiting for MinIO to be ready..."
sleep 15
echo ""

echo "Deploying migrate jobs..."

migrate_jobs=("user-service-migrate" "file-service-migrate" "course-service-migrate" "zoom-service-migrate" "enrollment-service-migrate")

for job in "${migrate_jobs[@]}"; do
  echo "Running ${job} migrations..."
  "${KUBECTL}" delete job "${job}" -n "${NAMESPACE}" --ignore-not-found
  "${KUBECTL}" apply -f "${MANIFESTS_DIR}/base/migrate/${job}-job.yaml"
  if "${KUBECTL}" wait --for=condition=complete "job/${job}" \
    -n "${NAMESPACE}" \
    --timeout=300s; then
    echo "${job} migrations completed successfully!"
  else
    echo "${job} migrations failed or timed out!"
    echo "Migration job logs:"
    "${KUBECTL}" logs -l job-name="${job}" -n "${NAMESPACE}" --tail=50 || true
    "${KUBECTL}" describe job "${job}" -n "${NAMESPACE}" || true
    exit 1
  fi
done
echo "Waiting for migrate jobs to be ready..."
sleep 10
echo ""

echo "Deploying services..."
"${KUBECTL}" apply -f "${MANIFESTS_DIR}/base/services"
echo "Waiting for services to be ready..."
sleep 20

echo "Waiting for frontend service to be ready..."
if "${KUBECTL}" wait --for=condition=ready pod -l app=frontend-service \
  -n "${NAMESPACE}" \
  --timeout=120s; then
  echo "Frontend service is ready!"
else
  echo "Frontend service failed to become ready!"
  "${KUBECTL}" get pods -l app=frontend-service -n "${NAMESPACE}"
  exit 1
fi
echo ""

echo "Deploying gateway..."
"${KUBECTL}" apply -f "${MANIFESTS_DIR}/base/gateway"
echo "Waiting for gateway to be ready..."
sleep 10
echo ""

echo "Deployment completed..."
echo ""
echo "Checking deployment status..."
echo ""
"${KUBECTL}" get pods -n "${NAMESPACE}"
echo ""
"${KUBECTL}" get services -n "${NAMESPACE}"
echo ""
"${KUBECTL}" get gateway -n "${NAMESPACE}"
echo ""
"${KUBECTL}" get httproute -n "${NAMESPACE}"
echo ""
"${KUBECTL}" get snippetsfilter -n "${NAMESPACE}"
echo ""

