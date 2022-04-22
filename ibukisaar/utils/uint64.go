package utils

import (
	"bytes"
	"encoding/binary"
)

type UInt64Slice []uint64

func (x UInt64Slice) Len() int           { return len(x) }
func (x UInt64Slice) Less(i, j int) bool { return x[i] < x[j] }
func (x UInt64Slice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

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
	binary.Write(&buffer, binary.LittleEndian, value)
	return buffer.Bytes()
}

func ToUInt64(buffer []byte) uint64 {
	var temp uint64
	bytesBuffer := bytes.NewBuffer(buffer)
	binary.Read(bytesBuffer, binary.LittleEndian, &temp)
	return temp
}
