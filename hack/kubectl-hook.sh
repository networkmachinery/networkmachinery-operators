#!/usr/bin/env bash

# TODO: make service name and ports configurable

BASH_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"


checkPrereqs() {
  command -v cfssl > /dev/null || echo "please install cfssl before proceeding with this script" && exit 1
  command -v cfssljson > /dev/null || echo "please install cfssljson before proceeding with this script" && exit 1
}

createCertAndKey() {

cat <<EOF | cfssl genkey - | cfssljson -bare server
{
  "hosts": [
    "network-validator.default.svc.cluster.local",
    "network-validator.default.svc"
  ],
  "CN": "networkconnectivity-test-controller-f6dc99784-whr4z",
  "key": {
    "algo": "ecdsa",
    "size": 256
  }
}
EOF
}

createCSR() {

cat <<EOF | kubectl apply -f -
apiVersion: certificates.k8s.io/v1beta1
kind: CertificateSigningRequest
metadata:
  name: network-validator.default
spec:
  request: $(cat server.csr | base64 | tr -d '\n')
  usages:
  - digital signature
  - key encipherment
  - server auth
EOF
}

approveCSRAndFetchCert(){
  kubectl certificate approve network-validator.default
  kubectl get csr network-validator.default -o jsonpath='{.status.certificate}' \
      | base64 --decode > server.crt
 }

createSecret(){
 kubectl create secret tls networkconnectivity-layer-validator-secret  --cert server.crt --key server-key.pem
}

createValidationWebhookConfig() {
  caCert=$(kubectl config view --raw -o go-template --template='{{ range .clusters }}{{ index .cluster "certificate-authority-data" }}{{ end }}')

  FILE=validation-webhook-config.tpl.yaml
  sed -e "s/\${CA_BUNDLE}/${caCert}/" "${BASH_DIR}"/templates/"${FILE}" | kubectl apply -f -
}

createWebhookService() {
  FILE=validation-webhook-service.tpl.yaml
  cat "${BASH_DIR}"/templates/"${FILE}" | kubectl apply -f -
}

cleanUP(){
  rm -f server.crt server-key.pem server.csr > /dev/null
}


if [[ "${BASH_SOURCE[0]}" = "$0" ]]; then
  trap cleanUP SIGINT SIGTERM EXIT

  echo "NOTE: You could also add me as a kubectl plugin by putting me in your PATH"
  echo "Configuring webhooks for your cluster"

  createCertAndKey
  createCSR
  approveCSRAndFetchCert
  createSecret
  createValidationWebhookConfig
  createWebhookService
fi

