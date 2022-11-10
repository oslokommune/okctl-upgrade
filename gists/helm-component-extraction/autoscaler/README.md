# Autoscaler

## Prepare

1. Export the `AWS_PROFILE` variable with the relevant profile name.
    ```bash
		export AWS_PROFILE=your-profile-name
		```
2. Log in to the AWS account with `aws sso login`

## Installation

### Prepare Cloudformation stacks

See [cloudformation/README.md](cloudformation/README.md) for details.

### Install the autoscaler

1. Edit the Makefile and set the variable(s) listed in the user input section
2. Run `make values` to generate the values.yaml file
3. Run `make install` to install the autoscaler
