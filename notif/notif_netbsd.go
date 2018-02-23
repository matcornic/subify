package notif

// Notifications are disabled for netbsd but we want a netbsd executable though

// SendSubtitleDownloadSuccess sends a notification when download went well
func SendSubtitleDownloadSuccess(successAPI string) {
	return
}

// SendSubtitleCouldNotBeDownloaded sends a notification when download went bad
func SendSubtitleCouldNotBeDownloaded(noSucessAPIs string) {
	return
}

// Error send a notification error
func Error(title, message string) error {
	return nil
}

// Info send a notification information
func Info(title, message string) error {
	return nil
}
