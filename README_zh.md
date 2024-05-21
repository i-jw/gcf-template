[English Version](README.md)
## 如何部署
### 设置GCP认证
```shell
gcloud auth application-default login
或者
GOOGLE_CREDENTIALS=service-account-key-xxxx.json
```
### 部署，注意替换本地变量中的项目ID和区域
```shell
cd http-info/
terraform init
terraform apply  --auto-approve

cd hello-world-go/
terraform init
terraform apply --auto-approve
```
### 销毁资源
```shell
terraform destroy --auto-approve
```