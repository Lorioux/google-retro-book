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
)


var OrgAllowedResources map[string]TFResourceType
var OrgPrimaryResourceTypes = []string{"Project", "Folder", "TagKey", "TagBinding", "TagKey", "TagValue"}
var parentPathMap = map[string][2]any{}
// var orgNodeTree = OrgNodeTree{}
var orgRootPath string 
var content []Asset
var control *sync.Map


// type OrgNodeTree struct {
// 	Name string
// 	ServiceTypes []string 
// 	MetaMapping map[string]string
// 	// MetaModules []MetaModule
// } 

func init(){
	control = &sync.Map{}
}

func ExecutePlayBook(root string, assetsFile string){
	orgRootPath = root
	file, err  := os.ReadFile(path.Join( root, assetsFile))
	// defer file.Close()
	// var broker map[string]interface{}
	
	if err == nil {
		if err := json.Unmarshal(file, &content); err != nil {
			log.Println(err)
		} else {
			// ConstructResourceHierarchy()
			WalkThroughOrgTreeResources()
			
		}
	} else {
		log.Fatalln("Error: ", err)
	}
	defer mappings.UpdateUnsupportedResourceMap()
}


func WalkThroughOrgTreeResources(){
	
	
	// mux := new(sync.RWMutex)
	mapPaths := &map[string][2]any{}
	for index, con := range content {
		resourceType := &TFResourceType{ Mirror: &con }
		control.Store(index, []any{false, nil})

		// if !strings.EqualFold(con.AssetType, "compute.googleapis.com/Instance"){
		// 	continue
		// }
		
		muw := new(sync.WaitGroup)
		// Check if resource type is in the allowed list.
		muw.Add(3)
		go resourceType.walkOrgTreeOnBackground(OrgPrimaryResourceTypes, index, muw)
		go resourceType.makeUserFriendlyPath(index, muw, *mapPaths)
		go resourceType.ConstructResourceHierarchy(index, muw)
		muw.Wait()
	}	
}


func (tfr *TFResourceType) walkOrgTreeOnBackground(rstypes []string,
	index int, 
	// resourceType *TFResourceType, 
	muw *sync.WaitGroup) {
	
	defer muw.Done()

	// log.Print("HERE WALK")
		
	if name, err := url.QueryUnescape(tfr.Mirror.PickResourceName()); err == nil {
		tfr.Name = name
	}
	if parent, err := url.QueryUnescape(tfr.Mirror.Ancestors[0]); err == nil {
		tfr.Parent = parent
	}
	
	go tfr.SetRequiredFields()
	
	// mux.Unlock()
	if index == len(content) - 1 {
		return
	} else if strings.Contains(tfr.Mirror.Ancestors[0], "organizations") && len(tfr.Mirror.Ancestors) == 1 {
		return
	} else {
		for _, r := range rstypes {
			if tfr.Mirror.PickResourceType() == r {
				break
			}
		}
		return
	}
}


func (tfr *TFResourceType) makeUserFriendlyPath(
	index int,
	muw *sync.WaitGroup,
	mapPaths map[string][2]any){

	defer muw.Done()
	var pathTemp string
	pathKey := tfr.Mirror.Ancestors[0]
	work := func(c [2]any, a string){
		switch c {
			case [2]any{} : {
				nameT := collectors.FetchResourceNameById(a)
				
				if r := recover(); r != nil {
					break
				}
				pathTemp = path.Join(nameT.(string), pathTemp)
				c = [2]any{true, pathTemp}	
			}
			default:
				pathTemp = path.Join(c[1].(string), pathTemp)
				c = [2]any{true, pathTemp}
				// continue
		}
	}
	
	// Only execute the step if path is not yet in the control at the 
	// specified index.
	contx, _ := control.Load(index)
	cx := contx.([]any)
	if !cx[0].(bool) && cx[1] == nil {
		ancestors := tfr.Mirror.Ancestors
		switch  tfr.Mirror.PickResourceType() {
		case "Project", "Folder", "Organization" : {
			for _, a := range ancestors {
				if !strings.Contains(a, "org") {continue}
				work(mapPaths[a], a)
			}
			pathTemp = path.Join(pathTemp, ".META", tfr.Mirror.PickServiceType())
		}
		default : {
			for _, a := range ancestors {
				// if !strings.Contains(a, "org") {continue}
				work(mapPaths[a], a)
			}
			if strings.Contains(ancestors[0], "org") {
				pathTemp = path.Join(pathTemp, ".META", tfr.Mirror.PickServiceType())
			} else {
				pathTemp = path.Join(pathTemp, ".ACTIVE", tfr.Mirror.PickServiceType())
			}
		}} 

		if pathTemp != "" {
			if _, ok:= control.Swap(index, []any{true, pathKey, pathTemp, tfr.Mirror.PickResourceType()}); !ok {
				log.Printf("STORED: %v", control)
			}
			
		}	
		// log.Print("HERE MAKE")	
	} 
}


func (rsc *TFResourceType) ConstructResourceHierarchy (index int,muw *sync.WaitGroup) {
	
	// mux.RLock()
	// defer mux.RUnlock()
	defer muw.Done()
	// resourceType := TFResourceType{}
	
	cm,_ := control.Load(index)
	
	for len(cm.([]any)) < 4 /*Minimum length of */ {
		cm,_ = control.Load(index)
	}
	cx  := cm.([]any)
	
	if rsc.CheckDirectoryExistsOrCreate(orgRootPath, cx[2].(string)) {
		// log.Printf("RH CONSTRUCT: %v --- DIR: %v => %v", index, cm[2], cm[3])
		if !rsc.CheckTFResourceTypeFileExistsOrCreate(orgRootPath, cx[2].(string), cx[3].(string)) {
			log.Panicf("FAILED TO CREATE THE RESOURCE TYPE %s", cx[2])
		}
		
		// for rsc.Path == "" && rsc.TfrType == "" && rsc.Mirror.IsSupported {
		// 	rsc = rsc.SetRequiredFields(orgRootPath, cm[2].(string), cm[3].(string))
		// 	// rsc.SetTFResourceFilePath()
		// }
		// log.Printf("HERE CONTRUCT: %v", rsc )
		if rsc.Mirror.IsSupported && rsc.GetTFResourceType() != "" {
			// go rsc.ConstructTFResourceTemplate()
			go rsc.AppendIntoTFResourceList()
		}
	}
}