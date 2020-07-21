#!/bin/bash

#Set kubectl for windows
kubectl="kubectl.exe"

#Set context
${kubectl} config set-context webhook --namespace=webhook --cluster=minikube --user=minikube
${kubectl} config use-context webhook

# Generate TLS secret
./webhook-create-cert.sh --namespace webhook --service webhook-server --secret node-lifetime-webhook-certs

# Read the PEM-encoded CA certificate, base64 encode it, and replace the `${CA_PEM_B64}` placeholder in the YAML
# template with it. Then, create the Kubernetes resources.
echo "Deploying ..."
ca_pem_b64="$(${kubectl} config view --raw --minify --flatten -o jsonpath='{.clusters[].cluster.certificate-authority-data}')"
sed -e 's@${CA_PEM_B64}@'"$ca_pem_b64"'@g' <"deployments/deployment.yaml" \
    | ${kubectl} apply -f -

echo "Deploy Complete"
