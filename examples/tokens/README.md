# AppID Config Tokens Datasource Example

Make sure `IAM_API_KEY`, `APPID_BASE_URL` and `IAM_BASE_URL` environment variables are set, then initialize:

```bash
terraform init
```

And execute:

```bash
TF_LOG_PROVIDER=INFO terraform apply -var-file="example.tfvars"
```