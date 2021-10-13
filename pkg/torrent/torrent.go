package torrent

import (
	"fmt"
	"log"
	"os"

	"github.com/athiban2001/go_torrent/pkg/bencoding"
)

func New(filename string) (*Torrent, error) {
	metaDataBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to read file : %v", err.Error())
	}

	parser := &bencoding.ParserData{
		Data:   metaDataBytes,
		Index:  0,
		Length: len(metaDataBytes),
	}

	metaData, err := parser.ReadDictionary()
	if err != nil {
		return nil, err
	}

	torrent, err := getTorrent(metaData)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println(torrent.Announce, torrent.CreationDate, len(torrent.AnnounceList), torrent.Info.PieceLength, len(torrent.Info.Files), len(torrent.Info.Pieces), torrent.Info.Private)
	return &Torrent{}, nil
}
