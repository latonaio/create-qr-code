package main

import (
	"context"
	"fmt"
	"image/png"
	"os"
	"os/signal"
	"syscall"

	"bitbucket.org/latonaio/aion-core/pkg/go-client/msclient"
	"bitbucket.org/latonaio/aion-core/pkg/log"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

func main() {
	// Create Kanban client
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	kanbanClient, err := msclient.NewKanbanClient(ctx, msName)
	if err != nil {
		log.Fatalf("failed to get kanban client: %v", err)
	}
	log.Printf("successful get kanban client")
	defer kanbanClient.Close()

	kanbanCh, err := kanbanClient.GetKanbanCh()
	if err != nil {
		log.Fatalf("failed to get kanban channel: %v", err)
	}
	log.Printf("successfull get kanban channel")

	errCh := make(chan error, 1)
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM)

	for {
		select {
		case s := <-signalCh:
			fmt.Printf("recieved signal: %s", s.String())
			goto END
		case err := <-errCh:
			log.Errorf("error: %v", err)
		case k := <-kanbanCh:
			if k == nil {
				continue
			}
			// Get metadata from Kanban
			fromMetadata, err := k.GetMetadataByMap()
			if err != nil {
				errCh <- fmt.Errorf("failed to get metadata: %v", err)
			}
			log.Printf("got metadata from kanban")
			log.Debugf("metadata: %v", fromMetadata)

			size, ok := fromMetadata["size"].(int)
			if !ok {
				errCh <- fmt.Errorf("failed to convert interface{} to string")
			}
			json_str, ok := fromMetadata["json_str"].(string)
			if !ok {
				errCh <- fmt.Errorf("failed to convert interface{} to string")
			}
			output_path, ok := fromMetadata["output_path"].(string)
			if !ok {
				errCh <- fmt.Errorf("failed to convert interface{} to string")
			}

			qrCode, _ := qr.Encode(json_str, qr.M, qr.Auto)
			qrCode, _ = barcode.Scale(qrCode, size, size)

			file, _ := os.Create(output_path)
			png.Encode(file, qrCode)
			file.Close()

			toMetadata := map[string]interface{}{"file_path": output_path}

			// Write metadata to Kanban
			if err := writeKanban(kanbanClient, toMetadata); err != nil {
				errCh <- fmt.Errorf("failed to write kanban: %v", err)
			}
			log.Printf("write metadata to kanban")
		}
	}
END:
}
