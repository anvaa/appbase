package filefunc

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

)

func IsExists(path string) bool {
	// if folder/file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// Create folder
func CreateFolder(path string) error {
	// create folder
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	// log.Println("Created folder:", path)
	return nil
}

// Create file
func CreateFile(path string) (*os.File, error) {
	// create or overwrite file
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	// log.Println("Created file:", path)
	return file, nil
}

// Delete file
func DeleteFile(filepath string) error {

	_, err := os.Stat(filepath)
	os.IsExist(err)
	if err != nil {
		return err
	}
	os.Remove(filepath)
	log.Println("Deleted file:", filepath)
	return nil
}

// Delete folder with content
func DeleteFolder_FR(path string) error {
	// force delete folder and all sub folders with content
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}
	log.Println("Deleted folder and content:", path)
	return nil
}

func ReadFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetFileList returns a list of files in a directory
func GetFileList(dir string) []string {
	var files []string
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files
}

// GetFileListByExt returns a list of files in a directory with a specific extension
func GetFileListByExt(dir string, ext string) []string {
	var files []string
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, ext) {
			files = append(files, path)
		}
		return nil
	})
	return files
}

func intToString(i int) string {
	// convert int to string
	s := fmt.Sprintf("%d", i)
	// remove trailing zeros
	s = strings.TrimRight(s, "0")
	// remove trailing dot
	s = strings.TrimRight(s, ".")
	return s
}

func floatToString(f float64) string {
	// convert float to string
	s := fmt.Sprintf("%f", f)
	// remove trailing zeros
	s = strings.TrimRight(s, "0")
	// remove trailing dot
	s = strings.TrimRight(s, ".")
	return s
}

func remCharForCSV(str string) string {
	// remove char from string
	remc := ","
	// replace with space
	str = strings.ReplaceAll(str, remc, " ")
	// remove double quotes
	str = strings.ReplaceAll(str, "\"", "")

	// remove new line
	str = strings.ReplaceAll(str, "\n", " ")
	str = strings.ReplaceAll(str, "\r", " ")
	// remove tab
	str = strings.ReplaceAll(str, "\t", " ")
	// remove multiple spaces
	str = strings.ReplaceAll(str, "  ", " ")

	return str
}

func WriteWebFSToDisk(folder string, wfs embed.FS) error {

	log.Println("Writing embed.FS to disk:", folder)
	// loop thru static and write folders and files to static on disk
	err := fs.WalkDir(wfs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			// create folder
			folder := folder + "/" + path
			if !IsExists(folder) {
				CreateFolder(folder)
			}
		} else {
			// create file
			folder := folder + "/" + path
			if !IsExists(folder) {
				file, err := CreateFile(folder)
				if err != nil {
					return err
				}
				data, err := wfs.ReadFile(path)
				if err != nil {
					return err
				}
				file.Write(data)
				file.Close()
			}
		}

		return nil
	},
	)
	if err != nil {
		return err
	}

	return nil

}
