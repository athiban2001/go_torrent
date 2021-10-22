package torrent

import (
	"crypto/sha1"
	"fmt"
	"net/url"
	"strconv"

	"github.com/athiban2001/go_torrent/pkg/gt"
)

func GetBencodedHash(info *gt.Info) string {
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

func getInfo(infoData map[string]interface{}) (*gt.Info, error) {
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

	info := &gt.Info{}
	info.Name = name
	info.Pieces = pieces
	info.PieceLength = piecesLength
	info.Private = private
	info.Files = files
	return info, nil
}

func getSize(files []*gt.File) int64 {
	size := int64(0)
	for _, file := range files {
		size += file.Length
	}

	return size
}

func getFiles(filesData []interface{}) ([]*gt.File, error) {
	var (
		files []*gt.File = []*gt.File{}
		file  *gt.File
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

func getFile(fileData interface{}) (*gt.File, error) {
	fileMap, _ := fileData.(map[string]interface{})
	pathsData, pathsDataOK := fileMap["path"].([]interface{})
	length, lengthOK := fileMap["length"].(int64)
	md5Sum, _ := fileMap["md5sum"].(string)

	if !pathsDataOK || !lengthOK {
		return nil, fmt.Errorf("invalid file data : path,length keys not found")
	}

	file := &gt.File{}
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
