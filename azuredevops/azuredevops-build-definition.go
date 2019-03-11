package azuredevops

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBuildDefinitionObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceBuildDefinitionCreate,
		Update: resourceBuildDefinitionUpdate,
		Read:   resourceBuildDefinitionRead,
		Delete: resourceBuildDefinitionDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}


func resourceBuildDefinitionCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}


func resourceBuildDefinitionRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceBuildDefinitionUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceBuildDefinitionDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}