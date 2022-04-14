package utils

import (
	"bytes"
	"encoding/binary"
)

func Get(value, shift uint64) (continuous, singleCount uint64) {
	meta := value >> shift & 0b1111
	continuous = meta / 5
	singleCount = meta % 5
	return
}

func Set(value, shift uint64, distance, cnt uint64) uint64 {
	return (value & ^(0b1111 << shift)) | ((distance*5 + cnt) << shift)
}

func ToBytes(value uint64) []byte {
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, value)
	return buffer.Bytes()
}

func ToUInt64(buffer []byte) uint64 {
	var temp uint64
	bytesBuffer := bytes.NewBuffer(buffer)
	binary.Read(bytesBuffer, binary.BigEndian, &temp)
	return temp
}
