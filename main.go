package main

import (
	"encoding/json"
	"log"
	"os"
	"playbooks"
)

func main() {
	// collectors.GetProjects()
	// assets := collectors.GetAssetsByProjectId(nil, "cloudlabs-371516")
	// log.Println(assets)
	// for k, v := range assets.(map[string]collectors.Project) {
	// 	log.Printf("KEY: %v\n%v", k, v.Parents)
	// }
	// mappings.CallProvider()
	// mappings.ReadFileToReverse("mappings/data_mapping.json", "mappings/data_reverse.json")
	// mappings.DoReverse("mappings/resource_mapping.json")
	// mappings.DoReverse()
	// tfResource := playbooks.TFResourceType{}
	// tfResource.CheckDirectoryExistsOrCreate("./Factory", "DigitalFactory")

	playbooks.ExecutePlayBook("./playbooks","assets.json")
	// defer log.Println(playbooks.ListOfAllTFResourcesPerPath)
	defer func(){
		if data, err := json.Marshal(playbooks.ListOfAllTFResourcesPerPath); err == nil {
			if os.WriteFile("./resources.json", data, 0777) != nil {
				log.Printf("FAILED TO CREATE THE OUTPUT FILE")
			}
		}
	}()
}