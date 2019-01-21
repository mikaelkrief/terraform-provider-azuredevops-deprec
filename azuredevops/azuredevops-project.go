package azuredevops

import (
	azuredevopssdk "github.com/mikaelkrief/go-azuredevops-sdk"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"time"
	"fmt"
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
			"abbreviation": {
				Type:     schema.TypeString,
				Required: false,
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
				Required: true,
				ForceNew: false,
				Default: "Prive",
			},
		},
	}
}

func periodicFunc(tick time.Time){
    fmt.Println("Tick at: ", tick)
}

func resourceProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*azuredevopssdk.Client)

	var versioncontrol = azuredevopssdk.Versioncontrol{}
	versioncontrol.SourceControlType = d.Get("source_control_type").(string)

	var processTemplate = azuredevopssdk.ProcessTemplate{}
	processTemplate.TemplateTypeId = d.Get("template_type_id").(string)

	var capabilities = azuredevopssdk.Capabilities{}
	capabilities.Versioncontrol = versioncontrol
	capabilities.ProcessTemplate = processTemplate

	var project = azuredevopssdk.Project{}
	project.Name = d.Get("name").(string)
	project.Description = d.Get("description").(string)
	project.Abbreviation = d.Get("abbreviation").(string)
	project.Capabilities = capabilities
	project.Visibility = d.Get("visibility").(string)

	log.Printf(project.Name)

	// client.ShowProject(project)
	id, err := client.CreateProject(project)
	if  err != nil {
		return fmt.Errorf("Error creating project  %q: %+v", project.Name, err)
	}

loop:
	for t := range time.NewTicker(10 * time.Second).C {
		periodicFunc(t)
		status, err := client.GetOperation(id)
	   if err !=nil {
		return fmt.Errorf("Error geting project creating operation  %q: %+v", id, err)
	   }
	   if(status == "succeeded"){
		break loop
	   }
	}
	
	d.SetId(id)
	log.Printf(id)
	if err != nil {
		return err
	}
	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*azuredevopssdk.Client)

	
	if d.HasChange("source_control_type") {
		return fmt.Errorf("The source controle type can't be changed")
	}

	if d.HasChange("template_type_id") {
		return fmt.Errorf("The template type Id can't be changed")
	}

	var name = d.Get("name").(string)
	var description = d.Get("description").(string)

	var project = azuredevopssdk.Project{
		Name: name,
		Description: description,
	}
	log.Printf("[INFO] end of update.%s",project.Name)
	return nil

}

func resourceProjectRead(d *schema.ResourceData, meta interface{}) error {
    //TODO: This really should hit a get API to see what the current state is.  Something like the following sudo code
    // client := meta.(*classis.Client)
    // spotGroup, err := client.ReadSpotGroup(d.Id())
    // if err then d.SetId("")
    // else set and changed values
	return nil
}

func resourceProjectDelete(d *schema.ResourceData, meta interface{}) error {

	return nil
}