package notif

import notifier "github.com/deckarep/gosx-notifier"

const notificationGroup = "com.matcornic.subify"

// Error send a notification error
func Error(title, message string) error {
	iconPath := downloadIcon()
	notification := notifier.Notification{
		Group:   notificationGroup,
		Title:   title,
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
		Title:   title,
		Message: message,
		Sound:   notifier.Pop,
		AppIcon: iconPath,
	}
	return notification.Push()
}
