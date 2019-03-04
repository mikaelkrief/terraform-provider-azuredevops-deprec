package azuredevops

import (
	"fmt"
	"github.com/mikaelkrief/go-azuredevops-sdk/core/5.1-preview"
	"github.com/satori/go.uuid"
	"log"
	"strings"
	"terraform-provider-azuredevops/azuredevops/utils"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceProjectObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Update: resourceProjectUpdate,
		Read:   resourceProjectRead,
		Delete: resourceProjectDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"source_control_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "Git",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val
					if v != "Git" && v != "TFVC" {
						errs = append(errs, fmt.Errorf("%q must be Git or TFVC, got: %q", key, v))
					}
					return
				},
			},
			"template_type_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
			"visibility": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "private",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if strings.ToLower(old) == strings.ToLower(new) {
						return true
					}
					return false
				},
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := strings.ToLower(val.(string))
					if v != "private" && v != "public" {
						errs = append(errs, fmt.Errorf("%q must be private or public, got: %q", key, v))
					}
					return
				},
			},
		},
	}
}

func resourceProjectCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*AzureDevOpsClient)
	client := meta.(*AzureDevOpsClient).coreClient
	clientoperation := meta.(*AzureDevOpsClient).operationClient
	ctx := meta.(*AzureDevOpsClient).StopContext

	var name = d.Get("name").(string)
	var description = d.Get("description").(string)
	var sourcecontroltype = d.Get("source_control_type").(string)

	process, err := getProcessByName(meta, d.Get("template_type_name").(string))
	if err != nil {
		return err
	}

	var processid = process.ID.String()
	teamProject := core.TeamProject{
		Name:        &name,
		Description: &description,
		Capabilities: map[string]map[string]*string{
			"processTemplate": {
				"templateTypeId": &processid,
			},
			"versioncontrol": {
				"sourceControlType": &sourcecontroltype,
			},
		},
	}

	if d.Get("template_type_name").(string) == "" {
		d.Set("template_type_name", process.Name)
	}

	operation, err := client.QueueCreateProject(ctx, c.organization, teamProject)

	utils.PrettyPrint(teamProject) //Log Request

	if err != nil {
		return fmt.Errorf("Error creating project  %q: %+v", *teamProject.Name, err)
	}

loop:
	for t := range time.NewTicker(10 * time.Second).C {
		utils.PeriodicFunc(t)
		op, err := clientoperation.GetOperation(ctx, *operation.ID, c.organization, nil)
		if err != nil {
			return fmt.Errorf("Error getting project creating operation  %q: %+v", op.ID, err)
		}
		if op.Status == "succeeded" {
			break loop
		}
	}

	projectCreated, err := client.GetProject(ctx, c.organization, *teamProject.Name, nil, nil)
	if err != nil {
		return fmt.Errorf("Error getting project %v: ", teamProject.Name)
	}
	var idproject = projectCreated.ID

	d.SetId(idproject.String())

	if err != nil {
		return err
	}
	return nil
}

func getProcessByName(meta interface{}, processname string) (process core.Process, err error) {
	c := meta.(*AzureDevOpsClient)
	client := c.coreClient
	ctx := c.StopContext
	var processToApply = core.Process{}

	if processname != "" {
		var ProcessName = processname
		process, err := client.GetProcessIdbyName(ctx, c.organization, ProcessName)
		if err != nil {
			return processToApply, fmt.Errorf("error get process template %q: %+v", ProcessName, err)
		}

		utils.PrettyPrint(process) //Log Request
		log.Printf(*process.Name)
		processToApply = *process
		return processToApply, nil
	} else {
		process, err := client.GetDefaultProcess(ctx, c.organization)
		if err != nil {
			return *process, fmt.Errorf("error get default process template %+v", err)
		}
		return *process, nil
	}

}

func resourceProjectRead(d *schema.ResourceData, meta interface{}) error {

	c := meta.(*AzureDevOpsClient)
	client := meta.(*AzureDevOpsClient).coreClient
	ctx := meta.(*AzureDevOpsClient).StopContext

	var projectname = d.Get("name").(string)
	var includecapa = true
	var includehisto = false
	project, err := client.GetProject(ctx, c.organization, projectname, &includecapa, &includehisto)
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error getting project  %q: %+v", projectname, err)
	}

	d.Set("name", project.Name)
	d.Set("description", project.Description)
	d.Set("visibility", project.Visibility)
	d.Set("source_control_type", project.Capabilities["versioncontrol"]["sourceControlType"])

	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, meta interface{}) error {

	c := meta.(*AzureDevOpsClient)
	client := meta.(*AzureDevOpsClient).coreClient
	clientoperation := meta.(*AzureDevOpsClient).operationClient
	ctx := meta.(*AzureDevOpsClient).StopContext

	d.Partial(true)

	if d.HasChange("source_control_type") {
		return fmt.Errorf("The source controle type %q can't be changed", d.Get("source_control_type").(string))
	}

	if d.HasChange("template_type_name") {
		return fmt.Errorf("the template type Id can't be changed")
	}

	var name = d.Get("name").(string)
	var description = d.Get("description").(string)
	teamProject := core.TeamProject{}

	if d.HasChange("name") {
		teamProject.Name = &name
	}
	if d.HasChange("description") {
		teamProject.Description = &description
	}

	var projectid = d.Id()
	idUUID, err := uuid.FromString(projectid)
	op, err := client.UpdateProject(ctx, c.organization, teamProject, idUUID)

	if err != nil {
		return fmt.Errorf("Error updating project  %q: %+v", name, err)
	}

loop:
	for t := range time.NewTicker(10 * time.Second).C {
		utils.PeriodicFunc(t)
		op, err := clientoperation.GetOperation(ctx, *op.ID, c.organization, nil)
		if err != nil {
			return fmt.Errorf("Error getting project creating operation  %q: %+v", op.ID, err)
		}
		if op.Status == "succeeded" {
			break loop
		}
	}

	d.SetId(projectid)
	log.Printf(projectid)
	if err != nil {
		return err
	}
	d.Partial(false)
	return nil
}

func resourceProjectDelete(d *schema.ResourceData, meta interface{}) error {

	c := meta.(*AzureDevOpsClient)
	client := meta.(*AzureDevOpsClient).coreClient
	clientoperation := meta.(*AzureDevOpsClient).operationClient
	ctx := meta.(*AzureDevOpsClient).StopContext

	var projectid = d.Id()
	idUUID, err := uuid.FromString(projectid)

	op, err := client.QueueDeleteProject(ctx, c.organization, idUUID)
	if err != nil {
		return fmt.Errorf("Error deleting project  %q: %+v", projectid, err)
	}

loop:
	for t := range time.NewTicker(10 * time.Second).C {
		utils.PeriodicFunc(t)
		op, err := clientoperation.GetOperation(ctx, *op.ID, c.organization, nil)
		if err != nil {
			return fmt.Errorf("Error getting project creating operation  %q: %+v", op.ID, err)
		}
		if op.Status == "succeeded" {
			break loop
		}
	}

	return nil
}
