package azuredevops

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	azuredevopssdk "github.com/mikaelkrief/go-azuredevops-sdk"
)

//Provider AzureDevOps
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AZURE_DEVOPS_ORGANIZATION", nil),
				Description: "Azure DevOps organization name",
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AZURE_DEVOPS_TOKEN", nil),
				Description: "Azure DevOps token",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"azuredevops_project": resourceProjectObject(),
		},
		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	organization := d.Get("organization").(string)
	token := d.Get("token").(string)

	return azuredevopssdk.NewClientWith(organization, token)
}
