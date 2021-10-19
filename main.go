package main

import (
	"fmt"

	"github.com/athiban2001/go_torrent/pkg/gt"
)

func main() {
	torrentFileName := "./multi_course.torrent"
	_, err := gt.NewTorrent(torrentFileName)
	if err != nil {
		fmt.Printf("unable to parse file : %v", err.Error())
	}
}
