package gt

import (
	"bytes"
	"encoding/gob"
	"io"

	"github.com/athiban2001/go_torrent/pkg/utils"
)

func GetConfig() (*GoTorrentConfig, error) {
	configFileData, err := utils.GetConfigFile()
	if err != nil {
		return nil, err
	}

	if len(configFileData) == 0 {
		gtConfig := &GoTorrentConfig{}
		gtConfig.PeerID = GetPeerID()
		if err := WriteConfig(gtConfig); err != nil {
			return nil, err
		}

		return gtConfig, nil
	}

	gtConfig := &GoTorrentConfig{}
	buffer := bytes.NewBuffer(configFileData)
	if err := gob.NewDecoder(buffer).Decode(gtConfig); err != nil && err != io.EOF {
		return nil, err
	}

	return gtConfig, nil
}

func WriteConfig(config *GoTorrentConfig) error {
	data := make([]byte, 0)
	buffer := bytes.NewBuffer(data)
	if err := gob.NewEncoder(buffer).Encode(config); err != nil {
		return err
	}

	return utils.WriteConfigFile(buffer.Bytes())
}
