package torrent

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/athiban2001/go_torrent/pkg/bencoding"
	"github.com/athiban2001/go_torrent/pkg/gt"
	"github.com/athiban2001/go_torrent/pkg/utils"
	"github.com/google/uuid"
)

func New(filename string) (*gt.Torrent, error) {
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

	gtConfig, err := gt.GetConfig()
	if err != nil {
		return nil, err
	}

	gtConfig.Torrents = append(gtConfig.Torrents, torrent)
	if err = gt.WriteConfig(gtConfig); err != nil {
		return nil, err
	}

	return torrent, nil
}

func Remove(torrent *gt.Torrent) error {
	gtConfig, err := gt.GetConfig()
	if err != nil {
		return err
	}
	newTorrentList := make([]*gt.Torrent, 0)

	for _, currentTorrent := range gtConfig.Torrents {
		if currentTorrent.ID != torrent.ID {
			newTorrentList = append(newTorrentList, currentTorrent)
		}
	}

	gtConfig.Torrents = newTorrentList
	return gt.WriteConfig(gtConfig)
}

func getTorrent(metaData map[string]interface{}, filename string) (*gt.Torrent, error) {
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

	torrent := &gt.Torrent{}
	torrent.ID = uuid.NewString()
	torrent.Name = filename
	torrent.CreationDate = time.Unix(epoch, 0)
	torrent.Comment = comment
	torrent.CreatedBy = createdBy
	torrent.Announce = announce
	torrent.AnnounceList = getAnnounceList(announceList)
	torrent.Encoding = encoding
	torrent.Downloaded = 0
	torrent.Uploaded = 0
	torrent.Info = info
	torrent.Left = getSize(torrent.Info.Files)
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
