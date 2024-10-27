# Setting Up a Local Kubernetes Cluster with Minikube

This guide explains how to set up a local Kubernetes cluster on your development machine using Minikube.

## Prerequisites

1. **Virtualization Software**: Minikube requires a hypervisor. Options include VirtualBox, Docker, HyperKit (macOS), or Hyper-V (Windows).
2. **Minikube**: Install Minikube as shown below.
3. **Kubectl**: The Kubernetes command-line tool.

## Installation

### MacOS

```bash
# Installing docker as a container runtime
brew install --cask docker
```

```bash
# Installing minikube as a kubernetes runtime
brew install minikube
```

## Starting minikube

```bash
minikube start --driver=docker
```

```bash
# Verifying that minikube is working
minikube status
```

## Enable minikube registry

```bash
minikube addons enable registry
```

## Switch to Minikube's Docker daemon

This is necessary to build the image inside Minikube and make them available right away.

```bash
eval $(minikube docker-env)
```

**Switching back:**

```bash
eval $(minikube docker-env --unset)
```
### Cache images in minikube

```bash
minikube image load alpine:latest
```
