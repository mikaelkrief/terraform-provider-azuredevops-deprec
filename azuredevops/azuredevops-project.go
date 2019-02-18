package azuredevops

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"terraform-azuredevops/azuredevops/utils"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	azuredevopssdk "github.com/mikaelkrief/go-azuredevops-sdk"
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
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return true
				},
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

func periodicFunc(tick time.Time) {
	fmt.Println("Tick at: ", tick)
}

func resourceProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*azuredevopssdk.Client)

	var project = azuredevopssdk.Project{}
	project.Name = d.Get("name").(string)
	project.Description = d.Get("description").(string)

	var versioncontrol = azuredevopssdk.Versioncontrol{}
	versioncontrol.SourceControlType = d.Get("source_control_type").(string)

	var processTemplate = azuredevopssdk.ProcessTemplate{}
	log.Printf("TEST")
	log.Printf(d.Get("template_type_name").(string))
	if d.Get("template_type_name").(string) != "" {
		var processname = d.Get("template_type_name").(string)
		process, err := client.GetProcessId(processname)
		if err != nil {
			return fmt.Errorf("Error get process template id %+v", err)
		}
		log.Printf(process.Id)
		processTemplate.TemplateTypeId = process.Id
	} else {
		process, err := client.GetDefaultProcess()
		if err != nil {
			return fmt.Errorf("Error get default process template %+v", err)
		}
		processTemplate.TemplateTypeId = process.Id
		d.Set("template_type_name", process.Name)
	}

	var capabilities = &azuredevopssdk.Capabilities{}
	capabilities.Versioncontrol = versioncontrol
	capabilities.ProcessTemplate = processTemplate

	project.Capabilities = capabilities

	utils.PrettyPrint(project) //Log Request

	id, err := client.CreateProject(project)
	if err != nil {
		return fmt.Errorf("Error creating project  %q: %+v", project.Name, err)
	}

loop:
	for t := range time.NewTicker(10 * time.Second).C {
		periodicFunc(t)
		status, err := client.GetOperation(id)
		if err != nil {
			return fmt.Errorf("Error getting project creating operation  %q: %+v", id, err)
		}
		if status == "succeeded" {
			break loop
		}
	}

	projectCreated, err := client.GetProject(project.Name)
	if err != nil {
		return fmt.Errorf("Error getting project  %q:", project.Name)
	}
	var idproject = projectCreated.Id

	d.SetId(idproject)
	log.Printf(idproject)
	if err != nil {
		return err
	}
	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*azuredevopssdk.Client)

	d.Partial(true)

	if d.HasChange("source_control_type") {
		return fmt.Errorf("The source controle type can't be changed")
	}

	if d.HasChange("template_type_name") {
		return fmt.Errorf("The template type Id can't be changed")
	}

	var project = azuredevopssdk.Project{}
	if d.HasChange("name") {
		project.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		project.Description = d.Get("description").(string)
	}

	b, err := json.Marshal(project)
	log.Printf("[INFO] project of update. %s", string(b))

	var projectid = d.Id()

	id, err := client.UpdateProject(projectid, project)
	if err != nil {
		return fmt.Errorf("Error updating project  %q: %+v", project.Name, err)
	}

loop:
	for t := range time.NewTicker(10 * time.Second).C {
		periodicFunc(t)
		status, err := client.GetOperation(id)
		if err != nil {
			return fmt.Errorf("Error getting project updating operation  %q: %+v", id, err)
		}
		if status == "succeeded" {
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

func resourceProjectRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*azuredevopssdk.Client)
	var projectname = d.Get("name").(string)
	project, err := client.GetProject(projectname)
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error getting project  %q: %+v", projectname, err)
	}

	d.Set("name", project.Name)
	d.Set("description", project.Description)
	d.Set("visibility", project.Visibility)
	d.Set("source_control_type", project.Capabilities.Versioncontrol.SourceControlType)

	return nil
}

func resourceProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*azuredevopssdk.Client)
	var projectid = d.Id()

	id, err := client.DeleteProject(projectid)
	if err != nil {
		return fmt.Errorf("Error deleting project  %q: %+v", projectid, err)
	}

loop:
	for t := range time.NewTicker(10 * time.Second).C {
		periodicFunc(t)
		status, err := client.GetOperation(id)
		if err != nil {
			return fmt.Errorf("Error getting project deleting operation  %q: %+v", id, err)
		}
		if status == "succeeded" {
			break loop
		}
	}
	return nil
}
