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

Run the following command to verify the checksum:

```shell
sha256sum remote-state-upgrader
```

Expected checksum: `199411f198d45a864c7b1123bf098d9752568d0f1641f8510bd249ecab68eba2`


### Linux
Download [remote-state-upgrader_v0.0.0_Linux_amd64.tar.gz](dist/remote-state-upgrader_v0.0.0_Linux_amd64.tar.gz) to
your machine and run

```shell
tar xvf remote-state-upgrader_v0.0.0_Darwin_amd64.tar.gz
```

Run the following command to verify the checksum:

```shell
sha256sum remote-state-upgrader
```

Expected checksum: `7ff544fe30a0c231342087217c8f05437326eebdf98ee1ce7a2cff61ad0e5d7f`

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
