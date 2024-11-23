package util

const ErrKey = "error_message"

func HexMap() [16]byte {
	return [16]byte{
		'0', '1', '2', '3',
		'4', '5', '6', '7',
		'8', '9', 'a', 'b',
		'c', 'd', 'e', 'f',
	}
}
