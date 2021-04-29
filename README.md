<a href="https://terraform.io">
    <img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" alt="Terraform logo" title="Terraform" align="right" height="50" />
</a>

# Terraform Provider for AppID

Check `examples` folder for working examples

## Development

### Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.13.6+ (to run acceptance tests)
- [Go](https://golang.org/doc/install) 1.15.8+ (to build the provider plugin)

### Quick Start

First, clone the repository:

```bash
git clone git@github.ibm.com:wh-return-to-work/terraform-provider-appid.git
```

To compile the provider, run make build. This will build the provider and put the provider binary in the current directory.

```bash
make build
```

Run `make install` to install provider binary under `~/.terraform.d/plugins/ibm.com/watson-health/appid/{VERSION}/{OS_ARCH}`.

After it is installed, terraform should be able to detect it during `terraform init` phase.

### Testing

To run unit tests, simply run:

```bash
make test
```

To run acceptance tests, make sure `IAM_API_KEY` environment variable is set and execute:

```bash
make testacc
```

**Note:** Acceptance tests create/destroy real resources, while they are named using `tf_acc_test` testing prefix, use some caution. Check `provider_test.go` contents for supported environment variables and their default values.

### Documentation

#### Environment setup

1. Install [pipenv](https://pipenv.readthedocs.io/en/latest/#install-pipenv-today)

        % brew install pipenv
        
2. Install [pyenv](https://github.com/pyenv/pyenv#installation)

        % brew install pyenv
        
3. Use `pyenv` to install `Python 3` if not installed already (run `pyenv versions` to check installed versions or `pyenv list` to list available for install):

        % pyenv install 3.8.1

4. Run `pipenv install` to install `mkdocs` dependencies

5. In order to activate the virtual environment associated with this project you can simply use the shell keyword:

        % pipenv shell

6. To generate or update Terraform documentation, run `go generate`.

7. To serve mk-docs locally, run `mkdocs serve`