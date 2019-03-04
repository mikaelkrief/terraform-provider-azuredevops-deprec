package azuredevops

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/mikaelkrief/go-azuredevops-sdk/build/5.1-preview"
	"github.com/mikaelkrief/go-azuredevops-sdk/core/5.1-preview"
	"github.com/mikaelkrief/go-azuredevops-sdk/operation/5.1-preview.1"
	"time"
)

// AzureDevOpsClient contains the handles to all the specific Azure Resource Manager
// resource classes' respective clients.
type AzureDevOpsClient struct {
	organization string
	token        string
	StopContext  context.Context

	buildClient     build.BuildClient
	coreClient      core.CoreClient
	operationClient operation.OperationClient
}

func (c *AzureDevOpsClient) configureClient(client *autorest.Client, auth autorest.Authorizer) {
	client.Authorizer = auth
	client.PollingDuration = 60 * time.Minute

}

// getAzDOClient is a helper method which returns a fully instantiated
// *getAzDOClient based on the Config's current settings.
func getAzDOClient(_organization string, _token string) (*AzureDevOpsClient, error) {

	// client declarations:
	client := AzureDevOpsClient{
		organization: _organization,
		token:        _token,
	}

	auth := autorest.NewBasicAuthorizer(_token)
	auth.WithAuthorization()

	client.registerBuildServiceClients(auth)
	client.registerCoreServiceClients(auth)
	client.registerOperationServiceClients(auth)

	return &client, nil
}

func (c *AzureDevOpsClient) registerBuildServiceClients(auth autorest.Authorizer) {
	buildctl := build.NewBuildClient()
	c.configureClient(&buildctl.Client, auth)
	c.buildClient = buildctl
}

func (c *AzureDevOpsClient) registerCoreServiceClients(auth autorest.Authorizer) {
	corectl := core.NewCoreClient()
	c.configureClient(&corectl.Client, auth)
	c.coreClient = corectl
}

func (c *AzureDevOpsClient) registerOperationServiceClients(auth autorest.Authorizer) {
	operationclt := operation.NewOperationClient()
	c.configureClient(&operationclt.Client, auth)
	c.operationClient = operationclt
}
