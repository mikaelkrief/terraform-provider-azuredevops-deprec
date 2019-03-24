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

//noinspection GoInvalidCompositeLiteral
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
			"buildnumber_format": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "$(date:yyyyMMdd)$(rev:.r)",
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
			"variables": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"variable": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
									"allow_override": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"is_secret": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
								},
							},
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
						"step": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"display_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"task_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"task_version": {
										Type:     schema.TypeString,
										Required: true,
									},
									"task_type": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "task",
									},
									"inputs": {
										Type:     schema.TypeMap,
										Required: true,
									},
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"always_run": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"continue_on_error": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"condition": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "succeeded()",
										/*ValidateFunc: validation.StringInSlice([]string{
											"succeededOrFailed()",
											"succeeded()",
											"always()",
											"failed()",
											OR ANY string for custom
										}, true),*/
									},
									"timeout_in_minutes": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  0,
									},
									"reference_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"environment_variables": {
										Type:     schema.TypeMap,
										Optional: true,
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
		if v, ok := conf["step"]; ok {
			phase.Steps = ExpandStepsPhase(v.([]interface{}))

		}
		phases = append(phases, *phase)
	}

	return &phases
}

func FlattenPhases(d *schema.ResourceData, input *build.Process) []map[string]interface{} {

	results := make([]map[string]interface{}, 0)
	result := make(map[string]interface{})

	for i, v := range *input.Phases {

		if v.Name != nil {
			result["name"] = *v.Name
		}

		if v.Steps != nil {
			result["step"] = FlattenStepsPhase(d, *v.Steps)
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
		if v, ok := conf["enabled"]; ok {
			step.Enabled = utils.Bool(v.(bool))
		}

		if v, ok := conf["always_run"]; ok {
			step.AlwaysRun = utils.Bool(v.(bool))
		}

		if v, ok := conf["continue_on_error"]; ok {
			step.ContinueOnError = utils.Bool(v.(bool))
		}

		if v, ok := conf["condition"]; ok {
			step.Condition = utils.String(v.(string))
		}

		if v, ok := conf["timeout_in_minutes"].(int); ok {
			step.TimeoutInMinutes = utils.Int32(int32(v))
		}

		if v, ok := conf["reference_name"]; ok {
			step.RefName = utils.String(v.(string))
		}

		if v, ok := conf["environment_variables"]; ok {
			step.Environment = expandMaps(v.(map[string]interface{}))
		}

		if v, ok := conf["task_id"]; ok {
			task := build.TaskDefinitionReference{}
			idUUID, err := uuid.FromString(v.(string))
			if err != nil {
				fmt.Errorf("Error convert uuid  %q: %+v", v.(string), err)
			}
			task.ID = &idUUID
			if v, ok := conf["task_version"]; ok {
				task.VersionSpec = utils.String(v.(string))
			}
			if v, ok := conf["task_type"]; ok {
				task.DefinitionType = utils.String(v.(string))
			}

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
			flattenMaps(d, v.Inputs, "inputs")
		}

		if v.DisplayName != nil {
			result["display_name"] = *v.DisplayName
		}

		if v.Enabled != nil {
			result["enabled"] = *v.Enabled
		}

		if v.AlwaysRun != nil {
			result["always_run"] = *v.AlwaysRun
		}

		if v.ContinueOnError != nil {
			result["continue_on_error"] = *v.ContinueOnError
		}

		if v.Condition != nil {
			result["condition"] = *v.Condition
		}

		if v.TimeoutInMinutes != nil {
			result["timeout_in_minutes"] = int(*v.TimeoutInMinutes)
		}

		if v.RefName != nil {
			result["reference_name"] = *v.RefName
		}

		if v.Environment != nil {
			flattenMaps(d, v.Inputs, "environment_variables")
		}

		if v.Task.ID != nil {
			result["task_id"] = *v.Task.ID
		}
		if v.Task.VersionSpec != nil {
			result["task_version"] = *v.Task.VersionSpec
		}
		if v.Task.DefinitionType != nil {
			result["task_type"] = *v.Task.DefinitionType
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

func flattenMaps(d *schema.ResourceData, mapMap map[string]*string, schemaKey string) {

	// If tagsMap is nil, len(tagsMap) will be 0.
	output := make(map[string]interface{}, len(mapMap))

	for i, v := range mapMap {
		output[i] = *v
	}

	d.Set(schemaKey, output)
}

/*func expandVariables(d *schema.ResourceData) (map[string]*build.DefinitionVariable, error) {
	vars := d.Get("variable").(*schema.Set).List()
	//varstab := vars.(map[string]interface{})
	variables := make(map[string]*build.DefinitionVariable, len(vars))

	for _, sgVar := range vars {

		sgVar := sgVar.(map[string]interface{})
		name := sgVar["name"].(string)
		value := sgVar["value"].(string)
		allowOverride := sgVar["allow_override"].(bool)
		isSecret := sgVar["is_secret"].(bool)
		variable := build.DefinitionVariable{
			Value:         &value,
			AllowOverride: &allowOverride,
			IsSecret:      &isSecret,
		}
		valuemap := variable

		variables[name] = &valuemap
	}

	return variables, nil
}*/

func expandVariables(input interface{}) map[string]*build.DefinitionVariable {
	configs := input.([]interface{})

	variables := configs[0].(map[string]interface{})

	if v, ok := variables["variable"]; ok {
		vars := v.([]interface{})
		_variables := make(map[string]*build.DefinitionVariable, len(vars))

		for _, sgVar := range vars {

			sgVar := sgVar.(map[string]interface{})
			name := sgVar["name"].(string)
			value := sgVar["value"].(string)
			allowOverride := sgVar["allow_override"].(bool)
			isSecret := sgVar["is_secret"].(bool)
			variable := build.DefinitionVariable{
				Value:         &value,
				AllowOverride: &allowOverride,
				IsSecret:      &isSecret,
			}
			valuemap := variable

			_variables[name] = &valuemap
		}

		return _variables

	}
	return nil
}

func FlattenVariables(d *schema.ResourceData, input map[string]*build.DefinitionVariable) map[string]interface{} {
	resultvar := make(map[string]interface{})
	results := make([]map[string]interface{}, 0)

	for key, v := range input {
		result := make(map[string]interface{})

		if v.IsSecret != nil {
			result["is_secret"] = *v.IsSecret
		}

		if v.AllowOverride != nil {
			result["allow_override"] = *v.AllowOverride
		}

		if v.Value != nil {
			result["value"] = *v.Value
		}

		result["name"] = key

		results = append(results, result)

	}
	resultvar["variables"] = results

	return resultvar

}

func resourceBuildDefinitionCreate(d *schema.ResourceData, meta interface{}) error {

	c := meta.(*AzureDevOpsClient)
	client := meta.(*AzureDevOpsClient).buildClient
	ctx := meta.(*AzureDevOpsClient).StopContext

	var project = d.Get("project_id").(string)
	var name = d.Get("name").(string)
	var buildnumberformat = d.Get("buildnumber_format").(string)
	var typeprocess = int32(1)
	buildVars := expandVariables(d.Get("variables"))
	/*if varsErr != nil {
		return fmt.Errorf("Error Building list of Variables: %+v", varsErr)
	}*/

	definition := build.Definition{
		Name:              &name,
		BuildNumberFormat: &buildnumberformat,
		Repository:        ExpandRepository(d.Get("repository")),
		Process: &build.Process{
			Type:   &typeprocess,
			Phases: ExpandPhases(d.Get("designer_phase")),
		},
		Queue:     ExpandPoolQueue(d.Get("queue")),
		Variables: buildVars,
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
	d.Set("repository", FlattenRepository(build.Repository))
	d.Set("queue", FlattenPoolQueue(build.Queue))
	d.Set("designer_phase", FlattenPhases(d, build.Process))
	d.Set("buildnumber_format", build.BuildNumberFormat)
	d.Set("variables", FlattenVariables(d, build.Variables))

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
	var buildnumberformat = d.Get("buildnumber_format").(string)
	definitionId, err := strconv.ParseInt(d.Id(), 10, 32)
	var defid = int32(definitionId)

	listRev, err := client.GetDefinitionRevisions(ctx, c.organization, project, int32(definitionId))
	var countRev = len(*listRev.Value)
	var newRevision = int32(countRev)
	buildVars := expandVariables(d.Get("variables"))

	definition := build.Definition{
		ID:                &defid,
		Name:              &name,
		BuildNumberFormat: &buildnumberformat,
		Repository:        ExpandRepository(d.Get("repository")),
		Process: &build.Process{
			Type:   &typeprocess,
			Phases: ExpandPhases(d.Get("designer_phase")),
		},
		Queue:     ExpandPoolQueue(d.Get("queue")),
		Variables: buildVars,
		Revision:  &newRevision,
	}

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
