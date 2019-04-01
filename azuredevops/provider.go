package azuredevops

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

//Provider AzureDevOps
func Provider() terraform.ResourceProvider {
	p := &schema.Provider{
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
			"azuredevops_project":          resourceProjectObject(),
			"azuredevops_build_definition": resourceBuildDefinitionObject(),
		},
	}

	p.ConfigureFunc = providerConfigure(p)

	return p
}

func providerConfigure(p *schema.Provider) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		organization := d.Get("organization").(string)
		token := d.Get("token").(string)

		client, err := getAzDOClient(organization, token)
		if err != nil {
			return nil, err
		}

		client.StopContext = p.StopContext()

		// replaces the context between tests
		p.MetaReset = func() error {
			client.StopContext = p.StopContext()
			return nil
		}

		return client, nil

	}

}
