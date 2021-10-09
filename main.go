package main

import (
	"fmt"
	"os"

	"github.com/athiban2001/go_torrent/pkg/torrent"
)

type Torrent struct {
}

func NewTorrent(filename string) (*Torrent, error) {
	torrentData, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to read file : %v", err.Error())
	}

	data, err := torrent.InitTorrent(torrentData)
	if err != nil {
		return nil, fmt.Errorf("unable to parse file : %v", err.Error())
	}
	fmt.Println(data["info"])

	return nil, nil
}

func main() {
	torrentFileName := "ubuntu-21.04-desktop-amd64.iso.torrent"
	NewTorrent(torrentFileName)
}
