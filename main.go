package main

import (
	"github.com/hashicorp/terraform/plugin"
	"./azuredevops"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: azuredevops.Provider})
}