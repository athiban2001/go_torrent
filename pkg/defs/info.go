package defs

import (
	"fmt"
)

type Info struct {
	Name        string  `json:"name,omitempty"`
	PieceLength int64   `json:"piece-length,omitempty"`
	Pieces      string  `json:"pieces,omitempty"`
	Private     bool    `json:"private,omitempty"`
	Files       []*File `json:"files"`
}

func GetInfo(infoData map[string]interface{}) (*Info, error) {
	name, nameOK := infoData["name"].(string)
	length, lengthOK := infoData["length"].(int64)
	md5sum, _ := infoData["md5sum"].(string)
	piecesLength, piecesLengthOK := infoData["piece length"].(int64)
	pieces, piecesOK := infoData["pieces"].(string)
	private, _ := infoData["private"].(bool)
	filesData, filesDataOK := infoData["files"].([]interface{})

	if lengthOK {
		filesData = append(filesData, map[string]interface{}{
			"path":   []interface{}{name},
			"length": length,
			"md5sum": md5sum,
		})
	}
	if !nameOK || !piecesLengthOK || !piecesOK || (!lengthOK && !filesDataOK) {
		return nil, fmt.Errorf("invalid field : name,piecesLength,pieces,length,files in info")
	}
	files, err := getFiles(filesData)
	if err != nil {
		return nil, err
	}

	info := &Info{}
	info.Name = name
	info.Pieces = pieces
	info.PieceLength = piecesLength
	info.Private = private
	info.Files = files
	return info, nil
}
