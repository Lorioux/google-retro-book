package mappings

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
)

const CloudRTypeToTFRNameMappingFile = "resource_mapping.json"

const BASE_PATH = "./mappings"
var full_path = path.Join(BASE_PATH, CloudRTypeToTFRNameMappingFile)
var UnsupportedResourceMap map[string]interface{}
var resourceMapObj, err = readMappingsFile();


type ReverseMapType map[string]string

//Pick a json file then reverse the key and value
//Save the result as json file
func DoReverse(dsfile string) {
    
    if ; err != nil {
        panic(err)
    }

    // Create a new map to store the swapped keys and values.
    newMap := make(map[string]string)

    // Iterate over the map and swap the keys and values.
    for _, _map := range resourceMapObj["maps"].([]ReverseMapType)[0:] {
        
        for key, value := range _map {
            newMap[value] = key
        }
		// fmt.Println(key, value)
    }

    // Print the swapped map.
    // fmt.Println(newMap)
	resourceMapObj["reverse"] = newMap
	data, err := json.Marshal(resourceMapObj)
	if err != nil {
		fmt.Print(err)
	}
    // dsfile = scfile
	if err:= os.WriteFile(full_path, data, 0777); err != nil {
        fmt.Printf("There was an error: %v", err)
    }
    // fmt.Println("Done!")
}

func readMappingsFile() (map[string]interface{}, error) {
    // Read the JSON file.
    file, err := os.ReadFile(full_path)
    if err != nil {
        return nil, err
    }
    var mapObj map[string]interface{}
    if err = json.Unmarshal(file, &mapObj); err != nil {
        return nil, err
    }
    UnsupportedResourceMap = mapObj["unsupported"].(map[string]interface{})
    return mapObj, nil
}

func MatchTFRNameByCloudRSType(cloudRSType string) any {
    
    
    if  err != nil {
        panic(err)
    }
    // for k,v := range mapObj["reverse"] {
        
    // }
    v := resourceMapObj["reverse"].(map[string]interface{})[cloudRSType]; 

    if v == nil && UnsupportedResourceMap[cloudRSType] == nil {
        log.Printf("UNSUPPORTED RESOURCE TYPE: %v ----- %v", v , cloudRSType )
        UnsupportedResourceMap[cloudRSType] = "YES"
        return nil
    }
    return v
}

func UpdateUnsupportedResourceMap() {
    resourceMapObj["unsupported"] = UnsupportedResourceMap
    if data, err := json.Marshal(resourceMapObj); err == nil {
        os.WriteFile(full_path, data, 0777)
    } else {
        log.Fatalf("FAILED TO UPDATE UNSUPPORTED MAP: %v", err)
    }  
}