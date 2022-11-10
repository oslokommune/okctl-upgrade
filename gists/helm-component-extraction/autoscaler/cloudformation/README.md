## Prepare

1. Export the `AWS_PROFILE` variable with the relevant profile name.
    ```bash
		export AWS_PROFILE=your-profile-name
		```
2. Log in to the AWS account with `aws sso login`

## Install the Cloudformation stacks

2. Edit the Makefile and set the `AWS_PROFILE` variable to the relevant profile name.
3. Run `make configure` to build the Cloudformation stacks.
4. Run `make install` to deploy the Cloudformation stacks.
