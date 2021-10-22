package main

import (
	"fmt"
	"log"

	"github.com/athiban2001/go_torrent/pkg/gt"
)

func main() {
	// torrentFileName := "./multi_course.torrent"
	// _, err := torrent.New(torrentFileName)
	// if err != nil {
	// 	log.Fatalf(err.Error())
	// }

	gtConfig, err := gt.GetConfig()
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println(gtConfig.Torrents[0].Name)
}
