[中文版](README_zh.md)
## How to Deploy
### GCP Credentials setup
```shell
gcloud auth application-default login
或者
GOOGLE_CREDENTIALS=service-account-key-xxxx.json
```
### Deploy the resource, replace the local variables with your's (project_id and region)
```shell
cd http-info/
terraform init
terraform apply  --auto-approve

cd hello-world-go/
terraform init
terraform apply  --auto-approve
```
### Delete all resource
```shell
terraform destroy --auto-approve
```