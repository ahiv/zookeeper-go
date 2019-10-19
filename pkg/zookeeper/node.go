package zookeeper

type Node interface {
	ResolveChild(path string) Node
	ListChildren() []Node
	Watch() Watch
	Update(value interface{})
	Read() interface{}
	Delete()
	IsEphemeral() bool
}
