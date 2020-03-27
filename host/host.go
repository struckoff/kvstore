package host

import (
	"bytes"
	"encoding/json"
	"errors"
	balancer "github.com/struckoff/SFCFramework"
	"github.com/struckoff/kvstore/node"
	"log"
	"net/http"
	"strings"
)

// Host represents bounding of network api with balancer lib and local node
type Host struct {
	bal *balancer.Balancer
	n *node.InternalNode
}

func NewHost(n *node.InternalNode, bal *balancer.Balancer) (*Host, error) {
	if err := bal.AddNode(n); err != nil{
		return nil, err
	}
	return &Host{bal, n}, nil
}

// Lookup tries to connect to the remote node using addresses in the given list.
// Function ends on first success.
// If all attempts fail it will return an error.
func (h *Host) Lookup(eps []string) error{
	for _, addr := range eps {
		p := strings.Join([]string{addr, "node"}, "/")
		m := h.n.Meta()
		buf := bytes.NewBuffer(nil)
		if err := json.NewEncoder(buf).Encode(m); err != nil{
			log.Println("[ERROR]", err)
			continue
		}
		r, err := http.Post(p,"application", buf )
		if err != nil{
			log.Println("[ERROR]",err)
			continue
		}
		if r.StatusCode >= 400{
			log.Println("[ERROR]",r.Status)
			continue
		}
		var metas []node.NodeMeta
		if err := json.NewDecoder(r.Body).Decode(&metas); err != nil{
			log.Println("[ERROR]",err)
			continue
		}
		for _, meta := range metas {
			en := node.NewExternalNode(meta)
			if err := h.bal.AddNode(en); err != nil {
				return err
			}
		}
		return nil
	}
	return errors.New("unable to connect to nodes")
}

// AddNode adds node to balancer
func (h *Host) AddNode(n node.Node) error {
	return h.bal.AddNode(n)
}

// RemoveNode removes node from balancer
func (h *Host) RemoveNode(id string) error {
	return h.bal.RemoveNodeByID(id)
}

// Returns node from balancer by given key.
func (h *Host) GetNode(key string) (node.Node, error) {
	di := DataItem(key)
	nb, err := h.bal.LocateData(di)
	if err != nil {
		return nil, err
	}
	n, ok := nb.(node.Node)
	if !ok {
		return nil, errors.New("wrong node type")
	}
	return n, nil
}
