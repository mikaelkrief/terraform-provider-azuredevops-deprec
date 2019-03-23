package azuredevops

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"azuredevops": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {

	//for test
	os.Setenv("AZURE_DEVOPS_ORGANIZATION", "MKtests")
	os.Setenv("AZURE_DEVOPS_TOKEN", "kkubaqyoeaachzox3jnhmmgb7w7xxkhzkp4ot6jf7osi7t243pgq")

	if v := os.Getenv("AZURE_DEVOPS_ORGANIZATION"); v == "" {
		t.Fatal("AZURE_DEVOPS_ORGANIZATION must be set for acceptance tests")
	}
	if v := os.Getenv("AZURE_DEVOPS_TOKEN"); v == "" {
		t.Fatal("AZURE_DEVOPS_TOKEN must be set for acceptance tests")
	}
}
