Azure DevOps Terraform Provider
===============================

Terraform provider for [Azure DevOps](https://azure.microsoft.com/en-us/services/devops/)

General Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.11.x (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/mikaelkrief/terraform-provider-azuredevops`

```sh
$ mkdir -p $GOPATH/src/github.com/mikaelkrief; cd $GOPATH/src/github.com/mikaelkrief
$ git clone git@github.com:mikaelkrief/terraform-provider-azuredevops
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/mikaelkrief/terraform-provider-azuredevops
$ sh "scripts/install.sh" <version>
```

Requirements
---------------------
- Have or create an Azure DevOps organization (account), see the [documentation](https://docs.microsoft.com/en-us/azure/devops/organizations/accounts/create-organization?view=azure-devops) for create it for free
- Generate a PAT (Personnal Access Token) for authentication , see the [documentation](https://docs.microsoft.com/en-us/azure/devops/organizations/accounts/use-personal-access-tokens-to-authenticate?view=azure-devops)

Using the provider
----------------------

```
# Configure the Azure DevOps provider
provider "azuredevops" {
  organization = "<name of your Azure DevOps organisation"
  token        = "<Your PAT authentication>"
}

# Create a Project
resource "azuredevops_project" "project" {
  name               = "My project"
  template_type_name = "agile"
  description = "description of my project"
  source_control_type = "Git"
}
```
