# Okctl remote state upgrade

## Description

This upgrade moves the `state.db` file from a local machine to an S3 bucket.

## Installation

### MacOs
Download [remote-state-upgrader_v0.0.0_Darwin_amd64.tar.gz](dist/remote-state-upgrader_v0.0.0_Darwin_amd64.tar.gz) to
your machine and run `tar xvf remote-state-upgrader_v0.0.0_Darwin_amd64.tar.gz`

Run the following and verify that the checksum matches `52d9434c30638f47b482326b3ea8319758dc7a2cabc5930954334bee776f88c1`
```shell
sha256sum remote-state-upgrader
```


### Linux 
Download [remote-state-upgrader_v0.0.0_Linux_amd64.tar.gz](dist/remote-state-upgrader_v0.0.0_Linux_amd64.tar.gz) to
your machine and run `tar xvf remote-state-upgrader_v0.0.0_Darwin_amd64.tar.gz`

Run the following and verify that the checksum matches `3c75073fcc620934f35f44d9f0b582ce9656fd75a8e469a840c6eb68cf9b3a4d`
```shell
sha256sum remote-state-upgrader
```

## Usage

1. Authenticate with your desired okctl environment by either 
   1. Running `okctl venv -c <cluster manifest>` (preferred)
   2. Running `saml2aws login`, see [saml2aws](https://github.com/Versent/saml2aws)
2. Run `./remote-state-upgrade <path to cluster manifest>`
 