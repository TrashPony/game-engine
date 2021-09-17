package game_math

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func IntFromBytes(b []byte) int32 {
	var x int32
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.BigEndian, &x)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}

	return x
}

func GetInt64Bytes(x int) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, int32(x))
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}

	return buf.Bytes()
}

func BoolToByte(bitSet bool) byte {
	if bitSet {
		return 1
	}

	return 0
}
