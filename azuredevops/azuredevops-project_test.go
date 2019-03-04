package azuredevops

import (
	"fmt"
	"terraform-provider-azuredevops/azuredevops/utils"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func Test_projectCheck(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testProjectCheckBasicMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testProjectExist("azuredevops_project.test1"),
				),
			},
			{
				Config: testProjectCheckBasicWithTemplate(),
				Check: resource.ComposeTestCheckFunc(
					testProjectExist("azuredevops_project.test2"),
				),
			},
		},
	})
}

func testProjectCheckBasicMinimal() string {
	return fmt.Sprintf(
		`resource "azuredevops_project" "test1" {
			name  = "project-%v"
  }
`, utils.String(5))
}

func testProjectCheckBasicWithTemplate() string {
	return fmt.Sprintf(
		`resource "azuredevops_project" "test2" {
			name               = "project2-%v"
			description = "description test for project"
			template_type_name="scrum"
		  }
`, utils.String(5))
}

func testProjectDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AzureDevOpsClient)
	client := testAccProvider.Meta().(*AzureDevOpsClient).coreClient
	ctx := testAccProvider.Meta().(*AzureDevOpsClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azuredevops_check" {
			continue
		}

		var includecapa = true
		var includehisto = false
		project, err := client.GetProject(ctx, conn.organization, rs.Primary.ID, &includecapa, &includehisto)

		if &project != nil {
			return fmt.Errorf("Bad: Project %q still exists", rs.Primary.ID)
		}

		if err != nil {
			return err
		}

		return nil
	}

	return nil
}

func testProjectExist(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		c := testAccProvider.Meta().(*AzureDevOpsClient)
		conn := testAccProvider.Meta().(*AzureDevOpsClient).coreClient
		ctx := testAccProvider.Meta().(*AzureDevOpsClient).StopContext

		project, err := conn.GetProject(ctx, c.organization, rs.Primary.ID, nil, nil)

		if err != nil {
			return err
		}

		if &project == nil {
			return err
		}
		return nil
	}
}
