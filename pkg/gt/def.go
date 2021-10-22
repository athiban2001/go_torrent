package gt

import (
	"time"
)

type GoTorrentConfig struct {
	PeerID   string     `json:"peerID"`
	Torrents []*Torrent `json:"torrents"`
}

type Torrent struct {
	ID           string     `json:"id,omitempty"`
	Name         string     `json:"name,omitempty"`
	Announce     string     `json:"announce,omitempty"`
	AnnounceList [][]string `json:"announce-list,omitempty"`
	CreationDate time.Time  `json:"creation date,omitempty"`
	Comment      string     `json:"comment,omitempty"`
	CreatedBy    string     `json:"created by,omitempty"`
	Encoding     string     `json:"encoding,omitempty"`
	Uploaded     int64      `json:"uploaded"`
	Downloaded   int64      `json:"downloaded"`
	Left         int64      `json:"left"`
	Progress     string     `json:"progress,omitempty"`
	Info         *Info      `json:"info"`
}

type Info struct {
	Name        string  `json:"name,omitempty"`
	PieceLength int64   `json:"piece-length,omitempty"`
	Pieces      string  `json:"pieces,omitempty"`
	Private     bool    `json:"private,omitempty"`
	Files       []*File `json:"files"`
}

type File struct {
	Length int64    `json:"length,omitempty"`
	MD5Sum string   `json:"md5sum,omitempty"`
	Path   []string `json:"path"`
}
