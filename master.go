package service_discovery

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

type Master struct {
	sync.RWMutex
	kapi   client.KeysAPI
	key    string
	nodes  map[string]string
	active bool
}

func NewMaster(serviceName string, endpoints []string) (*Master, error) {
	cfg := client.Config{
		Endpoints:               endpoints,
		HeaderTimeoutPerRequest: time.Second * 2,
	}
	c, err := client.New(cfg)
	if err != nil {
		return nil, err
	}
	master := &Master{
		kapi:   client.NewKeysAPI(c),
		key:    fmt.Sprintf("%s/%s/", kRoot, serviceName),
		nodes:  make(map[string]string),
		active: true,
	}
	master.fetch()

	/// `fetch` Timer may work well too?
	go master.watch()

	return master, err
}

func (m *Master) GetNodesStrictly() map[string]string {
	//log.Println("strictly active ->", m.active)
	if !m.active {
		return nil
	}
	return m.GetNodes()
}

func (m *Master) GetNodes() map[string]string {
	m.RLock()
	defer m.RUnlock()
	return m.nodes
}

func (m *Master) addNode(node, extInfo string) {
	m.Lock()
	defer m.Unlock()
	node = strings.TrimLeft(node, m.key)
	m.nodes[node] = extInfo
}

func (m *Master) delNode(node string) {
	m.Lock()
	defer m.Unlock()
	node = strings.TrimLeft(node, m.key)
	delete(m.nodes, node)
}

func (m *Master) watch() {
	watcher := m.kapi.Watcher(m.key, &client.WatcherOptions{
		Recursive: true,
	})
	for {
		resp, err := watcher.Next(context.Background())
		if err != nil {
			log.Println(err)
			m.active = false
			continue
		}
		m.active = true
		//log.Println("loop active ->", m.active)
		switch resp.Action {
		case "set", "update":
			m.addNode(resp.Node.Key, resp.Node.Value)
			break
		case "expire", "delete":
			m.delNode(resp.Node.Key)
			break
		default:
			log.Println("watchme!!!", "resp ->", resp)
		}
	}
}

func (m *Master) fetch() error {
	resp, err := m.kapi.Get(context.Background(), m.key, nil)
	if err != nil {
		return err
	}
	if resp.Node.Dir {
		for _, v := range resp.Node.Nodes {
			m.addNode(v.Key, v.Value)
		}
	}
	return err
}
