package main

import (
	"context"
	"fmt"
	"image/png"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"bitbucket.org/latonaio/aion-core/pkg/go-client/msclient"
	"bitbucket.org/latonaio/aion-core/pkg/log"
	"bitbucket.org/latonaio/aion-core/proto/kanbanpb"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

func main() {
	// Create Kanban client
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	kanbanClient, err := msclient.NewKanbanClient(ctx, msName, kanbanpb.InitializeType_START_SERVICE)
	if err != nil {
		log.Fatalf("failed to get kanban client: %v", err)
	}
	log.Printf("successful get kanban client")
	defer kanbanClient.Close()

	kanbanCh := kanbanClient.GetKanbanCh()
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
			fromMetadata, err := msclient.GetMetadataByMap(k)
			if err != nil {
				errCh <- fmt.Errorf("failed to get metadata: %v", err)
				continue
			}
			log.Printf("got metadata from kanban")
			log.Debugf("metadata: %v", fromMetadata)

			sizeStr, ok := fromMetadata["size"].(string)
			if !ok {
				errCh <- fmt.Errorf("failed to convert size to string")
				continue
			}
			size, err := strconv.Atoi(sizeStr)
			if err != nil {
				errCh <- fmt.Errorf("failed to convert size to int")
				continue
			}
			jsonStr, ok := fromMetadata["json_str"].(string)
			if !ok {
				errCh <- fmt.Errorf("failed to convert json_str to string")
				continue
			}
			outputPath, ok := fromMetadata["output_path"].(string)
			if !ok {
				errCh <- fmt.Errorf("failed to convert output_path to string")
				continue
			}
			option, ok := fromMetadata["option"].(string)
			if !ok {
				errCh <- fmt.Errorf("failed to convert option to string")
				continue
			}

			qrCode, err := qr.Encode(jsonStr, qr.M, qr.Auto)
			if err != nil {
				errCh <- fmt.Errorf("failed to encode json to QR code: %v", err)
				continue
			}
			qrCode, err = barcode.Scale(qrCode, size, size)
			if err != nil {
				errCh <- fmt.Errorf("failed to scale QR code: %v", err)
				continue
			}
			file, err := os.Create(outputPath)
			if err != nil {
				errCh <- fmt.Errorf("failed to create file: %v", err)
				continue
			}
			if err = png.Encode(file, qrCode); err != nil {
				errCh <- fmt.Errorf("failed to encode QR code to png: %v", err)
				continue
			}
			file.Close()

			toMetadata := map[string]interface{}{"file_path": outputPath, "option": option}

			// Write metadata to Kanban
			if err := writeKanban(kanbanClient, toMetadata); err != nil {
				errCh <- fmt.Errorf("failed to write kanban: %v", err)
				continue
			}
			log.Printf("write metadata to kanban")
		}
	}
END:
}
