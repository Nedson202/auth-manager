# Auth manager
Simple auth API

# Connecting to PG kubernetes
* For Local Developemt
  * kubectl exec -it [PG_POD] -- psql -h localhost -p 5432 -U [PG_USER] [PG_DB]

# Running Kubernetes in development mode
* tilt up -- you need to have tilt installed

# Proto generation
### Generate proto
protoc -I ./api/proto/v1 \
    -I. -I$GOPATH/src/github.com/nedson202/auth-manager/third_party \
    --go_out ./api/proto/v1 --go_opt=paths=source_relative \
    --grpc-gateway_out=logtostderr=true:api/proto/v1 --grpc-gateway_opt paths=source_relative \
    api/proto/v1/auth.proto

### Generate OpenApi spec
protoc -I ./api/proto/v1 \
    -I. -I$GOPATH/src/github.com/nedson202/auth-manager/third_party \
    --go-grpc_out ./api/proto/v1 --go-grpc_opt paths=source_relative \
    --swagger_out=logtostderr=true:api/swagger/v1 \
    api/proto/v1/auth.proto

# Managing secrets
Secrets are managed in GitOps via SealedSecrets: https://engineering.bitnami.com/articles/sealed-secrets.html
* Install commandline tool 
  * brew install kubeseal
* Generate SealedSecret definition
  * kubeseal --format=yaml --cert=kube/public-key-cert.pem <kube/[KUBERNETES SECRET DEFINITION] >kube/[SEALED SECRET OUTPUT].yaml
