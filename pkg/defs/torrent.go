package defs

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/athiban2001/go_torrent/pkg/bencoding"
	"github.com/athiban2001/go_torrent/pkg/utils"
	"github.com/google/uuid"
)

type Torrent struct {
	ID           string     `json:"id,omitempty"`
	Name         string     `json:"name,omitempty"`
	Announce     string     `json:"announce,omitempty"`
	AnnounceList [][]string `json:"announce-list,omitempty"`
	CreationDate time.Time  `json:"creation date,omitempty"`
	Comment      string     `json:"comment,omitempty"`
	CreatedBy    string     `json:"created by,omitempty"`
	Encoding     string     `json:"encoding,omitempty"`
	Progress     string     `json:"progress,omitempty"`
	Info         *Info      `json:"info"`
}

func NewTorrent(filename string) (*Torrent, error) {
	metaDataBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to read file : %v", err.Error())
	}

	parser := &bencoding.ParserData{
		Data:   metaDataBytes,
		Length: len(metaDataBytes),
	}

	metaData, err := parser.ReadDictionary()
	if err != nil {
		return nil, err
	}
	absFileName, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}
	torrent, err := getTorrent(metaData, absFileName)
	if err != nil {
		return nil, err
	}
	if err = utils.CopyDataFile(torrent.Name, torrent.ID+".torrent"); err != nil {
		return nil, err
	}

	gtConfig, err := GetConfig()
	if err != nil {
		return nil, err
	}

	gtConfig.Torrents = append(gtConfig.Torrents, torrent)
	if err = WriteConfig(gtConfig); err != nil {
		return nil, err
	}

	return torrent, nil
}

func getTorrent(metaData map[string]interface{}, filename string) (*Torrent, error) {
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
	info, err := GetInfo(infoData)
	if err != nil {
		return nil, err
	}

	torrent := &Torrent{}
	torrent.ID = uuid.NewString()
	torrent.Name = filename
	torrent.CreationDate = time.Unix(epoch, 0)
	torrent.Comment = comment
	torrent.CreatedBy = createdBy
	torrent.Announce = announce
	torrent.AnnounceList = getAnnounceList(announceList)
	torrent.Encoding = encoding
	torrent.Info = info
	return torrent, nil
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
