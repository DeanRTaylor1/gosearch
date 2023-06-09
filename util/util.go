package util

import (
	"encoding/json"
	"fmt"

	"log"
	"os"
)

type IndexedData struct {
	URL     string
	Content string // Or any other data structure used for storing indexed content
}

// Utility function, deprecated
func JSONToFile(j []byte, filename string) {
	fmt.Println("j length:", len(j)) // debugging line
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	l, err := f.Write(j)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

// Utility function, deprecated
func MapToJSON(m map[string]string, createFile bool, filename string) string {
	// Convert map[string]string to map[string]interface{}
	mi := make(map[string]interface{}, len(m))
	for k, v := range m {
		mi[k] = v
	}

	return MapToJSONGeneric(mi, createFile, filename)
}

// Utility function, deprecated
func MapToJSONGeneric(m map[string]interface{}, createFile bool, filename string) string {
	if len(m) == 0 {
		fmt.Println("map is empty")
		return ""
	}

	b, err := json.Marshal(m)
	if err != nil {
		fmt.Println("error:", err)
		return ""
	}
	if createFile {
		JSONToFile(b, filename)
	}
	return string(b)
}

// This function is used to find all the pre-existing indexes
func GetCurrentAvailableModelDirectories() []string {
	files, err := os.ReadDir("./indexes")
	if err != nil {
		log.Println(err)
		err := os.Mkdir("./indexes", os.FileMode(0777))
		if err != nil {
			log.Println(err)
		}
	}

	directories := []string{}
	for _, f := range files {
		if f.IsDir() {
			if f.Name() == "src" || f.Name() == ".git" {
				continue
			}
			directories = append(directories, f.Name())
		}
	}

	return directories
}

// This function checks that a dir exists
func CheckDirIsValid(dirName string) (bool, error) {
	_, err := os.Stat("./" + dirName)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil // Directory does not exist
		}
		return false, err // Some other error occurred
	}
	return true, nil // Directory exists
}

// This function is used to calculate the directory length or the total number of documents within a directory.
// Deprecated
func GetDirLength(dirName string) int {
	files, err := os.ReadDir("./" + dirName)
	if err != nil {
		log.Fatal(err)
	}

	return len(files)
}

// Const variables to enable terminal colors
const (
	TerminalReset  = "\033[0m"
	TerminalRed    = "\033[31m"
	TerminalGreen  = "\033[32m"
	TerminalYellow = "\033[33m"
	TerminalBlue   = "\033[34m"
	TerminalPurple = "\033[35m"
	TerminalCyan   = "\033[36m"
	TerminalWhite  = "\033[37m"
)
