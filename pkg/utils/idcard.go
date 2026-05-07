package utils

import "strings"

func MaskIDCard(id string) string {
	if len(id) <= 3 {
		return strings.Repeat("x", len(id))
	}

	firstChar := id[:1]
	lastTwoChars := id[len(id)-2:]
	
	middleX := strings.Repeat("x", len(id)-3)

	return firstChar + middleX + lastTwoChars
}
