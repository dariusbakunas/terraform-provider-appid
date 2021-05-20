
<a name="v0.4.2"></a>
## [v0.4.2](https://github.com/dariusbakunas/terraform-provider-appid/compare/v0.4.1...v0.4.2)

> 2021-05-20

### Fix

* add more nil checks


<a name="v0.4.1"></a>
## [v0.4.1](https://github.com/dariusbakunas/terraform-provider-appid/compare/v0.4.0...v0.4.1)

> 2021-05-20

### Fix

* handle empty fields in cd template data source
* move err block
* error block should be before enabling retries


<a name="v0.4.0"></a>
## [v0.4.0](https://github.com/dariusbakunas/terraform-provider-appid/compare/v0.3.0...v0.4.0)

> 2021-05-19

### Feat

* add password regex resource
* add appid password regex data source

### Fix

* make sure resources are re-created on tenant_id change


<a name="v0.3.0"></a>
## [v0.3.0](https://github.com/dariusbakunas/terraform-provider-appid/compare/v0.2.0...v0.3.0)

> 2021-05-18

### Feat

* add google idp resource/datasource
* add facebook idp datasource/resource


<a name="v0.2.0"></a>
## v0.2.0

> 2021-05-17

### Feat

* add releaser configuration
* add redirect_url data source

### Fix

* fix tests
* another ptr fix
* has to be ptr
* handle region/appid_base_url props

### Refactor

* enable api retries
* cleanup old api code
* switch remaining resources to new appid go sdk
* switch role to new appid go sdk
* switch redirect_urls to new appid go sdk
* switch saml idp to new appid go sdk
* switch custom idp to new go sdk
* switch cd idp to new appid go sdk
* switch cloud directory tmpl to new go sdk
* switch application resource to new appid go sdk
* switch token config to new appid go sdk client

### Reverts

* do not set provider defaults from env

### Pull Requests

* Merge pull request [#7](https://github.com/dariusbakunas/terraform-provider-appid/issues/7) from watson-health-development/appid-go-sdk
* Merge pull request [#6](https://github.com/dariusbakunas/terraform-provider-appid/issues/6) from watson-health-development/roles
* Merge pull request [#5](https://github.com/dariusbakunas/terraform-provider-appid/issues/5) from watson-health-development/roles
* Merge pull request [#4](https://github.com/dariusbakunas/terraform-provider-appid/issues/4) from watson-health-development/roles
* Merge pull request [#3](https://github.com/dariusbakunas/terraform-provider-appid/issues/3) from watson-health-development/docs
* Merge pull request [#2](https://github.com/dariusbakunas/terraform-provider-appid/issues/2) from watson-health-development/docs
* Merge pull request [#1](https://github.com/dariusbakunas/terraform-provider-appid/issues/1) from watson-health-development/custom-idp
* Merge pull request [#1](https://github.com/dariusbakunas/terraform-provider-appid/issues/1) from wh-return-to-work/RTW-2630

