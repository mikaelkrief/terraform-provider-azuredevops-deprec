package azuredevopssdk

import "fmt"
import "log"
import "net/http"
import "io/ioutil"
import "encoding/base64"
import "time"

const baseURL string = "https://dev.azure.com/"

//Client client struct
type Client struct {
	client       *http.Client
	organization string
	encToken     string
}

//NewClientWith initialize a new client
func NewClientWith(organization string, token string) (*Client, error) {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	return &Client{netClient, organization, basicAuth(":" + token)}, nil
}

func basicAuth(token string) string {
	auth := token
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func (s *Client) doRequest(req *http.Request) ([]byte, error) {
	// dispay log the truncated secret
	runeSecret := []rune(s.encToken)
	myShortSecret := string(runeSecret[0:6])
	log.Printf("[SECRET] --> " + myShortSecret)
	//----------------------------------------

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+s.encToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if 202 == resp.StatusCode {
		return body, nil
	}

	if 200 != resp.StatusCode {
		if resp.StatusCode == 203 {
			return nil, fmt.Errorf("%s", "BAD TOKEN")
		} else {
			return nil, fmt.Errorf("%s", body)
		}
	}
	return body, nil
}
