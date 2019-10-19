package zookeeper

type WatchEventKind int

const (
	WatchCloseEvent WatchEventKind = iota
	NodeUpdateEvent
	NodeDeleteEvent
)

type WatchEvent struct {
	Kind WatchEventKind
}

type WatchChannel <-chan WatchEvent

type Watch interface {
	Close()
	Channel() WatchChannel
}