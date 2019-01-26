package azuredevops

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	azuredevopssdk "go-azuredevops-sdk"

	"github.com/hashicorp/terraform/helper/schema"
	//azuredevopssdk "github.com/mikaelkrief/go-azuredevops-sdk"
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
				Required: true,
				ForceNew: false,
			},
			"source_control_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"template_type_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"visibility": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "private",
			},
		},
	}
}

func periodicFunc(tick time.Time) {
	fmt.Println("Tick at: ", tick)
}

func resourceProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*azuredevopssdk.Client)

	//var versioncontrol = azuredevopssdk.Versioncontrol{}
	//versioncontrol.SourceControlType = d.Get("source_control_type").(string)

	// var processTemplate = azuredevopssdk.ProcessTemplate{}
	// processTemplate.TemplateTypeId = d.Get("template_type_id").(string)

	//var capabilities = &azuredevopssdk.Capabilities{}
	//capabilities.Versioncontrol = versioncontrol
	//capabilities.ProcessTemplate = processTemplate

	var project = azuredevopssdk.Project{}
	project.Name = d.Get("name").(string)
	project.Description = d.Get("description").(string)
	//project.Capabilities = capabilities
	project.Visibility = d.Get("visibility").(string)
	project.Capabilities.Versioncontrol.SourceControlType = d.Get("source_control_type").(string)
	project.Capabilities.ProcessTemplate.TemplateTypeId = d.Get("template_type_id").(string)

	log.Printf(project.Name)

	// client.ShowProject(project)
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
		return fmt.Errorf("Error geting project  %q: %+v", project.Name)
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

	if d.HasChange("template_type_id") {
		return fmt.Errorf("The template type Id can't be changed")
	}

	//var name = d.Get("name").(string)
	//var description = d.Get("description").(string)

	var project = azuredevopssdk.Project{}
	project.Name = d.Get("name").(string)
	project.Description = d.Get("description").(string)
	//project.Capabilities = ""
	//project.Capabilities.ProcessTemplate.TemplateTypeId = ""

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
	d.Set("template_type_id", project.Capabilities.ProcessTemplate.TemplateTypeId)

	return nil
}

func resourceProjectDelete(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		log.Println(string(b))
	}
	return
}
