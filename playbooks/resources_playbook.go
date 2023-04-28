package playbooks

import (
	"fmt"
	"log"
	"mappings"
	"os"
	"path"
	"strings"
	"sync"
	"text/template"
)


type TFResourceFactory interface {
	// MakeResource () string
	// DoResourceMetadata () string
	// MakeServiceType () string

	/*Check if resource catalog file's parent folders exist in the directory tree. 
	  For instance, a directory pattern e.g. [orgnode/[...folder_name]/project_name]/Service_name/
	  where Service_name corresponds to subdomain (e.g. compute.googleapis.com) part "[compute]".
	  Example 1:
	    /orgnode/ 		# org node
		|--Orgpolicy/ 	# service_name e.g orgpolicy.googleapis.com/Policy
		|  |--*.tf
	  Example 2:
	    /Security/ 		# org unit
		|--Orgpolicy/ 	# service_name e.g orgpolicy.googleapis.com/Policy
		|  |--*.tf
	*/
	CheckDirectoryExistsOrCreate(root string, relative string) bool // create the directory IF NOT EXISTS returning the PATH and BOOL

	/*Check if resource catalog file's exists
	 For instance, a resource catalog file correspond to resource (e.g. compute.googlepais.com/[Instance|Network]) part: "Instance or Network"
	 Example 1:
    	/Projects/
		|--CloudLabs/           # project_name
		|  |--Compute/          # service_name
		|  |  |--Instance.tf    # resource catalog as terraform configuration file *.tf
		|  |  |--Network.tf     # resource catalog as terraform configuration file *.tf
		|  |--Dns/              # service_name e.g dns.googlepais.com/Policy
		|  |  |--Policy.tf      # resource catalog as terraform configuration file *.tf
	*/
	CheckTFResourceTypeFileExistsOrCreate () bool

	/* Construct the TF Google resource and append into the TF Resource Type File
	   For instance, a TF resource construct will have "[name]" and "[parent]" fields provided by Asset listing output.
	   Hence, if we run $(gcloud asset list --[organization|folder|project] ID --filter "assetType = 'compute.googleapis.com/Instance'" --limit 1 --format json)"
	   Example output: 
	   # assets.json
	   [
        {
          "ancestors": [
	          "projects/*******",
	          "folders/********",
	          "organizations/********"
          ],
          "assetType": "compute.googleapis.com/Instance",
          "name": "//compute.googleapis.com/projects/[project_id]/zones/europe-central2-a/instances/worker",
        },
		....]
		# Instance.tf
		...
		[truncated]
		resource "google_compute_instance" "worker" {
			name = "projects/[project_id]/zones/europe-central2-a/instances/worker" #(REQUIRED)
			parent = "projects/*******" #(optional)
		}
	*/
	ConstructTFGCPResource() string

	/* Reconstruct TF Google resource based on the import state values
	*/

	ConstructTFResourceTemplate () string

	IsSupportedTFResource(tftype string, gtype string) bool
}


type TFResourceType struct {
	Parent string     
	Name string 
	RequiredFields []string // map[string]interface{}
	TfrType string
	TfrLabel string
	Path string
	Type string
	Mirror *Asset
}

type Ancestors []string

type Asset struct {
	Ancestors Ancestors `json:"ancestors"`
	AssetType string `json:"assetType"`
	Name string `json:"name"`
	// updateTime string `json:update_time`
	IsSupported bool `json:"supported"`
}

type ServiceType struct{}

var KindResourceCounter map[string]int 
var file *os.File; 
var fileOerr error

var ancestors string

var ListOfAllTFResourcesPerPath sync.Map



func init(){
	file, fileOerr = os.OpenFile("./playbooks/catalog.tf", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	KindResourceCounter = map[string]int{}
	ListOfAllTFResourcesPerPath = sync.Map{}
}

/**
* Let resource be asset of type e.g. compute.googlepais.com/Instance, therefore:
* 	1. Let name be string formatted as "projects/{project_id}/zones/{zone_name}/instances/[name]", and
*   2. Let required_fields be a map of key=value, where key is a required fields different to [name and parent]. And
*   3. Let parent be (optional) e.g. projects/{project_id}, so
* So, there should be a parent folder where:
	1. Folder name is "Compute", so it should have a file with:
	2. File name is "Instance.tf"
*/

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
	
	return strings.ToLower(sType)  
}


func (a *Asset) SetSupported(s bool){
	switch {
	case s: 
		a.IsSupported = s
	default: 
		log.Printf("UNSUPPORTED RSTYPE: %v", a.AssetType)
	}	
}


func (tfr *TFResourceType) SetTFResourceType(s string){
	tfr.TfrType = s
}


func (tfr *TFResourceType) GetTFResourceType() string {
	return tfr.TfrType
}


func (tfr *TFResourceType) CheckDirectoryExistsOrCreate(p ...string) bool {
	// TODO: Concatenate the root and relative path  then test if exists
	// TODO: Create the tree IF NOT EXIST
	base := path.Join(p...)
	if _, err := os.ReadDir(base); err != nil {
		if err  := os.MkdirAll(base, 0777); err == nil {
			log.Printf("Creating directory tree: %v", base)
		} else {
			return false
		}
	}
	return true
}


func (tfr *TFResourceType) CheckTFResourceTypeFileExistsOrCreate (p ...string) bool {
	//TODO: For each asset type create a file in the directory tree
	filepath := path.Join(p...) + ".tf"
	if file, err := os.OpenFile(filepath, os.O_CREATE|os.O_SYNC, 0777); err == nil {
		tfr.SetTFResourceFilePath(filepath)
		return file.Close() == nil
	}
	return false
}


func (tfr *TFResourceType) SetRequiredFields (p ...string) *TFResourceType {
	tfr.TfrType = mappings.MatchTFRNameByCloudRSType(tfr.Mirror.AssetType).(string)

	if tfr.TfrType =="" {
		// log.Printf("REQUIRED.....: %v", tfr.Mirror.AssetType)
		return nil
	}
	tfr.Mirror.SetSupported(true)

	if rq, err := mappings.PickResourceRequiredFieldsByTFRName(tfr.TfrType); err == nil {
		tfr.RequiredFields = rq
	} else {
		return nil
	}
	// go tfr.SetTFResourceType(tfrType)
	KindResourceCounter[tfr.TfrType] += 1
	count := KindResourceCounter[tfr.TfrType]
	
	tfr.TfrLabel = fmt.Sprint(tfr.TfrType,"_",count)
	// go mappings.TFImportState(tfrType.(string), tfr.TfrLabel, tfr.Name)
	// log.Printf("RSTYPE: %v --- TF: %v", tfr.Mirror.AssetType, tfr.TfrType)
	go tfr.SetTFResourceFilePath(p...)
	return tfr
}


func (tfr *TFResourceType) ConstructTFResourceTemplate () any {

	// tfrType := mappings.MatchTFRNameByCloudRSType()
	const form = `{{block "main" .}}{{"\n"}}
	locals {
		states    = jsondecode(file("{{.TfrValuesFile}}"))["resources"]
		instances = [for resource in local.states : merge(resource["instances"]...) if resource["name"] == "{{.TfrLabel}}"]
	  }

	resource "{{.TfrType}}" "{{.TfrLabel}}"{
		for_each     = { for k, vm in local.instances : vm["index_key"] => vm["attributes"] }
	{{$name := .Name}}{{$parent := .Parent}}{{$hasname := false}}{{$hasparent := false}}
	# required fields
    {{range $i, $req := .RequiredFields}}{{$isname := ne $req "name"}}{{$isparent := ne $req "parent"}}{{if and $isname $isparent}}{{$req}} = nil{{"\n\t"}}{{else if $isname}}{{$hasname = $isname}}{{else if $isparent}}{{$hasparent = $isparent}}{{end}}{{end}}
    {{if $hasname}}name = "{{$name}}"{{else}}name = "{{$name}}"{{end}}
    {{if $hasparent}}parent = "{{$parent}}"{{end}}	{{"\n"}}}{{"\n\n"}}{{end}}`

	resource, err := template.New(strings.ToLower(tfr.TfrType)).Parse(form)
	// resource, err := resource.ParseFiles("./playbooks/.meta_template")
	if err != nil {
		log.Printf("FAILED TEMPLATE PARSE: %v", err)
	} 

	
	if fileOerr != nil {
		log.Printf("ERROR OPENING FILE: %v", err)
	}

	// for tfr.TfrType == "" {} 
	file, err := os.OpenFile(tfr.Path, os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		log.Panicf("FAILED TO OPEN FILE: %v TO APPEND TEMPLATE. %v", tfr.Path, err)
	}
	defer file.Close()
	log.Printf("ADDING TF RESOURCE %v --- INTO %v", tfr.TfrType, tfr.Path)
	if err := resource.Execute(file, tfr); err != nil {
		log.Panic(err)
	}

	return "Done!"
}


func (tfr *TFResourceType) ConstructTFResourceLocalValuesFile() any {
	// Construct TF Local Variable Template
	// Get the resource file path and add a JSON File
	rspath := path.Dir(tfr.Path)
	if file, err := os.OpenFile(path.Join(rspath, tfr.Mirror.PickResourceName()), os.O_CREATE|os.O_APPEND, 0777); err == nil {
		os.WriteFile(rspath, []byte{}, 0777)
		file.Close()
	}

	const local = `{{block "main .}}
	local {
		{{range .}}
	}{{"\n\n"}}`


	return nil
}


func (tfr *TFResourceType) SetTFResourceFilePath(p ...string){
	tfr.Path = path.Join(p...)
}


func (tfr *TFResourceType) AppendIntoTFResourceList(){
	
	if val, ok := ListOfAllTFResourcesPerPath.Load(tfr.Path); ok {
		val = append(val.([]any), tfr)
		ListOfAllTFResourcesPerPath.Swap(tfr.Path, val)
	} else if val == nil {
		ListOfAllTFResourcesPerPath.Store(tfr.Path, []any{})
		// log.Printf("FAILED TO STORE")
	}
}