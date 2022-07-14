# URL Signer
REST Server to sign Google Cloud Media CDN URLs written in Golang. Mostly http handler around sample in https://gist.github.com/mlevkov/8d1a481992494210cb2e5cc3a1c05221

May require recent version of Go to compile

## Getting Started

Before you get started, verify that your project is allow-listed for Media CDN services.

Clone this repository.

```
git clone https://github.com/alekssaul/urlsigner 
cd urlsigner
```

## Generate ed25519 certs

Run the certgen utility to generate public and private ed25519 certificates in URL safe Base64 encoded format

```sh
go run ./certgen
```

Move the certificate files into terraform assets folder

```sh
mv *.key ./deploy/terraform/assets/
```

## Deploy Media CDN Infrastructure

`deploy/terraform` folder contains terraform specs to bootstrap a test infrastructure for Media CDN.

Run terraform commands to initialize the terraform plugins.

```sh
cd deploy/terraform
terraform init
```

Run terraform plan to validate infrastructure changes
```sh
terraform plan
```

Deploy the Media CDN settings
```sh
terrform apply
```

## Deploy URL Signer service to Cloud Run

Set the KEYSET and PRIVATEKEY environmental variable to output of terraform

```sh
export KEYSET=$(terraform output --raw keyset)
export PRIVATEKEY=$(terraform output --raw keyset_primary_private)
```

Deploy the service to cloud run

```sh
gcloud run deploy --set-env-vars=KEYSET=$KEYSET  --update-secrets=PRIVATEKEY=$PRIVATEKEY
```

