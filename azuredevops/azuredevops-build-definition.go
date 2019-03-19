package azuredevops

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mikaelkrief/go-azuredevops-sdk/build/5.1-preview"
	"github.com/satori/go.uuid"
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
						"steps": {
							Type:        schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"display_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"enabled":{
										Type:     schema.TypeBool,
										Optional:true,
										Default: true,
									},
									"task_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"inputs" : {
										Type:     schema.TypeMap,
										Required: true,
									},

								},
							},
						},

					},

				},
			},
			"queue": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pool_name": {
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

func FlattenRepository(input *build.Repository) interface{} {

	result := make(map[string]interface{})

	if input == nil {
		log.Printf("[DEBUG] Repository is nil")
		return result
	}
	if input.Name != nil {
		result["name"] = *input.Name
	}
	if input.Type != nil {
		result["type"] = *input.Type
	}
	if input.DefaultBranch != nil {
		result["branch"] = *input.DefaultBranch
	}
	return result
}

func ExpandPoolQueue(input interface{}) *build.AgentPoolQueue {
	configs := input.([]interface{})
	defPool := build.AgentPoolQueue{}
	defAgent := build.TaskAgentPoolReference{}

	if len(configs) == 0 {
		return &defPool
	}

	repo := configs[0].(map[string]interface{})

	if v, ok := repo["pool_name"]; ok {
		defPool.Name = utils.String(v.(string))
		defAgent.Name = utils.String(v.(string))
	}

	defPool.Pool = &defAgent

	return &defPool
}

func FlattenPoolQueue(input *build.AgentPoolQueue) interface{} {

	result := make(map[string]interface{})

	if input == nil {
		log.Printf("[DEBUG] PoolQueue is nil")
		return result
	}
	if input.Name != nil {
		result["pool_name"] = *input.Name
	}
	return result
}

func ExpandPhases(input interface{}) *[]build.Phase {
	configs := input.([]interface{})
	var phases = []build.Phase{}

	if len(configs) == 0 {
		return &phases
	}

	for i := 0; i < len(configs); i++ {
		conf := configs[i].(map[string]interface{})
		phase := new(build.Phase)
		if v, ok := conf["name"]; ok {
			phase.Name = utils.String(v.(string))

		}
		if v, ok := conf["steps"]; ok {
			phase.Steps = ExpandStepsPhase(v.([]interface{}))

		}
		phases = append(phases, *phase)
	}

	return &phases
}

func FlattenPhases(d *schema.ResourceData,input *build.DesignerProcess) []map[string]interface{} {

	results := make([]map[string]interface{}, 0)
	result := make(map[string]interface{})

	for i, v := range *input.Phases {

		if v.Name != nil {
			result["name"] = *v.Name
		}

		if v.Steps != nil {
			result["steps"] = FlattenStepsPhase(d, *v.Steps)
		}
		results = append(results, result)
		fmt.Printf("i %d", i)
	}
	return results
}

func ExpandStepsPhase(input interface{}) *[]build.DefinitionStep {
	configs := input.([]interface{})
	var steps = []build.DefinitionStep{}

	if len(configs) == 0 {
		return &steps
	}

	for i := 0; i < len(configs); i++ {
		conf := configs[i].(map[string]interface{})
		step := new(build.DefinitionStep)
		if v, ok := conf["display_name"]; ok {
			step.DisplayName = utils.String(v.(string))
		}
		if v, ok := conf["inputs"]; ok {
			step.Inputs = expandMaps(v.(map[string]interface{}))
		}

		if v, ok := conf["task_id"]; ok {
			task := build.TaskDefinitionReference{}
			idUUID, err := uuid.FromString(v.(string))
			if err != nil {
				fmt.Errorf("Error convert uuid  %q: %+v", v.(string), err)
			}
			task.ID = &idUUID
			step.Task = &task
		}

		steps = append(steps, *step)
	}
	return &steps
}

func FlattenStepsPhase(d *schema.ResourceData, input []build.DefinitionStep) []map[string]interface{} {

	results := make([]map[string]interface{}, 0)

	for i, v := range input {
		result := make(map[string]interface{})


		if v.Inputs != nil {
			flattenInputs(d, v.Inputs)
		}

		if v.DisplayName != nil {
			result["display_name"] = *v.DisplayName
		}

		if v.Task.ID != nil {
			result["task_id"] = *v.Task.ID
		}
		fmt.Printf("i %d", i)
		results = append(results, result)
	}
	return results
}

func mapValueToString(v interface{}) (string, error) {
	switch value := v.(type) {
	case string:
		return value, nil
	case int:
		return fmt.Sprintf("%d", value), nil
	default:
		return "", fmt.Errorf("unknown map type %T in map value", value)
	}
}

func expandMaps(maps map[string]interface{}) map[string]*string {
	output := make(map[string]*string, len(maps))

	for i, v := range maps {
		//Validate should have ignored this error already
		value, _ := mapValueToString(v)
		output[i] = &value
	}

	return output
}

func flattenInputs(d *schema.ResourceData, mapMap map[string]*string) {

	// If tagsMap is nil, len(tagsMap) will be 0.
	output := make(map[string]interface{}, len(mapMap))

	for i, v := range mapMap {
		output[i] = *v
	}

	d.Set("inputs", output)
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
			Phases: ExpandPhases(d.Get("designer_phase")),
		},
		Queue: ExpandPoolQueue (d.Get("queue")),
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
	d.Set("repository", FlattenRepository (build.Repository))
	d.Set("queue", FlattenPoolQueue (build.Queue))
	d.Set("designer_phase", FlattenPhases(d, build.DesignerProcess))

	return nil
}

func resourceBuildDefinitionUpdate(d *schema.ResourceData, meta interface{}) error {

	c := meta.(*AzureDevOpsClient)
	client := meta.(*AzureDevOpsClient).buildClient
	ctx := meta.(*AzureDevOpsClient).StopContext

	d.Partial(true)

	var project = d.Get("project_id").(string)
	var name = d.Get("name").(string)
	var typeprocess = int32(1)

	definition := build.Definition{
		Name:       &name,
		Repository: ExpandRepository(d.Get("repository")),
		DesignerProcess: &build.DesignerProcess{
			Type:   &typeprocess,
			Phases: ExpandPhases(d.Get("designer_phase")),
		},
		Queue: ExpandPoolQueue (d.Get("queue")),
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
