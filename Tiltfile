# -*- mode: Go -*-

k8s_yaml('kube/postgres-config.yaml')
k8s_yaml('kube/postgres.yaml')
k8s_yaml('kube/auth-manager-deployment.yaml')
k8s_yaml('kube/auth-manager-service.yaml')
k8s_resource('user-api', port_forwards=5000)
