# Azure DevOps Go SDK

This code contain the Go SDK Client for Azure DevOps Rest API
[documentation](https://docs.microsoft.com/en-us/rest/api/azure/devops/?view=azure-devops-rest-5.0)


go-azuredevops-sdk provides Go packages for managing and using Azure DevOps services.
It officially supports the last two major releases of Go.  Older versions of
Go will be kept running in CI until they no longer work due to changes in any
of the SDK's external dependencies.  The CHANGELOG will be updated when a
version of Go is removed from CI.


# Install and Use:

## Install

```sh
$ go get -u github.com/mikaelkrief/go-azuredevops-sdk/...
```

or if you use dep, within your repo run:

```sh
$ dep ensure -add github.com/mikaelkrief/go-azuredevops-sdk/
```

If you need to install Go, follow [the official instructions](https://golang.org/dl/).


## Implemented Resources :

- Project
- Process template
- Operation

