package services

import (
	"log"
	"sync"
)

// NotificationKind classifies a notification.
type NotificationKind string

const (
	KindInfo    NotificationKind = "info"
	KindSuccess NotificationKind = "success"
	KindWarning NotificationKind = "warning"
	KindError   NotificationKind = "error"
)

// Notification is a single app-level message.
type Notification struct {
	Kind    NotificationKind
	Title   string
	Message string
}

// NotifyHandler is a callback invoked when a notification arrives.
type NotifyHandler func(n Notification)

// Notifier is a simple in-process pub/sub for UI notifications.
type Notifier struct {
	mu       sync.RWMutex
	handlers []NotifyHandler
}

// NewNotifier constructs an empty Notifier.
func NewNotifier() *Notifier {
	return &Notifier{}
}

// Subscribe registers a handler to receive notifications.
func (n *Notifier) Subscribe(h NotifyHandler) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.handlers = append(n.handlers, h)
}

// Notify broadcasts a notification to all subscribers.
func (n *Notifier) Notify(kind NotificationKind, title, message string) {
	note := Notification{Kind: kind, Title: title, Message: message}
	n.mu.RLock()
	handlers := make([]NotifyHandler, len(n.handlers))
	copy(handlers, n.handlers)
	n.mu.RUnlock()

	for _, h := range handlers {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("[notifier] handler panic: %v", r)
				}
			}()
			h(note)
		}()
	}
}

// Info is a shorthand for KindInfo notifications.
func (n *Notifier) Info(title, msg string)    { n.Notify(KindInfo, title, msg) }
func (n *Notifier) Success(title, msg string) { n.Notify(KindSuccess, title, msg) }
func (n *Notifier) Warning(title, msg string) { n.Notify(KindWarning, title, msg) }
func (n *Notifier) Error(title, msg string)   { n.Notify(KindError, title, msg) }