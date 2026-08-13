package util

import "strings"

type lang struct {
	latin    string
	georgian string
}

var translitMap = []lang{
	{
		latin:    "ch'",
		georgian: "ჭ",
	},
	{
		latin:    "ts'",
		georgian: "წ",
	},
	{
		latin:    "ch",
		georgian: "ჩ",
	},
	{
		latin:    "kh",
		georgian: "ხ",
	},
	{
		latin:    "dz",
		georgian: "ძ",
	},
	{
		latin:    "ts",
		georgian: "ც",
	},
	{
		latin:    "p'",
		georgian: "პ",
	},
	{
		latin:    "sh",
		georgian: "შ",
	},
	{
		latin:    "q'",
		georgian: "ყ",
	},
	{
		latin:    "k'",
		georgian: "კ",
	},
	{
		latin:    "gh",
		georgian: "ღ",
	},
	{
		latin:    "t'",
		georgian: "ტ",
	},
	{
		latin:    "zh",
		georgian: "ჟ",
	},
	{
		latin:    "i",
		georgian: "ი",
	},
	{
		latin:    "l",
		georgian: "ლ",
	},
	{
		latin:    "n",
		georgian: "ნ",
	},
	{
		latin:    "r",
		georgian: "რ",
	},
	{
		latin:    "s",
		georgian: "ს",
	},
	{
		latin:    "m",
		georgian: "მ",
	},
	{
		latin:    "u",
		georgian: "უ",
	},
	{
		latin:    "p",
		georgian: "ფ",
	},
	{
		latin:    "k",
		georgian: "ქ",
	},
	{
		latin:    "o",
		georgian: "ო",
	},
	{
		latin:    "a",
		georgian: "ა",
	},
	{
		latin:    "t",
		georgian: "თ",
	},
	{
		latin:    "z",
		georgian: "ზ",
	},
	{
		latin:    "v",
		georgian: "ვ",
	},
	{
		latin:    "e",
		georgian: "ე",
	},
	{
		latin:    "d",
		georgian: "დ",
	},
	{
		latin:    "g",
		georgian: "გ",
	},
	{
		latin:    "b",
		georgian: "ბ",
	},
	{
		latin:    "j",
		georgian: "ჯ",
	},
	{
		latin:    "h",
		georgian: "ჰ",
	},
}

func ToGeorgian(input string) string {
	var builder strings.Builder
	input = strings.ToLower(input)

	i := 0
	for i < len(input) {
		matched := false
		for _, v := range translitMap {
			if strings.HasPrefix(input[i:], v.latin) {
				builder.WriteString(v.georgian)
				i += len(v.latin)
				matched = true
				break
			}
		}
		if !matched {
			builder.WriteByte(input[i])
			i++
		}
	}

	return builder.String()
}
