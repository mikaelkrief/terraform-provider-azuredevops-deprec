package azuredevopssdk

import "fmt"
import "log"
import "net/http"
import "encoding/json"
import "bytes"

func (s *Client) CreateProject(project Project) (string, error) {
	var a [1]interface{}
	a[0] = project

	projectBytes, _ := json.Marshal(a[0])
	projectReader := bytes.NewReader(projectBytes)
	url := fmt.Sprintf(baseURL+"%s/_apis/projects?api-version=5.0-preview.3", s.organization)
	log.Printf(url)
	req, err := http.NewRequest("POST", url, projectReader)
	if err != nil {
		return "", err
	}
	bytes, err := s.doRequest(req)
	if err != nil {
		return "", err
	}

	var resp ResponseProject
	json.Unmarshal(bytes, &resp)
	return resp.Id, nil
}

func (s *Client) UpdateProject(projectid string, projectUpdated Project) (string, error) {
	var a [1]interface{}
	a[0] = projectUpdated

	projectBytes, _ := json.Marshal(a[0])
	projectReader := bytes.NewReader(projectBytes)
	url := fmt.Sprintf(baseURL+"%s/_apis/projects/%s?api-version=5.0-preview.3", s.organization, projectid)
	log.Printf(url)
	req, err := http.NewRequest("PATCH", url, projectReader)
	if err != nil {
		return "", err
	}
	bytes, err := s.doRequest(req)
	if err != nil {
		return "", err
	}

	var resp ResponseProject
	json.Unmarshal(bytes, &resp)
	return resp.Id, nil
}

func (s *Client) GetProject(name string) (Project, error) {
	url := fmt.Sprintf(baseURL+"%s/_apis/projects/%s?includeCapabilities=true&api-version=5.0-preview.3", s.organization, name)
	log.Printf(url)
	var project Project
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return project, err
	}
	bytes, err := s.doRequest(req)
	if err != nil {
		return project, err
	}

	json.Unmarshal(bytes, &project)
	// log.Printf("project: %+v\n", project)
	return project, nil
}


func (s *Client) DeleteProject(projectid string) (string, error) {

	url := fmt.Sprintf(baseURL+"%s/_apis/projects/%s?api-version=5.0-preview.3", s.organization, projectid)
	log.Printf(url)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return "", err
	}
	bytes, err := s.doRequest(req)
	if err != nil {
		return "", err
	}

	var resp ResponseProject
	json.Unmarshal(bytes, &resp)
	return resp.Id, nil
}

func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		log.Println(string(b))
	}
	return
}
