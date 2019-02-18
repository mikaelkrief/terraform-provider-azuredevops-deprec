package azuredevopssdk

import "fmt"
import "log"
import "net/http"
import "encoding/json"

func (s *Client) GetOperation(id string) (string, error) {
	url := fmt.Sprintf(baseURL+"%s/_apis/operations/%s?api-version=4.1", s.organization, id)
	log.Printf(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	bytes, err := s.doRequest(req)
	if err != nil {
		return "", err
	}

	var resp ResponseOperation
	json.Unmarshal(bytes, &resp)
	return resp.Status, nil
}
