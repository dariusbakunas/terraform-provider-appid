module github.ibm.com/dbakuna/terraform-provider-appid

go 1.15

// TEMPORARY: until it is merged into official IBM repository
replace github.com/IBM/appid-go-sdk => github.com/dariusbakunas/appid-go-sdk v0.0.0-20210517152640-71d307a222c8

require (
	cloud.google.com/go/storage v1.15.0 // indirect
	github.com/IBM/appid-go-sdk v0.0.0-20210510121547-84997d5df141
	github.com/IBM/go-sdk-core v1.1.0
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/apparentlymart/go-cidr v1.1.0 // indirect
	github.com/aws/aws-sdk-go v1.38.40 // indirect
	github.com/fatih/color v1.11.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/go-querystring v1.1.0
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-hclog v0.16.1 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-plugin v1.4.1 // indirect
	github.com/hashicorp/go-uuid v1.0.2 // indirect
	github.com/hashicorp/hcl/v2 v2.10.0 // indirect
	github.com/hashicorp/terraform-json v0.11.0 // indirect
	github.com/hashicorp/terraform-plugin-docs v0.4.0
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.6.1
	github.com/hashicorp/yamux v0.0.0-20210316155119-a95892c5f864 // indirect
	github.com/klauspost/compress v1.12.2 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/oklog/run v1.1.0 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/ulikunitz/xz v0.5.10 // indirect
	github.com/zclconf/go-cty v1.8.3 // indirect
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a // indirect
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	golang.org/x/net v0.0.0-20210510120150-4163338589ed // indirect
	golang.org/x/oauth2 v0.0.0-20210514164344-f6687ab2804c
	golang.org/x/sys v0.0.0-20210514084401-e8d321eab015 // indirect
	golang.org/x/tools v0.1.1 // indirect
	google.golang.org/api v0.46.0 // indirect
	google.golang.org/genproto v0.0.0-20210513213006-bf773b8c8384 // indirect
	google.golang.org/grpc v1.37.1 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
