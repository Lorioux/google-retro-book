package playbooks

import "google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/source"

type TFImportModulesFactory interface{
	/*Check  TF importing resource meta file EXISTS OR CREATE NEW.
    For instance, a TF importing resource meta file contains a reference of service folders as modules"
    Example 1:
      /Projects/
      |--CloudLabs/           # project_name
      |  |--meta_modules.tf 	# Terraform modules for imports
    Example 2:
      /[orgnode | orgunit]/   # org node
      |--meta_modules.tf      # Terraform importing resource modules
    */
	CheckTFMetaModulesFileExistsOrCreate () bool

	/*Construct a meta module for importing resource type. For instance. meta modules have source field
    pointing to the resources relative path.
    Example 1:
	    /Projects/
	    |--CloudLabs/       
	    |  |--Compute/      
	    |  |  |--Instance.tf
         |--meta_modules.tf   # content example
             [truncated]
             module "compute" {
                source = "./compute"
             }
             [truncated]
	*/
  ConstructTFMetaModule() string 
}


/*Also it can have a `resources block, optional` as:
 resources = list(object({
  resource_type = string  # e.g. google_compute_instance, google_compue_disk
  identifier    = string  # e.g. workernode, workernode-disk
  resource_id   = string  # e.g. projects/-/zones/-/instances/workernode
 }))   
*/
type TFMetaModule struct {
  source string
  //resources []interface{}
}


func (t TFMetaModule) ConstructTFMetaModule() string  {
  t.source = ""
  return ""
}