package torrent

import (
	"fmt"
	"time"
)

type Torrent struct {
	Announce     string     `json:"announce,omitempty"`
	AnnounceList [][]string `json:"announce-list,omitempty"`
	CreationDate time.Time  `json:"creation date,omitempty"`
	Comment      string     `json:"comment,omitempty"`
	CreatedBy    string     `json:"created by,omitempty"`
	Encoding     string     `json:"encoding,omitempty"`
	Info         *Info
}

type Info struct {
	PieceLength int64  `json:"piece-length,omitempty"`
	Pieces      string `json:"pieces,omitempty"`
	Private     bool   `json:"private,omitempty"`
	Files       []*File
}

type File struct {
	Length int64    `json:"length,omitempty"`
	MD5Sum string   `json:"md5sum,omitempty"`
	Path   []string `json:"path"`
}

func getTorrent(metaData map[string]interface{}) (*Torrent, error) {
	infoData, infoDataOK := metaData["info"].(map[string]interface{})
	announce, announceOK := metaData["announce"].(string)
	announceList, _ := metaData["announce-list"].([]interface{})
	epoch, _ := metaData["creation date"].(int64)
	comment, _ := metaData["comment"].(string)
	createdBy, _ := metaData["created by"].(string)
	encoding, _ := metaData["encoding"].(string)

	if !infoDataOK || !announceOK {
		return nil, fmt.Errorf("invalid field : info,announce")
	}
	info, err := getInfo(infoData)
	if err != nil {
		return nil, err
	}

	torrent := &Torrent{}
	torrent.CreationDate = time.Unix(epoch, 0)
	torrent.Comment = comment
	torrent.CreatedBy = createdBy
	torrent.Announce = announce
	torrent.AnnounceList = getAnnounceList(announceList)
	torrent.Encoding = encoding
	torrent.Info = info
	return torrent, nil
}

func getInfo(infoData map[string]interface{}) (*Info, error) {
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
	info.Pieces = pieces
	info.PieceLength = piecesLength
	info.Private = private
	info.Files = files
	return info, nil
}

func getAnnounceList(announceListData []interface{}) [][]string {
	announceList := [][]string{}

	for _, announceMinListData := range announceListData {
		announceMinList := []string{}

		announceNanoListData, _ := announceMinListData.([]interface{})
		for _, announceData := range announceNanoListData {
			announce, announceOK := announceData.(string)
			if announceOK {
				announceMinList = append(announceMinList, announce)
			}
		}

		if len(announceMinList) > 0 {
			announceList = append(announceList, announceMinList)
		}
	}

	return announceList
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
