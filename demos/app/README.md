# Dockerize the Go Application for minikube

**Build the Docker image**

```bash
docker build --tag $(minikube ip):56540/secret-app:latest .
```