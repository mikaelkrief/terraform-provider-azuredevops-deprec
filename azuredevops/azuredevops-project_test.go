package azuredevops

import (
	"fmt"
	azuredevopssdk "go-azuredevops-sdk"
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
				Config: testProjectCheck_Basic,
				Check: resource.ComposeTestCheckFunc(
					testProjectExist("azuredevops_project.test"),
				),
			},
		},
	})
}

const testProjectCheck_Basic = `
resource "azuredevops_project" "test" {
	name               = "test Terraform 2"
	template_type_name = "agile1"
  }
`

func testProjectDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*azuredevopssdk.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azuredevops_check" {
			continue
		}

		project, err := conn.GetProject(rs.Primary.ID)

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

		conn := testAccProvider.Meta().(*azuredevopssdk.Client)

		project, err := conn.GetProject(rs.Primary.ID)

		if err != nil {
			return err
		}

		if &project == nil {
			return err
		}

		return nil

	}
}
