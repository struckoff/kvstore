package kvstore

type Balancer interface {
	AddNode(n Node) error
	RemoveNode(id string) error
	LocateKey(key string) (Node, error)
	Nodes() ([]Node, error)
}
