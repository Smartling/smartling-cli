/*
	Smartling SDK v2 files api sample

	Sample shows usage of smartling files api
	http://docs.smartling.com/pages/API/v2/FileAPI/
	Replace userIdentifier and tokenSecret with your credentials
	
*/

package main

import (
	"time"
	"log"

	// should be replaced with "github.com/Smartling/smartling-cli" 
	// relative path is here only for making it easy to run out of the box
	"./../../smartling"
)

const (
	userIdentifier = "" // put your user identifier here
	tokenSecret    = "" // put your token secret here
	projectId      = "" // put your project id here
)

func main() {

	log.Printf("Initializing smartling client and performing autorization")
	client := smartling.NewClient(userIdentifier, tokenSecret)

	log.Printf("Listing project (%v) files", projectId)

	log.Printf("Listing constraints:")
	log.Printf("UriMask : 'master' (file URI must contain 'master' substring)")
	log.Printf("LastUploadedBefore : one month back")
	log.Printf("FileTypes : json and Java properties files")
	listRequest := smartling.FileListRequest {
		UriMask: "master",
		LastUploadedBefore: smartling.Iso8601Time(time.Now().Add(time.Hour * 31 * 24 * -1)),
		FileTypes: []smartling.FileType {smartling.Json, smartling.JavaProperties},
	}
	log.Println()

	files, err := client.ListFiles(projectId, listRequest)
	if err != nil {
		log.Printf("%v", err.Error())
		return
	}
	log.Printf("Success!")

	log.Printf("Found %v files", files.TotalCount)

	for _, f := range files.Items {
		log.Printf("%v", f.FileUri)
	}
}
