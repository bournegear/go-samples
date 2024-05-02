package main

//The goal of this script is to query Jira API and returns a summary of the notification scheme of a Project
//The script returns the Notification Name, Notification Description, the Roles it notifies.
//Other than the env file when ran it prompts for the project key as an input to retrieve the summary

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func goDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load("creds.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func main() {
	// This returns the API key as a String
	var project_key string
	api_key := goDotEnvVariable("ATLASSIAN_API_KEY")
	username := goDotEnvVariable("USERNAME")
	base_url := goDotEnvVariable("BASE_URL")

	//prompt user for Project Key
	fmt.Print("Enter Project Key: ")
	fmt.Scanln(&project_key)

	api_route := fmt.Sprintf("/project/%s/notificationscheme?expand=notificationSchemeEvents", project_key)
	// Build the full URL to GET
	url := base_url + api_route

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	//Set Authentication and headers
	request.Header.Set("Accept", "application/json")
	request.SetBasicAuth(username, api_key)
	// Make the request
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("Error Sending request:", err)
		return
	}
	defer resp.Body.Close()
	//Process the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}
	// Process the JSON
	var hierarchyData interface{}
	err = json.Unmarshal(body, &hierarchyData)
	if err != nil {
		fmt.Println("error unmarshalling JSON", err)
		return
	}

	notificationSchemeMap := hierarchyData.(map[string]interface{})

	//Print the Notification Scheme Name
	fmt.Println(notificationSchemeMap["name"])
	fmt.Println("===============================")

	eventList, ok := notificationSchemeMap["notificationSchemeEvents"].([]interface{})
	if !ok {
		fmt.Println("Error: notificationSchemeEvents is not a list")
		return
	}

	for _, event := range eventList {
		//Empty array to print all the roles that  receive a notification
		listOfRole := []string{}
		eventMap, ok := event.(map[string]interface{})
		if !ok {
			fmt.Println("Error: Event item is not a map")
			continue
		}
		eventName := eventMap["event"].(map[string]interface{})["name"].(string)
		fmt.Println("Notification Name: ", eventName)
		eventDescription := eventMap["event"].(map[string]interface{})["description"].(string)
		fmt.Println("Notification Description:", eventDescription)

		notificationTypes, ok := eventMap["notifications"].([]interface{})
		if !ok {
			fmt.Println("Error: notifications is not a map")
		}
		for _, notification := range notificationTypes {
			notificationMap, ok := notification.(map[string]interface{})
			if !ok {
				fmt.Println("Error: Notifications are not a map")
				continue
			}
			notifies := notificationMap["notificationType"].(string)
			listOfRole = append(listOfRole, notifies)
			// fmt.Println("Notifies: ", notifies)
		}
		fmt.Println("Roles: ", listOfRole)
		fmt.Println("")
	}

}
