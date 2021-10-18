package main

import (
	"fmt"

	"github.com/athiban2001/go_torrent/pkg/defs"
)

func main() {
	torrentFileName := "./multi_course.torrent"
	_, err := defs.NewTorrent(torrentFileName)
	if err != nil {
		fmt.Printf("unable to parse file : %v", err.Error())
	}
}
