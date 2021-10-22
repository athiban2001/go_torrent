package utils

import "encoding/binary"

func IntToBytes(value int64, is64 bool) []byte {
	var bytes []byte
	if is64 {
		bytes = make([]byte, 8)
		binary.LittleEndian.PutUint64(bytes, uint64(value))
	} else {
		bytes = make([]byte, 4)
		binary.LittleEndian.PutUint32(bytes, uint32(value))
	}

	return bytes
}

func BytesToInt(bytes []byte, is64 bool) int64 {
	if is64 {
		value := binary.BigEndian.Uint64(bytes)
		return int64(value)
	} else {
		value := binary.BigEndian.Uint32(bytes)
		return int64(value)
	}
}
