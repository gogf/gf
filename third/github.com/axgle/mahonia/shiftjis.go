package mahonia

// Converters for the Shift-JIS encoding.

import (
	"sync"
)

func init() {
	RegisterCharset(&Charset{
		Name:    "Shift_JIS",
		Aliases: []string{"MS_Kanji", "csShiftJIS", "SJIS"},
		NewDecoder: func() Decoder {
			sjisOnce.Do(makeSjisTable)
			return sjisTable.Decoder()
		},
		NewEncoder: func() Encoder {
			sjisOnce.Do(makeSjisTable)
			return sjisTable.Encoder()
		},
	})
}

var sjisOnce sync.Once

var sjisTable MBCSTable

func makeSjisTable() {
	var b [2]byte

	for jis0208, unicode := range jis0208ToUnicode {
		if unicode == 0 {
			continue
		}

		j1 := byte(jis0208 >> 8)
		j2 := byte(jis0208)

		if j1 < 95 {
			b[0] = (j1+1)/2 + 112
		} else {
			b[0] = (j1+1)/2 + 176
		}

		if j1&1 == 1 {
			b[1] = j2 + 31
			if j2 >= 96 {
				b[1]++
			}
		} else {
			b[1] = j2 + 126
		}

		sjisTable.AddCharacter(rune(unicode), string(b[:]))
	}

	for jis0201, unicode := range jis0201ToUnicode {
		if unicode == 0 {
			continue
		}

		sjisTable.AddCharacter(rune(unicode), string(byte(jis0201)))
	}

	for i := '\x00'; i < 32; i++ {
		sjisTable.AddCharacter(i, string(byte(i)))
	}

	sjisTable.AddCharacter(0x7f, "\x7f")
}
