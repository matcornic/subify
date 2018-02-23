package notif

import (
	"os/exec"
)

func sendMessage(title, message string) error {
	subifyIcon := downloadIcon()
	exec.Command("notify-send", "-i", subifyIcon, title, message).Run()
}

// Error send a notification error
func Error(title, message string) error {
	return sendMessage(title, message)
}

// Info send a notification information
func Info(title, message string) error {
	return sendMessage(title, message)
}
