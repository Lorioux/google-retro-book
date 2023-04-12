package main

import (
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
	// mappings.ReadFileDoReverse("mappings/resource_mapping.json", "mappings/resource_reverse.json")
	// tfResource := playbooks.TFResourceType{}
	// tfResource.CheckDirectoryExistsOrCreate("./Factory", "DigitalFactory")

	playbooks.ExecutePlayBook("./playbooks","assets.json")
}

