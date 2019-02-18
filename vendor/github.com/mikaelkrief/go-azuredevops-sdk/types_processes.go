package azuredevopssdk

type Processe struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsDefault   bool   `json:"isDefault,omitempty"`
	Type        string `json:"type,omitempty"`
}

type Processes struct {
	Count        int        `json:"count"`
	ProcesseList []Processe `json:"value,omitempty"`
}
