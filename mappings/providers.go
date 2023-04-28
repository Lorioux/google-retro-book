package mappings

import (
	"errors"
	"fmt"
	meta "google"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ProviderRequiredParams []string
	ResourceRequiredParams map[string][]string
	Provider *schema.Provider
	ResourceMap map[string]*schema.Resource
	 
)
var Cfg meta.Config

func init(){
	ProviderRequiredParams = []string{}
	ResourceRequiredParams = map[string][]string{}
	Provider = meta.Provider()
	ResourceMap = Provider.ResourcesMap
}

func ExecuteCallProvider(){
	// provider := google.NewComputeEngineCredentials()
	
	// ResourceRequiredParams = map[string][]string{}

	for param := range Provider.Schema {
		if !Provider.Schema[param].Optional {
			ProviderRequiredParams = append(ProviderRequiredParams, param)
		}	
	}
	// log.Printf("[PROVIDER PARAMS] : %v", ProviderRequiredParams)

	// for resource_name, resource_schema := range Provider.ResourcesMap {
	// 	if (resource_name != "google_compute_instance") /*==  (resource_name != "google_storage_bucket_iam_binding")*/ {
	// 		// log.Println(resource_name)
	// 		continue
	// 	}
	// 	for param_name, param := range resource_schema.Schema {
	// 		if param.Required {
	// 			ResourceRequiredParams[resource_name] = append(ResourceRequiredParams[resource_name], param_name)
	// 		}
			
	// 	}
	// }
	
	// for rname, rparams := range ResourceRequiredParams {
	// 	log.Printf("[%v] : %v", rname, rparams)
	// }

}


func PickResourceRequiredFieldsByTFRName( tfrName string) ([]string, error) {
	// tfrName := MatchTFRNameByCloudRSType()
	
	resource := ResourceMap[tfrName]
	var fields []string
	// time.Sleep(1 * time.Second)
	if resource == nil {
		return nil, errors.New("NIL")
	}
	for param_name, param := range resource.Schema {
		if param.Required {
			fields = append(fields, param_name)
		} else if strings.EqualFold(param_name, "parent") && param.Optional{
			fields = append(fields, param_name)
		}
	}
	if fields == nil {
		return nil, nil
	}
	// log.Printf("[TF_RESOURCE_REF]: %v ---- [REQUIRED_FIELDS] : %v", tfrName, fields)
	ResourceRequiredParams[tfrName] = fields
	return ResourceRequiredParams[tfrName], nil
}


func TFImportState(tfrType string, label string, id string){
	// resource := ResourceMap[tfrType]
	// cfg := &meta.Config{
	// 	ImpersonateServiceAccount: "cw-prod-resman-000@cw-prod-iac-core-000.iam.gserviceaccount.com",
	// }
	
	// data := resource.Data(nil)
	// data.SetId(id)
	
	// Provider.SetMeta(cfg)
	// info := &terraform.InstanceInfo{Id: data.Id(), Type: tfrType}
	// if state, _ := Provider.ImportState(collectors.GetContext(), info, data.Id()); state != nil {
	// 	for _, s := range state {
	// 		log.Print(s)
	// 	}
	// }
	// resource_addr := fmt.Sprintf("%s.%s \"%s\"", tfrType, label, id)
	
	// tfInitCmd := exec.Cmd{
	// 	Path: "/usr/bin/terraform",
	// 	Args: []string{"terraform", "init"},
	// 	Dir: fmt.Sprint(os.Getenv("PWD"), "/playbooks"),
	// 	Stdout: os.Stdout,
	// } 
	tfImportCmd := exec.Cmd{
		Path: "/usr/bin/terraform",
		Args: []string{"terraform", "import"},
		Dir: fmt.Sprint(os.Args[1], "/playbooks"),
		Stdout: os.Stderr,
	}
	
	// if err := tfInitCmd.Run(); err != nil {
	// 	log.Println(err)
	// } else {
	// 	log.Printf("PROCESS ID: %d",tfInitCmd.Process.Pid)
	// 	tfInitCmd.Wait()
	// } 
	if err := tfImportCmd.Run(); err != nil{
		log.Print(tfImportCmd.Dir)
		log.Println(err)
	} else {
		
		log.Printf("PROCESS ID: %d",tfImportCmd.Process.Pid)
		// log.Printf("PROCESS ID: %v", proc.)
		
		tfImportCmd.Wait()
	}


	// if instances, err := resource.Importer.State(data, cfg); err == nil {

	// 	for _,i := range instances {
	// 		log.Print(i.State())
			
	// 	}
	// 	// for _, r := range Provider.Resources(){
	// 	// 	log.Println(r.Name, r.Importable, r.SchemaAvailable)
	// 	// 	break	
	// 	// }
	// }
	
	// cfg.NewAppEngineClient(cfg.UserAgent).Projects
	// 
	
	// log.Println(resource.Exists(data, &meta.Config{
	// 	ImpersonateServiceAccount: "cw-prod-resman-000@cw-prod-iac-core-000.iam.gserviceaccount.com",
	// }))
	// resource.TestResourceData()
}

