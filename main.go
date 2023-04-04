package main

import (
	"collectors"
	"log"
	"mappings"
)

func main() {
	// c := make(chan []map[string]collectors.Project)
	// core.Run(c)
	// core.CreateFolders(nil)
	// defer Test(c)
	// go test(collectors.AllProjectsAssets)
	// log.Print("Hello")
	// time.Sleep(10 * time.Second)
	// log.Print(collectors.AllProjectsAssets)
	// mappings.CallProvider()
	// mappings.ReadFileToReverse("mappings/data_mapping.json", "mappings/data_reverse.json")
	mappings.ReadFileDoReverse("mappings/resource_mapping.json", "mappings/resource_reverse.json")
}

func Test(a chan []map[string]collectors.Project){
	result := <- a
	log.Println(result)
}