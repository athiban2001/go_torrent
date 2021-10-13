package main

import (
	"fmt"

	"github.com/athiban2001/go_torrent/pkg/torrent"
)

func main() {
	torrentFileName := "./multi_course.torrent"
	data, err := torrent.New(torrentFileName)
	if err != nil {
		fmt.Printf("unable to parse file : %v", err.Error())
	}
	fmt.Println("data : ", data)
}
