package core_aes

func RtLength(keybits int) int {
	return (keybits)/8 + 28
}
