# Okctl remote state upgrade

## Description

This upgrade moves the `state.db` file from a local machine to an S3 bucket.

## Installation

### MacOs
Download [remote-state-upgrader_v0.0.0_Darwin_amd64.tar.gz](dist/remote-state-upgrader_v0.0.0_Darwin_amd64.tar.gz) to
your machine and run `tar xvf remote-state-upgrader_v0.0.0_Darwin_amd64.tar.gz`

Run the following and verify that the checksum matches `b17e217abaf09f87f535d382624f78cb4f11a9f91aba14bb7f63d06b41a27f85`
```shell
sha256sum remote-state-upgrader
```


### Linux 
Download [remote-state-upgrader_v0.0.0_Linux_amd64.tar.gz](dist/remote-state-upgrader_v0.0.0_Linux_amd64.tar.gz) to
your machine and run `tar xvf remote-state-upgrader_v0.0.0_Darwin_amd64.tar.gz`

Run the following and verify that the checksum matches `d58f57d388074d7511537c112deb5c6908f5f4c014a158b04fa3ea722b4b288a`
```shell
sha256sum remote-state-upgrader
```

## Usage

1. Authenticate with your desired okctl environment by either 
   1. Running `okctl venv -c <cluster manifest>` (preferred)
   2. Running `saml2aws login`, see [saml2aws](https://github.com/Versent/saml2aws)
2. Run `./remote-state-upgrade <path to cluster manifest>`
 