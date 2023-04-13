package playbooks

import (
	"collectors"
	"encoding/json"
	"log"
	"mappings"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

type Ancestors []string

type Asset struct {
	Ancestors Ancestors `json:"ancestors"`
	AssetType string `json:"assetType"`
	Name string `json:"name"`
	// updateTime string `json:update_time`
}

func (a *Asset) PickResourceName () string {
	name := strings.Split(a.Name, "/")[3:]
	return strings.Join(name, "/")
}

func (a *Asset) PickResourceType () string {
	aType := strings.Split(a.AssetType, "/")[1]
	return aType //strings.ToTitle(aType)
}

func (a *Asset) PickServiceType () string {
	sType := strings.Split(a.AssetType, ".")[0]
	
	return strings.ToTitle(sType)  
}

var content []Asset

func ExecutePlayBook(root string, assetsFile string){
	orgRootPath = root
	file, err  := os.ReadFile(path.Join( root, assetsFile))
	// defer file.Close()
	// var broker map[string]interface{}
	
	if err == nil {
		if err := json.Unmarshal(file, &content); err != nil {
			log.Println(err)
		} else {
			// for k, v := range
			// log.Printf("My name: %v and Type is: %v and Service is: %v", 
			// 	content[0].PickResourceName(), 
			// 	content[0].PickResourceType(),
			// 	content[0].PickServiceType(),
			// )
			// for c := range content {
			// 	if strings.Contains(content[c].PickServiceType(), "MANAGER"){
			// 		log.Printf("My name: %v and Type is: %v and Service is: %v", 
			// 			content[c].PickResourceName(), 
			// 			content[c].PickResourceType(),
			// 			content[c].PickServiceType(),
			// 		)
			// 	}
			// 	// log.Println(content[c].PickServiceType())
			// }
			SortOrgNodeAllowedResource()
		}
		
	} else {
		log.Fatalln("Error: ", err)
	}
}

var OrgAllowedResources map[string]Asset
var OrgPrimaryResourceTypes = []string{"Project", "Folder", "TagKey", "TagBinding", "TagKey", "TagValue"}

type OrgNodeTree struct {
	Name string
	ServiceTypes []string 
	MetaMapping map[string]string
	// MetaModules []MetaModule
} 

type ServiceType struct{}

// func (o *OrgNodeTree) AddServiceType(p string) {
// 	if entries, err := os.ReadDir(p); err != nil {
// 		if err := os.MkdirAll(p, 0777); err != nil {
// 			log.Fatalf("Error: %v", err)
// 		} else {
// 			log.Print(entries)
// 		}
// 	}
// }

func SortOrgNodeAllowedResource(){
	cn := make(chan int)
	for index,con := range content {
		resourceType := &TFResourceType{}

		if !strings.EqualFold(con.AssetType, "cloudresourcemanager.googleapis.com/Folder"){
			continue
		}
		mux := &sync.RWMutex{}
		// Check if resource type is in the allowed list.
		time.Sleep(3 * time.Second) 
		go doWorkOnCoRoutine(OrgPrimaryResourceTypes, con, cn, index, resourceType, mux)
	}
	ConstructResourceHierarchy()
	defer mappings.UpdateUnsupportedResourceMap()
}


func doWorkOnCoRoutine(rs []string, a Asset, cx chan int, index int, resourceType *TFResourceType, mux *sync.RWMutex) {
	
	mux.Lock()	
	if name, err := url.QueryUnescape(a.PickResourceName()); err == nil {
		resourceType.Name = name
	}
	if parent, err := url.QueryUnescape(a.Ancestors[0]); err == nil {
		resourceType.Parent = parent
	}
	
	resourceType.SetRequiredFields(a.AssetType)
	
	mux.Unlock()
	if index == len(content) - 1 {
		cx <- -1
		return
	} else if strings.Contains(a.Ancestors[0], "organizations") && len(a.Ancestors) == 1 {
		cx <- index
		// log.Printf("Here1: %v", index)
		return
	} else {
		// log.Printf("Here2: %v", index)
		for _, r := range rs {
			if a.PickResourceType() == r {
				cx <- index
				break
			}
		}
		return
	}
}


var parentPathMap = map[string][2]any{}
// var orgNodeTree = OrgNodeTree{}
var orgRootPath string 

func ConstructResourceHierarchy () {
	resourceType := TFResourceType{}
	cx := make(chan []string)
	go func(con []Asset, mapPaths map[string][2]any, cn chan []string){
		for index, c := range con {

			var pathTemp string
			pathKey := c.Ancestors[0]
			if mapPaths[pathKey][0] == nil {
				// log.Printf("1:::: %v --- %v", mapPaths[pathKey][0], mapPaths[pathKey][1])
				for _, a := range c.Ancestors {
					nameT := collectors.FetchResourceNameById(a)
					pathTemp = path.Join(nameT.(string), pathTemp)	
				}
				// && mapPaths[pathKey][1] != nil
			} else {
				// log.Printf("2:::: %v --- %v", mapPaths[pathKey][0], mapPaths[pathKey][1])
				if index == len(content) - 1 {
					cn <- nil
				}
				continue
			}
			
			// log.Printf("Path: %v", pathTemp)
			if pathTemp != "" {
				cn <- []string{pathKey, pathTemp}
				mapPaths[pathKey] = [2]any{true, pathTemp}
				// log.Printf("3:::: %v --- %v", mapPaths[pathKey][0], mapPaths[pathKey][1])
			}
			
			if index == len(content) - 1 {
				cn <- nil
			}
		}
	}(content, parentPathMap, cx )

	cm := <- cx 
	for cm != nil {
		
		// log.Println(cm)
		resourceType.CheckDirectoryExistsOrCreate(orgRootPath, cm[1])
		cm = <- cx 
	}
}