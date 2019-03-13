package azuredevops

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mikaelkrief/go-azuredevops-sdk/build/5.1-preview"
	"strconv"
	"terraform-provider-azuredevops/azuredevops/utils"
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
			"repository": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: false,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: false,
						},
					},
				},
			},
			"phase": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: false,
						},
					},
				},
			},
		},
	}
}

func ExpandRepository(input interface{}) *build.Repository {
	configs := input.([]interface{})
	defRepo := build.Repository{}

	if len(configs) == 0 {
		return &defRepo
	}

	repo := configs[0].(map[string]interface{})

	if v, ok := repo["name"]; ok {
		defRepo.Name = utils.String(v.(string))
	}
	if v, ok := repo["type"]; ok {
		defRepo.Type = utils.String(v.(string))
	}

	return &defRepo
}

func ExpandPhases(input interface{}) *[]build.Phase {
	configs := input.([]interface{})
	var phases = []build.Phase{}

	if len(configs) == 0 {
		return &phases
	}

	repo := configs[0].(map[string]interface{})
	phase := new(build.Phase)
	if v, ok := repo["name"]; ok {
		phase.Name = utils.String(v.(string))
	}
	phases = append(phases, *phase)

	return &phases
}

func resourceBuildDefinitionCreate(d *schema.ResourceData, meta interface{}) error {

	c := meta.(*AzureDevOpsClient)
	client := meta.(*AzureDevOpsClient).buildClient
	ctx := meta.(*AzureDevOpsClient).StopContext

	var project = d.Get("project_id").(string)
	var name = d.Get("name").(string)
	var typeprocess = int32(1)

	definition := build.Definition{
		Name:       &name,
		Repository: ExpandRepository(d.Get("repository")),
		DesignerProcess: &build.DesignerProcess{
			Type:   &typeprocess,
			Phases: ExpandPhases(d.Get("phase")),
		},
	}

	utils.PrettyPrint(definition)

	builddef, err := client.CreateDefinition(ctx, c.organization, definition, project, nil, nil)

	utils.PrettyPrint(builddef) //Log Request

	if err != nil {
		return fmt.Errorf("Error creating build definition  %q: %+v", "", err)
	}

	//var buildid = builddef.ID
	//buildCreated, err2 := client.GetDefinition(ctx, c.organization, project, buildid, nil, nil, "", nil)

	//if err2 != nil {
	//	return fmt.Errorf("Error getting build definition %v: ", builddef.Name)
	//}
	var iddefinition = builddef.ID

	d.SetId(string(*iddefinition))


	return nil
}

func resourceBuildDefinitionRead(d *schema.ResourceData, meta interface{}) error {

	c := meta.(*AzureDevOpsClient)
	client := meta.(*AzureDevOpsClient).buildClient
	ctx := meta.(*AzureDevOpsClient).StopContext

	var defname = d.Get("name").(string)

	projectid := d.Get("project_id").(string)
	defid, err := strconv.ParseInt(d.Id(), 10, 32)
	build, err := client.GetDefinition(ctx, c.organization, projectid, int32(defid), nil, nil, "", nil)

	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error getting build definition  %q: %+v", defname, err)
	}

	d.Set("name", build.Name)

	return nil
}

func resourceBuildDefinitionUpdate(d *schema.ResourceData, meta interface{}) error {

	c := meta.(*AzureDevOpsClient)
	client := meta.(*AzureDevOpsClient).buildClient
	ctx := meta.(*AzureDevOpsClient).StopContext

	d.Partial(true)

	var project = d.Get("project_id").(string)
	var name = d.Get("name").(string)

	definition := build.Definition{
		Name: &name,
	}

	definitionId, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return err
	}
	builddef, err := client.UpdateDefinition(ctx, c.organization, definition, project, int32(definitionId), nil, nil)

	if err != nil {
		return fmt.Errorf("Error updating build definition  %q: %+v", name, err)
	}

	utils.PrettyPrint(builddef)
	d.SetId(d.Id())

	if err != nil {
		return err
	}
	d.Partial(false)
	return nil
}

func resourceBuildDefinitionDelete(d *schema.ResourceData, meta interface{}) error {

	c := meta.(*AzureDevOpsClient)
	client := meta.(*AzureDevOpsClient).buildClient
	ctx := meta.(*AzureDevOpsClient).StopContext

	var project = d.Get("project_id").(string)

	definitionId, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return err
	}

	builddef, err := client.DeleteDefinition(ctx, c.organization, project, int32(definitionId))

	if err != nil {
		return fmt.Errorf("Error deleting build definition  %q: %+v", definitionId, err)
	}
	utils.PrettyPrint(builddef)

	return nil
}
