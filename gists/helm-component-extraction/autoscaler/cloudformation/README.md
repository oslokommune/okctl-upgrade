## Prepare

1. Export the `AWS_PROFILE` variable with the relevant profile name.
    ```bash
		export AWS_PROFILE=your-profile-name
		```
2. Log in to the AWS account with `aws sso login`.

## Install the Cloudformation stacks

1. Edit the Makefile and configure the user input section.
2. Run `make configure` to build the Cloudformation stacks.
3. Run `make install` to deploy the Cloudformation stacks.
