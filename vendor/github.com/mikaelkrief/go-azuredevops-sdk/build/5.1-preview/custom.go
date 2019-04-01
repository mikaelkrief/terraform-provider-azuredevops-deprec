package build

import "encoding/json"

func (d *DesignerProcess) UnmarshalJSON(data []byte) error {

	type DesignerProcess2 DesignerProcess
	var dp2 DesignerProcess2
	if err := json.Unmarshal(data, &dp2); err != nil {
		return err
	}


	*d = DesignerProcess(dp2)
	return nil
}
