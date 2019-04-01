output "project_test_id" {
  description = "id of the project."
  value       = "${azuredevops_project.project.*.id}"
}

output "project_test_name" {
  description = "name of the project."
  value       = "${azuredevops_project.project.*.name}"
}