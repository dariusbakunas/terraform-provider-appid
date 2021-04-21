# AppID Config Tokens Datasource Example

Make sure `TF_VAR_iam_access_token` environment variable is set, then initialize:

```bash
terraform init
```

And execute:

```bash
TF_LOG_PROVIDER=INFO terraform apply -var-file="example.tfvars"
```