# Okctl remote state upgrade

## Description

This upgrade moves the `state.db` file from a local machine to an S3 bucket.

## Installation

### MacOs
Download [remote-state-upgrader_v0.0.0_Darwin_amd64.tar.gz](dist/remote-state-upgrader_v0.0.0_Darwin_amd64.tar.gz) to
your machine and run 
```shell
tar xvf remote-state-upgrader_v0.0.0_Darwin_amd64.tar.gz
```

Run the following and verify that the checksum matches 
`d861a1e258ded5ef27c66ca9ab6f86adc9ceeb93e4bda0793d08decf931a4f48`

```shell
sha256sum remote-state-upgrader
```


### Linux 
Download [remote-state-upgrader_v0.0.0_Linux_amd64.tar.gz](dist/remote-state-upgrader_v0.0.0_Linux_amd64.tar.gz) to
your machine and run 

```shell
tar xvf remote-state-upgrader_v0.0.0_Darwin_amd64.tar.gz
```

Run the following and verify that the checksum matches 
`8b38434a088a59b08bd52e40fdcac1f0b24c4c70adab6fae0f9eb947e859a0b0`

```shell
sha256sum remote-state-upgrader
```

## Usage

1. Authenticate with your desired okctl environment by either 
   1. Running `okctl venv -c <cluster manifest>` (preferred)
   2. Running `saml2aws login`, see [saml2aws](https://github.com/Versent/saml2aws)
2. Run `./remote-state-upgrade <path to cluster manifest>`
 
## Developers

After making changes to the code:
1. Run `make release`
2. Make sure to commit the `*.tar.gz` files in `dist/`
3. Update the checksums in this README by running `sha256sum` on the binaries produced by `make release` residing in
   1. `/dist/linux_linux_amd64`
   2. `/dist/darwin_darwin_amd64`

When making changes to the code, remember to update the checksums here in the readme after running `make release`. 
Remember to check in the tar.gz files in dist