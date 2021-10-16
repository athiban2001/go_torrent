package main

import (
	"fmt"

	"github.com/athiban2001/go_torrent/pkg/types"
)

func main() {
	torrentFileName := "./multi_course.torrent"
	data, err := types.NewTorrent(torrentFileName)
	if err != nil {
		fmt.Printf("unable to parse file : %v", err.Error())
	}
	fmt.Println(data.Info.GetInfoHash())
}
