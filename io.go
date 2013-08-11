package nbt

import (
	"encoding/binary"
	"io"
)

var (
	byteOrder = binary.BigEndian
)

func read(rd io.Reader, value interface{}) {
	binary.Read(rd, byteOrder, value)
}

func readByte(rd io.Reader) (result byte) {
	read(rd, &result)
	return result
}

func readShort(rd io.Reader) (result int16) {
	read(rd, &result)
	return result
}

func readInt(rd io.Reader) (result int32) {
	read(rd, &result)
	return result
}

func readLong(rd io.Reader) (result int64) {
	read(rd, &result)
	return result
}

func readFloat(rd io.Reader) (result float32) {
	read(rd, &result)
	return result
}

func readDouble(rd io.Reader) (result float64) {
	read(rd, &result)
	return result
}

func write(w io.Writer, value interface{}) error {
	return binary.Write(w, byteOrder, value)
}

func readByteArray(rd io.Reader) (result []byte) {
	if length := readInt(rd); length > 0 {
		bytes := make([]byte, length)
		io.ReadFull(rd, bytes)
		return bytes
	} else {
		return []byte{}
	}
}

func writeByteArray(w io.Writer, value []byte) error {
	write(w, int32(len(value)))
	return write(w, value)
}

func readString(rd io.Reader) string {
	if length := readShort(rd); length > 0 {
		bytes := make([]byte, length)
		binary.Read(rd, byteOrder, bytes)
		return string(bytes)
	} else {
		return ""
	}
}

func writeString(w io.Writer, value string) error {
	write(w, int16(len(value)))
	return write(w, []byte(value))
}

func readIntArray(rd io.Reader) (result []int32) {
	if length := readInt(rd); length > 0 {
		ints := make([]int32, length)
		binary.Read(rd, byteOrder, ints)
		return ints
	} else {
		return []int32{}
	}
}

func writeIntArray(w io.Writer, value []int32) error {
	write(w, int32(len(value)))
	return write(w, value)
}
