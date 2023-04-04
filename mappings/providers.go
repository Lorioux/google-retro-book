package mappings

import (
	goo "google"
	"log"
)

var (
	ProviderRequiredParams []string
	ResourceRequiredParams map[string][]string
)

func CallProvider(){
	// provider := google.NewComputeEngineCredentials()
	provider := goo.Provider()
	ResourceRequiredParams = map[string][]string{}

	for param := range provider.Schema {
		if !provider.Schema[param].Optional {
			ProviderRequiredParams = append(ProviderRequiredParams, param)
		}
		
	}
	log.Printf("[PROVIDER PARAMS] : %v", ProviderRequiredParams)

	for resource_name, resource_schema := range provider.ResourcesMap {
		if (resource_name != "google_compute_instance") /*==  (resource_name != "google_storage_bucket_iam_binding")*/ {
			// log.Println(resource_name)
			continue
		}
		for param_name, param := range resource_schema.Schema {
			if param.Required {
				ResourceRequiredParams[resource_name] = append(ResourceRequiredParams[resource_name], param_name)
			}
			
		}
	}
	
	for rname, rparams := range ResourceRequiredParams {
		log.Printf("[%v] : %v", rname, rparams)
	}

}