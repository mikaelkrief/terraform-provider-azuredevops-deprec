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

  buildnumber_format = "1.0$(rev:.r)" #Optionnal

  repository {
    name   = "${azuredevops_project.project.name}"
    type   = "TfsGit"
    branch = "master"
  }

  variables {
    variable {
      name  = "test3"
      value = "ok"
    }

    variable {
      name      = "test4"
      value     = "ok2"
      is_secret = true
    }


  }

  designer_phase {
    name = "phase1"

    step {
      display_name = "teststep"
      task_id      = "d9bafed4-0b18-4f58-968d-86655b4d2ce9"
      task_version = "2.*"

      inputs = {
        failOnStderr     = "false"
        script           = "echo Write your commands here\necho Use the environment variables input below to pass secret variables to this script"
        workingDirectory = ""
      }
    }

    step {
      display_name       = "teststep2"
      task_id            = "d9bafed4-0b18-4f58-968d-86655b4d2ce9"
      task_version       = "2.*"
      enabled            = false                                  #Optionnal
      continue_on_error  = false                                  #Optionnal
      condition          = "always()"                             #Optionnal
      timeout_in_minutes = 50                                     #Optionnal

      inputs = {
        failOnStderr     = "false"
        script           = "echo Write your step 2"
        workingDirectory = "$$(Buid.SourcesDirectory)"
      }
    }
  }

  designer_phase {
    name = "phase2"

    step {
      display_name   = "teststep3"
      task_id        = "d9bafed4-0b18-4f58-968d-86655b4d2ce9"
      task_version   = "2.*"
      enabled        = true                                   #Optionnal
      always_run     = false                                  #Optionnal
      reference_name = "testouput"                            #Optionnal

      inputs = {
        failOnStderr     = "true"
        script           = "echo Write your commands here\necho Use the environment variables input below to pass secret variables to this script"
        workingDirectory = ""
      }

      #Optionnal
      environment_variables = {
        var1 = "key1"
        var2 = "key2"
      }
    }
  }

  queue {
    pool_name = "Hosted VS2017"
  }
}
