package gui

// clearScreen clears the terminal UI
func clearScreen() {
	print("\033[H\033[2J")
}
