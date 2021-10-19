package gt

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/athiban2001/go_torrent/pkg/utils"
)

type GoTorrentConfig struct {
	PeerID   string     `json:"peerID"`
	Torrents []*Torrent `json:"torrents"`
}

func GetConfig() (*GoTorrentConfig, error) {
	configFileData, err := utils.GetConfigFile()
	if err != nil {
		return nil, err
	}

	gtConfig := &GoTorrentConfig{}
	if err = json.Unmarshal([]byte(configFileData), gtConfig); err != nil {
		fmt.Println(configFileData)
		return nil, err
	}

	return gtConfig, nil
}

func WriteConfig(config *GoTorrentConfig) error {
	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}

	configDirPath, err := utils.GetConfigDir()
	if err != nil {
		return err
	}

	configFilePath := path.Join(configDirPath, "config.json")
	file, err := os.Create(configFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = file.Write(data); err != nil {
		return err
	}

	return nil
}
