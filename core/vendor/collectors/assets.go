package collectors

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	rsm "cloud.google.com/go/resourcemanager/apiv3"
	rpb "cloud.google.com/go/resourcemanager/apiv3/resourcemanagerpb"
	cam "google.golang.org/api/cloudasset/v1"
)

var (
	ActiveProjects []map[string]string

	AllProjectsAssets []map[string]Project
)

type AncestorsTree struct {
	Name      string
	Id        string
	Parent    string
	Children  []AncestorsTree
	Ancestors []string
}



func GetContext() context.Context {
	ctx := context.Background()
	return ctx
}

// get all projects assets
func GetAssetsByProjectId(c chan []map[string]Project) {

	camService, _ := cam.NewService(GetContext())

	// create parent folders for the resources in the projects listed previously
	for _, project := range ActiveProjects {
		// get the project ID
		pid := project["projectId"]

		// if !strings.Contains(pid, "cloudlabs") {
		// 	continue
		// }

		// list all assets in the project
		call := camService.Assets.List(fmt.Sprintf("projects/%s", pid))
		call.ContentType("RESOURCE")
		var holder map[string]interface{}
		var project = NewProjectInstance()
		
		// start the call
		if response, err := call.Do(); err == nil {

			for index, assets := range response.Assets {

				resolve, _ := assets.MarshalJSON()
				// log.Println(resolve)

				if err := json .Unmarshal(resolve, &holder); err != nil {
					log.Printf("[WARNING] %s: %v", err, assets)
				}

				// log.Println(holder)
				parent := holder["ancestors"].([]interface{})[0]

				key := parent.(string)
				project.SetName(key)

				if index == 0 {
					// Ancestors[key] = append(Ancestors[key], holder["ancestors"].([]interface{})[1:])
					project.SetParents(holder["ancestors"].([]interface{})[1:])
				}
				// Add assets
				project.AddAssets(holder["assetType"].(string), holder["name"].(string))
				// Ancestors[key] = append(Ancestors[key], holder["assetType"])
			}
		} else {
			log.Fatal("Failed to retrieve assets in the project: ", pid, "\n", err)
		}
		AllProjectsAssets = append(AllProjectsAssets, map[string]Project{
			pid : project,
		})
	}
	c <- AllProjectsAssets
}

func FetchResourceNameById(s string) any {
	var holder any
	if strings.Contains(s, "projects") {
		// Retrieve the project display name
		// call := CrmService.Projects .Get(s)
		req := &rpb.GetProjectRequest{Name: s}
		cli, err := rsm.NewProjectsClient(GetContext())
		if err != nil {
			log.Fatal(err)
		}
		if project, err := cli.GetProject(GetContext(), req); err == nil {
			holder = []string{project.DisplayName, project.ProjectId,}
		}
		log.Default()
		defer cli.Close()
	}

	if strings.Contains(s, "folders"){
		req := &rpb.GetFolderRequest{Name: s}
		cli, err := rsm.NewFoldersClient(GetContext())
		if err != nil {
			log.Fatal("[ERROR] Folders client failure...")
		}
		defer cli.Close()
		if folder, err := cli.GetFolder(GetContext(), req); err == nil {
			holder = folder.DisplayName
			log.Default()
		}
	}

	if strings.Contains(s, "organizations") {
		req := &rpb.GetOrganizationRequest{Name: s}
		cli, err := rsm.NewOrganizationsClient(GetContext())
		if err != nil {
			log.Fatal("[ERROR] Folders client failure...")
		}
		defer cli.Close()
		if org, err := cli.GetOrganization(GetContext(), req); err == nil {
			holder = org.DisplayName
			log.Default()
		}
	}
	return holder
}