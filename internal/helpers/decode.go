package helpers

// ISO8859ToUTF8 is a function to convert iso 8859-1 to utf-8
func ISO8859ToUTF8(iso8859 []byte) []byte {
	buffer := make([]rune, len(iso8859))
	for i, b := range iso8859 {
		buffer[i] = rune(b)
	}

	return []byte(string(buffer))
}
