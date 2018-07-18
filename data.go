package main

func appendShort(l []byte, r uint16) []byte {
	return append(l, byte(r>>8), byte(r))
}

func extractShort(d []byte) ([]byte, uint16, bool) {
	if len(d) < 2 {
		return nil, 0, false
	}

	return d[2:], uint16(d[0])<<8 |
		uint16(d[1]), true
}
