package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
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

	// Time 2 weeks ago
	endTime := startTime.AddDate(0, 0, 14)

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

	// Bail out early if no files are found
	if len(files.Files) == 0 {
		fmt.Println("No matching files found, exiting")
		return
	}

	// Initialize a waitgroup
	var wg sync.WaitGroup

	// Iterate over and delete the files
	for _, file := range files.Files {
		fmt.Printf("Deleting %s\n", file.Name)
		wg.Add(1)
		go deleteFile(client, file.Id, &wg)
	}

	fmt.Println("Waiting for all goroutines to finish")
	wg.Wait()
	fmt.Println("All files deleted successfully")
}

// deleteFile deletes the file with the given ID from the connected drive
func deleteFile(client *drive.Service, fileID string, wg *sync.WaitGroup) {
	defer wg.Done()
	err := client.Files.Delete(fileID).SupportsAllDrives(true).Do()
	if err != nil {
		log.Fatalf("Failed to delete file: %v", err)
	}
}
