package nodes

import (
	"github.com/TrashPony/game-engine/router/mechanics/game_objects/node"
	"sync"
)

var nodes *store

type store struct {
	nodes map[string]*node.Node
	mx    sync.RWMutex
}

func Nodes() *store {

	if nodes == nil {
		nodes = newStore()
	}

	return nodes
}

func newStore() *store {
	g := &store{
		nodes: make(map[string]*node.Node),
	}

	return g
}

func (s *store) AddNode(name, url string, maxSessions int) (*node.Node, error) {

	newNode := &node.Node{Name: name, Url: url, MaxSessions: maxSessions}
	err := newNode.Connect()
	if err != nil {
		newNode.Stop()
		return nil, err
	}

	s.addNode(newNode)
	return newNode, nil
}

func (s *store) GetNode() *node.Node {
	s.mx.RLock()
	defer s.mx.RUnlock()

	for _, n := range s.nodes {
		return n
	}

	return nil
}

func (s *store) GetNodeByName(name string) *node.Node {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.nodes[name]
}

func (s *store) addNode(newNode *node.Node) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.nodes[newNode.Name] = newNode
}

func (s *store) RemoveNode(newNodeName string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	delete(s.nodes, newNodeName)
}

func (s *store) RangeNodes() <-chan *node.Node {
	s.mx.RLock()

	nodes := make(chan *node.Node, len(s.nodes))

	go func() {
		defer func() {
			s.mx.RUnlock()
			close(nodes)
		}()

		for _, n := range s.nodes {
			nodes <- n
		}
	}()

	return nodes
}
