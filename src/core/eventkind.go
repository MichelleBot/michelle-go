package core

type EventKind string

const (
	EventUnknown    EventKind = "unknown"
	EventMessage    EventKind = "message"
	EventReaction   EventKind = "reaction"
	EventEdit       EventKind = "edit"
	EventRevoke     EventKind = "revoke"
	EventPoll       EventKind = "poll"
	EventReceipt    EventKind = "receipt"
	EventPresence   EventKind = "presence"
	EventGroup      EventKind = "group"
	EventCall       EventKind = "call"
	EventNewsletter EventKind = "newsletter"
	EventSystem     EventKind = "system"
)
