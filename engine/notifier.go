package engine

// A Notifier can be anything that is responsible for sending a notification to
// a user, e.g. email, push notification, text message, or carrier pidgeon.
type Notifier interface {
  Notify(who, rule_name, msg string) error
}
