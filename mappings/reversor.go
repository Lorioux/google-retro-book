package mappings

import (
	"encoding/json"
	"fmt"
	"os"
)

//Pick a json file then reverse the key and value
//Save the result as json file
func ReadFileDoReverse(scfile string, dsfile string) {
    // Read the JSON file.
    data, err := os.ReadFile(scfile)
    if err != nil {
        panic(err)
    }

    // Decode the JSON into a map.
    var mapObj map[string]interface{}
    err = json.Unmarshal(data, &mapObj)
    if err != nil {
        panic(err)
    }

    // Create a new map to store the swapped keys and values.
    newMap := make(map[string]string)

    // Iterate over the map and swap the keys and values.
    for _, _map := range mapObj["maps"].([]interface{})[0:] {
        
        for key, value := range _map.(map[string]interface{}) {
            newMap[value.(string)] = key
        }
		// fmt.Println(key, value)
    }

    // Print the swapped map.
    // fmt.Println(newMap)
	mapObj["reverse"] = newMap
	data, err = json.Marshal(mapObj)
	if err != nil {
		fmt.Print(err)
	}
    dsfile = scfile
	os.WriteFile(dsfile, data, 0777)
}