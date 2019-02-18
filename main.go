package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/mikaelkrief/terraform-provider-azuredevops/azuredevops"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: azuredevops.Provider})
}
