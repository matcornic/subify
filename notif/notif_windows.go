package notif

import (
	toast "github.com/jacobmarshall/go-toast"
)

// Error send a notification error
func Error(title, message string) error {
	iconPath := downloadIcon()
	notification := toast.Notification{
		AppID:   "com.matcornic.subify",
		Title:   title,
		Message: message,
	}
	if iconPath != "" {
		notification.Icon = iconPath
	}
	return notification.Push()
}

// Info send a notification information
func Info(title, message string) error {
	iconPath := downloadIcon()
	notification := toast.Notification{
		AppID:   "com.matcornic.subify",
		Title:   title,
		Message: message,
	}
	if iconPath != "" {
		notification.Icon = iconPath
	}
	return notification.Push()
}
