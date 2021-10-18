package defs

import (
	"crypto/sha1"
	"fmt"
	"net/url"
	"strconv"
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

func (info *Info) GetHash() string {
	str := "d"

	if len(info.Files) == 1 {
		str += "6:lengthi" + strconv.Itoa(int(info.Files[0].Length)) + "e"
		if info.Files[0].MD5Sum != "" {
			str += "6:md5sum" + strconv.Itoa(len(info.Files[0].MD5Sum)) + ":" + info.Files[0].MD5Sum
		}
	} else {
		str += "5:filesl"
		for _, file := range info.Files {
			str += "d6:lengthi"
			str += strconv.Itoa(int(file.Length))
			str += "e"
			if file.MD5Sum != "" {
				str += "6:md5sum" + strconv.Itoa(len(file.MD5Sum)) + ":" + file.MD5Sum
			}
			str += "4:pathl"
			for _, p := range file.Path {
				str += strconv.Itoa(len(p)) + ":" + p
			}
			str += "ee"
		}
		str += "e"
	}
	str += "4:name" + strconv.Itoa(len(info.Name)) + ":" + info.Name
	str += "12:piece lengthi" + strconv.Itoa(int(info.PieceLength)) + "e"
	str += "6:pieces" + strconv.Itoa(len(info.Pieces)) + ":" + info.Pieces
	if info.Private {
		str += "7:privatei1e"
	}
	str += "e"

	hasher := sha1.New()
	hasher.Write([]byte(str))
	hash := hasher.Sum(nil)
	return url.QueryEscape(string(hash))
}
