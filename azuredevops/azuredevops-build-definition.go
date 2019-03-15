package azuredevops

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mikaelkrief/go-azuredevops-sdk/build/5.1-preview"
	"log"
	"strconv"
	"terraform-provider-azuredevops/azuredevops/utils"
	"time"
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
						"branch": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: false,
						},

					},
				},
			},
			"designer_phase": {
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

	if v, ok := repo["branch"]; ok {
		defRepo.DefaultBranch = utils.String(v.(string))
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
	var poolAgent = "Hosted VS2017"

	definition := build.Definition{
		Name:       &name,
		Repository: ExpandRepository(d.Get("repository")),
		DesignerProcess: &build.DesignerProcess{
			Type:   &typeprocess,
			Phases: ExpandPhases(d.Get("designer_phase")),
		},
		Queue: &build.AgentPoolQueue{
			Name: &poolAgent,
			Pool: &build.TaskAgentPoolReference{
				Name: &poolAgent,
			},

		},
	}

	utils.PrettyPrint(definition)

	builddef, err := client.CreateDefinition(ctx, c.organization, definition, project, nil, nil)

	utils.PrettyPrint(builddef) //Log Request

	if err != nil {
		return fmt.Errorf("Error creating build definition  %q: %+v", "", err)
	}

	time.Sleep(100 * time.Millisecond)

	var buildid = builddef.ID
	buildCreated, err2 := client.GetDefinition(ctx, c.organization, project, *buildid, nil, nil, "", nil)

	if err2 != nil {
		return fmt.Errorf("Error getting build definition %v: ", builddef.Name)
	}

	var id = strconv.Itoa(int(int32(*buildCreated.ID)))
	d.SetId(id)

	log.Print(d.Id())

	return resourceBuildDefinitionRead(d, meta)
}

func resourceBuildDefinitionRead(d *schema.ResourceData, meta interface{}) error {

	c := meta.(*AzureDevOpsClient)
	client := meta.(*AzureDevOpsClient).buildClient
	ctx := meta.(*AzureDevOpsClient).StopContext

	var defname = d.Get("name").(string)
	var id = d.Id()
	projectid := d.Get("project_id").(string)
	defid, err := strconv.ParseInt(id, 10, 32)
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
	definitionId, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	builddef, err2 := client.DeleteDefinition(ctx, c.organization, project, int32(definitionId))

	if err2 != nil {
		return fmt.Errorf("Error deleting build definition  %q: %+v", definitionId, err2)
	}
	utils.PrettyPrint(builddef)

	return nil
}
