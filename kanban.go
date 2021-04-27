package main

import (
	"fmt"

	"bitbucket.org/latonaio/aion-core/pkg/go-client/msclient"
)

const msName = "create-qr-code"

func writeKanban(kanbanClient msclient.MicroserviceClient, data map[string]interface{}) error {
	var options []msclient.Option
	options = append(options, msclient.SetMetadata(data))
	options = append(options, msclient.SetProcessNumber(kanbanClient.GetProcessNumber()))
	req, err := msclient.NewOutputData(options...)
	if err != nil {
		return fmt.Errorf("failed to construct output request: %v", err)
	}
	if err := kanbanClient.OutputKanban(req); err != nil {
		return fmt.Errorf("failed to output to kanban: %v", err)
	}
	return nil
}
