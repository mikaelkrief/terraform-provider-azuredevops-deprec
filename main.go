package main

import (
	"github.com/mikaelkrief/terraform-provider-azuredevops/azuredevops"

	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: azuredevops.Provider})
}
