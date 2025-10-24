package filefunc

import (
	"embed"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

)

func IsExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// CreateFolder creates a directory and all necessary parents.
func CreateFolder(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

// CreateFile creates or overwrites a file at the given path.
func CreateFile(path string) (*os.File, error) {
	return os.Create(path)
}

// DeleteFile removes the specified file.
func DeleteFile(path string) error {
	if err := os.Remove(path); err != nil {
		return err
	}
	log.Println("Deleted file:", path)
	return nil
}

// DeleteFolder_FR force removes a directory and all its contents.
func DeleteFolder_FR(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return err
	}
	log.Println("Deleted folder and content:", path)
	return nil
}

func GetFileInfo(path string) os.FileInfo {
	info, err := os.Stat(path)
	if err != nil {
		return nil
	}
	return info
}

// ReadFile reads the contents of a file.
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// GetFileList returns a list of files in a directory (recursively).
func GetFileList(dir string) []string {
	var files []string
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files
}

// GetFileListByExt returns a list of files in a directory with a specific extension (recursively).
func GetFileListByExt(dir, ext string) []string {
	var files []string
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && strings.HasSuffix(path, ext) {
			files = append(files, path)
		}
		return nil
	})
	return files
}

func WriteWebFSToDisk(folder string, wfs embed.FS) error {

	log.Println("Writing embed.FS to disk:", folder)
	err := fs.WalkDir(wfs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		targetPath := filepath.Join(folder, path)

		if d.IsDir() {
			if !IsExists(targetPath) {
				if err := CreateFolder(targetPath); err != nil {
					return err
				}
			}
		} else {
			if !IsExists(targetPath) {
				data, err := wfs.ReadFile(path)
				if err != nil {
					return err
				}
				if err := os.WriteFile(targetPath, data, 0644); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil

}
