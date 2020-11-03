# Auth manager
A service that manages user authentication. Built with GO.

# Connecting to PG kubernetes
* Implement service NodePort for postgres deployment
* For Local Developemt
  * minikube start
  * minikube service --url [postgres kubernetes service]
  * `$ psql -h localhost -U [postgres user] --password -p [port from minikube postgres service tunnel] [postgres database]`

# Running Kubernetes in development mode
* minikube start
* tilt up -- you need to have tilt installed
