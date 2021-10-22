package utils

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"time"
)

var configFileData []byte = make([]byte, 0)
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

func GetConfigFile() ([]byte, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return nil, err
	}

	configFilePath := path.Join(configDir, "config.gt")

	file, err := os.OpenFile(configFilePath, os.O_CREATE|os.O_RDWR, 0664)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := os.Stat(configFilePath)
	if err != nil {
		return nil, err
	}

	if !stat.ModTime().Equal(configFileLastModified) {
		data, err := os.ReadFile(configFilePath)
		if err != nil {
			return nil, err
		}
		configFileData = data
		configFileLastModified = stat.ModTime()
	}

	return configFileData, nil
}

func WriteConfigFile(data []byte) error {
	configDirPath, err := GetConfigDir()
	if err != nil {
		return err
	}

	filePath := path.Join(configDirPath, "config.gt")
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0664)
	if err != nil {
		return err
	}

	if _, err = file.Write(data); err != nil {
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	configFileData = data
	configFileLastModified = stat.ModTime()

	return nil
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

func RemoveDataFile(filename string) error {
	configDirPath, err := GetConfigDir()
	if err != nil {
		return err
	}

	filePath := path.Join(configDirPath, "data", filename)
	return os.Remove(filePath)
}
