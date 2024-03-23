package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()

	// Initialize drive service
	client, err := drive.NewService(ctx, option.WithCredentialsFile("service_account.json"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// ID of the folder to be cleaned
	var folderID string
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <folder_id>\n", os.Args[0])
		os.Exit(1)
	} else {
		folderID = os.Args[1]
	}

	// Time 4 weeks ago
	startTime := time.Now().AddDate(0, 0, -28)

	// Time 3 weeks ago
	endTime := startTime.AddDate(0, 0, 7)

	// List out the files matching our criteria
	files, err := client.
		Files.
		List().
		Q(
			fmt.Sprintf(
				"'%s' in parents AND modifiedTime > '%s' AND modifiedTime < '%s'",
				folderID,
				startTime.Format(time.RFC3339),
				endTime.Format(time.RFC3339),
			),
		).
		SupportsAllDrives(true).
		IncludeTeamDriveItems(true).
		Do()
	if err != nil {
		log.Fatalf("Failed to list files: %v", err)
	}

	// Iterate over and delete the files
	for _, file := range files.Files {
		fmt.Printf("Deleting %s\n", file.Name)
		err = client.Files.Delete(file.Id).SupportsAllDrives(true).Do()
		if err != nil {
			log.Fatalf("Failed to delete file: %v", err)
		}
	}
}
