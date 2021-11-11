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
`3ff1d6c3631e836a96f03d756f36b230d17c9e24dffd95b1b835011bfd5522b5`

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
`b4f09c06aa97323cd4c51a8130f618c5adb693e0e23bb9dc2b7c808014b87090`

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