package playbooks





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
	CheckDirectoryExists() (string, bool) // create the directory IF NOT EXISTS returning the PATH and BOOL

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
	CheckTFResourceTypeFileExists () bool

	/*Check if terraform resource import meta mappings file's exists
	 For instance, a resource catalog file correspond to resource (e.g. compute.googlepais.com/[Instance|Network]) part: "Instance or Network"
	 Example 1:
		/Projects/
		|--CloudLabs/       # project_name
		|  |--meta_mappings # TF importing meta mappings file 
	*/
	CheckTFImportMetaFileExists () bool

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

}




type TFResourceType struct {
	parent string     
	name string 
	required_fields map[string]interface{}
}

/**
* Let resource be asset of type e.g. compute.googlepais.com/Instance, therefore:
* 	1. Let name be string formatted as "projects/{project_id}/zones/{zone_name}/instances/[name]"
*/