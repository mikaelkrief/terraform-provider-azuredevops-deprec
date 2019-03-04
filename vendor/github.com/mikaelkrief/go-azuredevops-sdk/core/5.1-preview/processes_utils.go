package core

import (
	"context"
	"fmt"
	"strings"
)

//GetDefaultProcess : get the default process
func (client BaseClient) GetDefaultProcess(ctx context.Context, organization string) (result *Process, err error) {

	processes, err := client.GetProcesses(ctx, organization)
	if err != nil {
		return nil, err
	}

	var listprocess = processes.Value
	for index := 0; index < len(*listprocess); index++ {
		var proc = (*listprocess)[index]
		if *proc.IsDefault {
			process := proc
			return &process, nil
		}
	}
	return nil, nil

}

//GetProcessId : get thr process if by name
func (client BaseClient) GetProcessIdbyName(ctx context.Context, organization string, processname string) (result *Process, err error) {

	processes, err := client.GetProcesses(ctx, organization)
	if err != nil {
		return nil, err
	}

	var listprocess = processes.Value
	for index := 0; index < len(*listprocess); index++ {
		var proc = (*listprocess)[index]
		if strings.ToLower(processname) == strings.ToLower(*proc.Name) {
			process := proc
			return &process, nil
		}
	}

	return nil, fmt.Errorf("error the template process %+v doesn't exist", processname)

}
