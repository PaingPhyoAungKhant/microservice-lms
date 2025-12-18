#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "Build process started"

echo "Checking if Minikube is running..."
if ! minikube status | grep -q "Running"; then
  echo "Minikube is not running"
  echo "Please start minikube and try again"
  echo "minikube start"
  exit 1
fi

echo "Setting up Minikube Docker environment..."
eval "$(minikube docker-env)"

echo "Building user service image..."
docker build \
  -t asto-lms/user-service:latest \
  -f "${SCRIPT_DIR}/services/user-service/Dockerfile" \
  "${SCRIPT_DIR}"

echo "Building user service migrate image..."
docker build \
  -t asto-lms/user-service-migrate:latest \
  -f "${SCRIPT_DIR}/services/user-service/Dockerfile.migrate" \
  "${SCRIPT_DIR}"

echo "Building auth service image..."
docker build \
  -t asto-lms/auth-service:latest \
  -f "${SCRIPT_DIR}/services/auth-service/Dockerfile" \
  "${SCRIPT_DIR}"

echo "Building notification service image..."
docker build \
  -t asto-lms/notification-service:latest \
  -f "${SCRIPT_DIR}/services/notification-service/Dockerfile" \
  "${SCRIPT_DIR}"

echo "Building file service image..."
docker build \
  -t asto-lms/file-service:latest \
  -f "${SCRIPT_DIR}/services/file-service/Dockerfile" \
  "${SCRIPT_DIR}"

echo "Building file service migrate image..."
docker build \
  -t asto-lms/file-service-migrate:latest \
  -f "${SCRIPT_DIR}/services/file-service/Dockerfile.migrate" \
  "${SCRIPT_DIR}"

echo "Building course service image..."
docker build \
  -t asto-lms/course-service:latest \
  -f "${SCRIPT_DIR}/services/course-service/Dockerfile" \
  "${SCRIPT_DIR}"

echo "Building course service migrate image..."
docker build \
  -t asto-lms/course-service-migrate:latest \
  -f "${SCRIPT_DIR}/services/course-service/Dockerfile.migrate" \
  "${SCRIPT_DIR}"

echo "Building zoom service image..."
docker build \
  -t asto-lms/zoom-service:latest \
  -f "${SCRIPT_DIR}/services/zoom-service/Dockerfile" \
  "${SCRIPT_DIR}"

echo "Building zoom service migrate image..."
docker build \
  -t asto-lms/zoom-service-migrate:latest \
  -f "${SCRIPT_DIR}/services/zoom-service/Dockerfile.migrate" \
  "${SCRIPT_DIR}"

echo "Building enrollment service image..."
docker build \
  -t asto-lms/enrollment-service:latest \
  -f "${SCRIPT_DIR}/services/enrollment-service/Dockerfile" \
  "${SCRIPT_DIR}"

echo "Building enrollment service migrate image..."
docker build \
  -t asto-lms/enrollment-service-migrate:latest \
  -f "${SCRIPT_DIR}/services/enrollment-service/Dockerfile.migrate" \
  "${SCRIPT_DIR}"

echo "Building frontend image..."
docker build \
  -t asto-lms/frontend:latest \
  -f "${SCRIPT_DIR}/../frontend/Dockerfile" \
  "${SCRIPT_DIR}/../frontend"

echo "Build Process Completed"
echo "Deploy with: ./deploy.sh"

