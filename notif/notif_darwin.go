package notif

import (
	"fmt"

	notifier "github.com/deckarep/gosx-notifier"
)

const notificationGroup = "com.matcornic.subify"

// SendSubtitleDownloadSuccess sends a notification when download went well
func SendSubtitleDownloadSuccess(successAPI string) {
	Info("I found a subtitle for your video 😎", fmt.Sprintf("Thank you %s ❤️", successAPI))
}

// SendSubtitleCouldNotBeDownloaded sends a notification when download went bad
func SendSubtitleCouldNotBeDownloaded(noSucessAPIs string) {
	Error("‼️ I didn't found any subtitle 😭", fmt.Sprintf("No match for your video in : %s. Try later !", noSucessAPIs))
}

// Error send a notification error
func Error(title, message string) error {
	iconPath := downloadIcon()
	notification := notifier.Notification{
		Group:   notificationGroup,
		Title:   fmt.Sprintf("Subify - %s", title),
		Message: message,
		Sound:   notifier.Basso,
		AppIcon: iconPath,
	}
	return notification.Push()
}

// Info send a notification information
func Info(title, message string) error {
	iconPath := downloadIcon()
	notification := notifier.Notification{
		Group:   notificationGroup,
		Title:   fmt.Sprintf("Subify - %s", title),
		Message: message,
		Sound:   notifier.Pop,
		AppIcon: iconPath,
	}
	return notification.Push()
}
