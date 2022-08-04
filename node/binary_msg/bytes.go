package binary_msg

func GetIntBytes(x int) []byte {

	b := []byte{
		byte(0xff & (x >> 24)),
		byte(0xff & (x >> 16)),
		byte(0xff & (x >> 8)),
		byte(0xff & x),
	}

	return b
}

func ReuseByteSlice(s *[]byte, startIndex int, appendSlice []byte) {

	if len(*s) < startIndex {
		panic("small slice")
	}

	for i, b := range appendSlice {
		if len(*s) > startIndex+i {
			(*s)[startIndex+i] = b
		} else {
			*s = append(*s, b)
		}
	}
}

func BoolToByte(bitSet bool) byte {
	if bitSet {
		return 1
	}

	return 0
}
