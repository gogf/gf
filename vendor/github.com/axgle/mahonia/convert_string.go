package mahonia

// ConvertString converts a  string from UTF-8 to e's encoding.
func (e Encoder) ConvertString(s string) string {
	dest := make([]byte, len(s)+10)
	destPos := 0

	for _, rune := range s {
	retry:
		size, status := e(dest[destPos:], rune)

		if status == NO_ROOM {
			newDest := make([]byte, len(dest)*2)
			copy(newDest, dest)
			dest = newDest
			goto retry
		}

		if status == STATE_ONLY {
			destPos += size
			goto retry
		}

		destPos += size
	}

	return string(dest[:destPos])
}

// ConvertString converts a string from d's encoding to UTF-8.
func (d Decoder) ConvertString(s string) string {
	bytes := []byte(s)
	runes := make([]rune, len(s))
	destPos := 0

	for len(bytes) > 0 {
		c, size, status := d(bytes)

		if status == STATE_ONLY {
			bytes = bytes[size:]
			continue
		}

		if status == NO_ROOM {
			c = 0xfffd
			size = len(bytes)
			status = INVALID_CHAR
		}

		bytes = bytes[size:]
		runes[destPos] = c
		destPos++
	}

	return string(runes[:destPos])
}
