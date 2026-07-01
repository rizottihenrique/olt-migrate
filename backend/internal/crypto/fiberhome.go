package crypto

import (
	"strings"
)

// Mapas de criptografia e descriptografia da Fiberhome
var (
	decryptMap = map[rune]rune{
		'=': 'a', '<': 'b', ';': 'c', ':': 'd', '9': 'e',
		'8': 'f', '7': 'g', '6': 'h', '5': 'i', '4': 'j',
		'3': 'k', '2': 'l', '1': 'm', '0': 'n', '/': 'o',
		'.': 'p', '-': 'q', ',': 'r', '+': 's', '*': 't',
		')': 'u', '(': 'v', '\'': 'w', '&': 'x', '%': 'y',
		'$': 'z',
		'n': '0', 'm': '1', 'l': '2', 'k': '3', 'j': '4',
		'i': '5', 'h': '6', 'g': '7', 'f': '8', 'e': '9',
	}

	encryptMap = map[rune]rune{
		'a': '=', 'b': '<', 'c': ';', 'd': ':', 'e': '9',
		'f': '8', 'g': '7', 'h': '6', 'i': '5', 'j': '4',
		'k': '3', 'l': '2', 'm': '1', 'n': '0', 'o': '/',
		'p': '.', 'q': '-', 'r': ',', 's': '+', 't': '*',
		'u': ')', 'v': '(', 'w': '\'', 'x': '&', 'y': '%',
		'z': '$',
		'0': 'n', '1': 'm', '2': 'l', '3': 'k', '4': 'j',
		'5': 'i', '6': 'h', '7': 'g', '8': 'f', '9': 'e',
	}
)

// DecryptPassword descriptografa a senha obfuscada das OLTs Fiberhome
func DecryptPassword(encrypted string) string {
	var builder strings.Builder
	for _, ch := range encrypted {
		if plain, ok := decryptMap[ch]; ok {
			builder.WriteRune(plain)
		} else {
			builder.WriteRune(ch)
		}
	}
	return builder.String()
}

// EncryptPassword criptografa uma senha em texto plano para o formato obfuscado Fiberhome
func EncryptPassword(plaintext string) string {
	var builder strings.Builder
	for _, ch := range plaintext {
		if enc, ok := encryptMap[ch]; ok {
			builder.WriteRune(enc)
		} else {
			builder.WriteRune(ch)
		}
	}
	return builder.String()
}
