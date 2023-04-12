package mappings

import (
	"errors"
	meta "google"
	"strings"
)

var (
	ProviderRequiredParams = []string{}
	ResourceRequiredParams = map[string][]string{}
	Provider = meta.Provider()
	ResourceMap = Provider.ResourcesMap
)

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

