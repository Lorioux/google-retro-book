package core

import (
	// "gopkg.in/yaml.v3";

	coll "collectors"
	"encoding/json"
	"log"
	"os"
	"os/exec"
)

var results string

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Folders struct {
	open_c  string
	content string
	close_c string

	output string
}

func resolve(f *Folders) string {
	f.output = "" + f.open_c + f.content + f.close_c
	return f.output
}

func Run(c chan []map[string]coll.Project) {
	coll.GetProjects()
	coll.GetAssetsByProjectId(c)

	flds := PrepareListOfFolders("folders.txt")
	if flds != nil {
		CreateFolders(flds)
	}

}

// get the list of folders
func PrepareListOfFolders(filepath string) *Folders {
	content, err := os.ReadFile("./folders/" + filepath)

	check(err)

	var size = len(content) - 1

	flds := &Folders{
		open_c:  "[",
		content: "",
		close_c: "]",
	}

	for e := 0; e <= size; e++ {

		var c = string(content[e])

		if c == "]" && e < size-1 {
			flds.content = flds.content + c + ","
		} else {
			flds.content = flds.content + c
		}

	}

	return flds
}

// create folders as per list
func CreateFolders(flds *Folders) {

	results = resolve(flds)
	// write to output to a file
	if err := os.WriteFile("./folders/folders.json", []byte(results), 0777); err != nil {
		panic(err)
	}

	re := jsonWalker(results)

	for i := 0; i < len(re); i++ {
		if err := re[i]; err != nil {
			var folder = re[i][0]
			fpath := "./folders/"

			// check folder exists
			if CheckFolderExists(&fpath, folder) {
				// fmt.Printf("Exists: %v\n", folder)
				continue
			}

			cmd := exec.Command("mkdir", string(fpath+folder))

			// create folder then return nil
			if err := cmd.Run(); err != nil {
				log.Printf("Command exited with error: %v", err)
				panic(err)
			}

		}

	}
}

// check folders exists
func CheckFolderExists(Fpath *string, Fname string) bool {
	entries, err := os.ReadDir(*Fpath)
	if err != nil {
		log.Panic(err)
		return true
	}

	for i := 0; i < len(entries); i++ {
		if entries[i].IsDir() && entries[i].Name() == Fname {
			return true
		}
	}
	return false
}

type ListOfFolders struct {
	content [][]string
}

func jsonWalker(s string) [][]string {

	f := &ListOfFolders{
		content: nil,
	}

	for {

		if err := json.Unmarshal([]byte(s), &f.content); err == nil {
			break
		} else if err != nil {
			log.Fatal(err)
		}
	}

	return f.content
}
