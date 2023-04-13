package playbooks

import (
	"fmt"
	"log"
	"mappings"
	"os"
	"path"
	"strings"
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
}


type TFResourceType struct {
	Parent string     
	Name string 
	RequiredFields []string // map[string]interface{}
	TfrType string
	TfrLabel string
}

var KindResourceCounter = map[string]int{}

var ancestors string

/**
* Let resource be asset of type e.g. compute.googlepais.com/Instance, therefore:
* 	1. Let name be string formatted as "projects/{project_id}/zones/{zone_name}/instances/[name]", and
*   2. Let required_fields be a map of key=value, where key is a required fields different to [name and parent]. And
*   3. Let parent be (optional) e.g. projects/{project_id}, so
* So, there should be a parent folder where:
	1. Folder name is "Compute", so it should have a file with:
	2. File name is "Instance.tf"
*/

func (tfr *TFResourceType) CheckDirectoryExistsOrCreate(root string, relative string) bool {
	// TODO: Concatenate the root and relative path  then test if exists
	// TODO: Create the tree IF NOT EXIST
	base := path.Join(root, relative)
	if _, err := os.ReadDir(base); err != nil {
		// log.Printf("Null dir: %v --- %v", err, entries)
		if err  := os.MkdirAll(base, 0777); err == nil {
			log.Print("Created directory tree")
		}
	}
	return true
}

func (tfr *TFResourceType) SetRequiredFields (assetType string) {
	tfrType := mappings.MatchTFRNameByCloudRSType(assetType)
	if tfrType == nil {
		return
	}
	if rq, err := mappings.PickResourceRequiredFieldsByTFRName(tfrType.(string)); err == nil {
		tfr.RequiredFields = rq
	} else {
		return
	}
	KindResourceCounter[tfrType.(string)] += 1
	// log.Printf("REQUIRED: %v", tfr.required_fields)
	// if strings.EqualFold(tfrType.(string), "google_compute_instance") {
		count := KindResourceCounter[tfrType.(string)]
		tfr.TfrType = tfrType.(string)
		tfr.TfrLabel = fmt.Sprint(tfrType,"_",count)
		tfr.ConstructTFResourceTemplate(tfrType.(string))
	// }
}



func (tfr TFResourceType) ConstructTFResourceTemplate (tfrName string) string {
	const form = `{{block "main" .}}
resource "{{.TfrType}}" "{{.TfrLabel}}"{
	{{$name := .Name}}{{$parent := .Parent}}{{$hasname := false}}{{$hasparent := false}}
	# required fields
    {{range $i, $req := .RequiredFields}}{{$isname := ne $req "name"}}{{$isparent := ne $req "parent"}}{{if and $isname $isparent}}{{$req}} = nil{{"\n"}}{{else if $isname}}{{$hasname = $isname}}{{else if $isparent}}{{$hasparent = $isparent}}{{end}}{{end}}
    {{if $hasname}}name = "{{$name}}"{{else}}{{end}}
    {{if $hasparent}}parent = "{{$parent}}"{{end}}	
}{{"\n\n"}}{{end}}`

	resource, err := template.New(strings.ToLower(tfrName)).Parse(form)
	// resource, err := resource.ParseFiles("./playbooks/.meta_template")
	if err != nil {
		log.Printf("FAILED TEMPLATE PARSE: %v", err)
	} 
	// log.Printf("HERE: %v", tfr)
	file, err := os.OpenFile("./playbooks/templates.tf", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		log.Printf("ERROR OPENING FILE: %v", err)
	}
	// file = os.Stdin
	// var data []byte
	if err:= resource.Execute(file, tfr); err != nil {
		log.Panic(err)
	}
	// if _, err = file.Read(data); err == nil {
	// 	os.WriteFile("./playbooks/templates.tf", data, 0777)
		
	// } else {
	// 	log.Printf("EMPTY %v", err)
	// }

	return "Done!"
}

