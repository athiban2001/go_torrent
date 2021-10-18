package defs

import (
	"fmt"
)

type File struct {
	Length int64    `json:"length,omitempty"`
	MD5Sum string   `json:"md5sum,omitempty"`
	Path   []string `json:"path"`
}

func getFiles(filesData []interface{}) ([]*File, error) {
	var (
		files []*File = []*File{}
		file  *File
		err   error
	)

	for _, fileData := range filesData {
		if file, err = getFile(fileData); err != nil {
			return nil, err
		}

		files = append(files, file)
	}

	return files, nil
}

func getFile(fileData interface{}) (*File, error) {
	fileMap, _ := fileData.(map[string]interface{})
	pathsData, pathsDataOK := fileMap["path"].([]interface{})
	length, lengthOK := fileMap["length"].(int64)
	md5Sum, _ := fileMap["md5sum"].(string)

	if !pathsDataOK || !lengthOK {
		return nil, fmt.Errorf("invalid file data : path,length keys not found")
	}

	file := &File{}
	file.Length = length
	file.MD5Sum = md5Sum
	file.Path = []string{}
	for _, pathData := range pathsData {
		path, pathOK := pathData.(string)
		if pathOK {
			file.Path = append(file.Path, path)
		}
	}

	return file, nil
}
