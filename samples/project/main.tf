provider "azuredevops" {

}

resource "azuredevops_project" "project" {
  name               = "Terra-${random_id.project_name.hex}-3"
  template_type_name = "agile"

  description = "my project terraform ${random_id.project_name.hex}"

  template_type_name = "Scrum"

  /*source_control_type = "Git"
  template_type_name  = "Scrum1"
  visibility = "private"*/
}

resource "random_id" "project_name" {
  byte_length = 2
}