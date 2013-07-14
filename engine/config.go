package engine

type Config interface {
  // TODO(jmeady): Remove this. Notification mechanism config such as this
  //               should be handled more generically.
  SmtpConfig() (user, pass, host string)
  DefaultNotifier() (Notifier, error)
  DefaultNotificationRecipient() string
  LogLevel() string
  Rules() []Rule
}
