package azuredevopssdk

type Project struct {
	Id           string        `json:"id,omitempty"`
	Name         string        `json:"name,omitempty"`
	Description  string        `json:"description"`
	Visibility   string        `json:"visibility,omitempty"`
	Capabilities *Capabilities `json:"capabilities,omitempty"`
}

type Capabilities struct {
	Versioncontrol  Versioncontrol  `json:"versioncontrol,omitempty"`
	ProcessTemplate ProcessTemplate `json:"processTemplate,omitempty"`
}

type Versioncontrol struct {
	SourceControlType string `json:"sourceControlType,omitempty"`
}

type ProcessTemplate struct {
	TemplateTypeId string `json:"templateTypeId,omitempty"`
}

//ResponseProject
type ResponseProject struct {
	Id     string `json:"id"`
	Status string `json:"status"`
	Url    string `json:"url"`
}
