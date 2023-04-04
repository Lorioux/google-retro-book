package playbooks

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

	/*
	*/
}