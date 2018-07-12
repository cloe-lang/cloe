package ir

func intCodeToByteCode(is []int) []byte {
	bs := make([]byte, 0, len(is))

	for _, i := range is {
		bs = append(bs, byte(i))
	}

	return bs
}
