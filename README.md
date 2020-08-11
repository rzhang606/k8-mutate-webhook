# Mutating Admission Webhook for Kubernetes

## Structure

**Deployments:**
 - deployment.yaml: the actual server
   - webhooks need TLS certificate signed by kube, refer to the webhook-create-cert script for secret creation
 - webhook configuration: catch pod creation and updates, cabundle is filled by deploy.sh
 - service: port assigned to 443 from 8443 on container

**Go files**
 - main.go
    - main server definition
    - route handlers
 - mutator.go
    - non-network related logic for the handlers
    - separated for easier unit testing

**Other**
 - webhook-create-cert.sh
   - generates certificates, creates and approves kubernetes CertificateSigningRequest, then creates a kubernetes secret
   - the secret is consumed by the volume attached to the deployment pods
 - deploy.sh
  - runs webhook-create-cert and fills the cabundle fetched from kubernetes, then deploys


## Instructions to use

1) run 'docker build -t [name] .'
2) run 'docker tag [image-name] [image store]'
3) run 'docker push [image store]'
4) connect to the cluster and modify the context in the deploy.sh file
5) run the deploy.sh file (requires kubectl set up, openssl command)
