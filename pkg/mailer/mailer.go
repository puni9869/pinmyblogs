// Package mailer provides email sending services for user notifications.
package mailer

// MailerAPI defines the interface for sending emails.
type MailerAPI interface {
	Send()
}
