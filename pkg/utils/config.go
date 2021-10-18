package utils

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math/rand"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

var configFileData string = "{}"
var configFileLastModified time.Time = time.Now()

func GetConfigDir() (string, error) {
	configPath := ""

	if runtime.GOOS == "windows" && os.Getenv("APP_DATA") != "" {
		configPath = path.Join(os.Getenv("APP_DATA"), "go_torrent")
	} else if (runtime.GOOS == "darwin" || runtime.GOOS == "linux") && os.Getenv("HOME") != "" {
		configPath = path.Join(os.Getenv("HOME"), ".config", "go_torrent")
	} else {
		return "", fmt.Errorf("appropriate environment variables not set")
	}

	if err := os.MkdirAll(path.Join(configPath, "data"), os.ModeDir|os.ModePerm); err != nil {
		return "", err
	}

	return configPath, nil
}

func GetConfigFile() (string, error) {
	var stat os.FileInfo
	var err error

	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}

	configFilePath := path.Join(configDir, "config.json")
	stat, err = os.Stat(configFilePath)
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		peerID := getPeerID()
		jsonInitData := "{\"peerID\":\"" + peerID + "\"}"

		file, err := os.Create(configFilePath)
		if err != nil {
			return "", err
		}
		defer file.Close()

		if _, err = file.WriteString(jsonInitData); err != nil {
			return "", err
		}
		if stat, err = file.Stat(); err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}

	if !stat.ModTime().Equal(configFileLastModified) {
		data, err := os.ReadFile(configFilePath)
		if err != nil {
			return "", err
		}
		configFileData = string(data)
		configFileLastModified = stat.ModTime()
	}

	return configFileData, nil
}

func CopyDataFile(src, destFileName string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	configDirPath, err := GetConfigDir()
	if err != nil {
		return err
	}

	out, err := os.Create(path.Join(configDirPath, "data", destFileName))
	if err != nil {
		return err
	}

	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	if err = out.Close(); err != nil {
		return err
	}

	return nil
}

func getPeerID() string {
	prefix := "-GT0001-"
	possibleDigits := []rune("0123456789")

	rand.Seed(time.Now().UnixNano())
	b := strings.Builder{}

	for i := 1; i <= 12; i++ {
		b.WriteRune(possibleDigits[rand.Intn(len(possibleDigits))])
	}

	return prefix + b.String()
}
