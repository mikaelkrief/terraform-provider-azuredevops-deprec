provider "azuredevops" {}

resource "random_id" "project_name" {
  byte_length = 2
}

resource "azuredevops_project" "project" {
  name               = "Terra-${random_id.project_name.hex}-3"
  template_type_name = "agile"
  description        = "my project terraform ${random_id.project_name.hex}"
  template_type_name = "Scrum"
}

resource "azuredevops_build_definition" "test1" {
  name       = "build-def-${random_id.project_name.hex}"
  project_id = "${azuredevops_project.project.name}"

  repository {
    name   = "${azuredevops_project.project.name}"
    type   = "TfsGit"
    branch = "master"
  }

  designer_phase {
    name = "phase1"

    steps {
      display_name = "teststep"
      task_id      = "d9bafed4-0b18-4f58-968d-86655b4d2ce9"

      inputs = {
        failOnStderr     = "false"
        script           = "echo Write your commands here\necho Use the environment variables input below to pass secret variables to this script"
        workingDirectory = ""
      }
    }

    steps {
      display_name = "teststep2"
      task_id      = "d9bafed4-0b18-4f58-968d-86655b4d2ce9"

      inputs = {
        failOnStderr     = "false"
        script           = "echo Write your commands here\necho Use the environment variables input below to pass secret variables to this script"
        workingDirectory = ""
      }
    }
  }

  designer_phase {
    name = "phase2"

    steps {
      display_name = "teststep3"
      task_id      = "d9bafed4-0b18-4f58-968d-86655b4d2ce9"

      inputs = {
        failOnStderr     = "false"
        script           = "echo Write your commands here\necho Use the environment variables input below to pass secret variables to this script"
        workingDirectory = ""
      }
    }
  }

  queue {
    pool_name = "Hosted VS2017"
  }
}
