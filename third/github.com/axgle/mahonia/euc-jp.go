package mahonia

// Converters for the EUC-JP encoding

import (
	"sync"
)

func init() {
	RegisterCharset(&Charset{
		Name:    "EUC-JP",
		Aliases: []string{"extended_unix_code_packed_format_for_japanese", "cseucpkdfmtjapanese"},
		NewDecoder: func() Decoder {
			eucJPOnce.Do(makeEUCJPTable)
			return eucJPTable.Decoder()
		},
		NewEncoder: func() Encoder {
			eucJPOnce.Do(makeEUCJPTable)
			return eucJPTable.Encoder()
		},
	})
}

var eucJPOnce sync.Once
var eucJPTable MBCSTable

func makeEUCJPTable() {
	var b [3]byte

	b[0] = 0x8f
	for jis0212, unicode := range jis0212ToUnicode {
		if unicode == 0 {
			continue
		}

		b[1] = byte(jis0212>>8) | 128
		b[2] = byte(jis0212) | 128
		eucJPTable.AddCharacter(rune(unicode), string(b[:3]))
	}

	for jis0208, unicode := range jis0208ToUnicode {
		if unicode == 0 {
			continue
		}

		b[0] = byte(jis0208>>8) | 128
		b[1] = byte(jis0208) | 128
		eucJPTable.AddCharacter(rune(unicode), string(b[:2]))
	}

	b[0] = 0x8e
	for i := 128; i < 256; i++ {
		unicode := jis0201ToUnicode[i]
		if unicode == 0 {
			continue
		}

		b[1] = byte(i)
		eucJPTable.AddCharacter(rune(unicode), string(b[:2]))
	}

	for i := '\x00'; i < 128; i++ {
		var unicode rune
		if i < 32 || i == 127 {
			unicode = i
		} else {
			unicode = rune(jis0201ToUnicode[i])
			if unicode == 0 {
				continue
			}
		}

		eucJPTable.AddCharacter(unicode, string(byte(i)))
	}
}
