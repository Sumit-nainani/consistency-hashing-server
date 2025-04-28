package hashring

import (
	"fmt"
	"hash/crc32"
	"sort"
	"sync"
)

const (
	RING_SIZE = 64
)

var (
	once     sync.Once
	instance *HashRing
)

type HashRing struct {
	sync.RWMutex
	Nodes     []int
	IpHash    map[string]int
	KeyToIP   map[int]string
	KeyToNode map[string]string
	NodeToIP  map[string]string
}

func GetRingInstance() *HashRing {
	once.Do(func() {
		instance = &HashRing{
			IpHash:    make(map[string]int),
			KeyToIP:   make(map[int]string),
			KeyToNode: make(map[string]string),
			NodeToIP:  make(map[string]string),
		}
	})
	return instance
}

func (hr *HashRing) GetLenNodes() int {
	hr.RLock()
	defer hr.RUnlock()
	return len(hr.Nodes)
}

func getHashValue(key uint32) int {
	return int(key % RING_SIZE)
}

// Finding the best server for serving all requests whose key values is less than or qual to server's key value
// Using Binary search algorithm
func getNodeIp(ringNode []int, keyToIP map[int]string, hashValue int) string {
	l := 0
	h := len(ringNode) - 1

	for l < h {
		m := (l + h) / 2
		if ringNode[m] >= hashValue {
			h = m
		} else {
			l = m + 1
		}
	}
	if ringNode[l] >= hashValue {
		return keyToIP[ringNode[l]]
	} else if ringNode[h] >= hashValue {
		return keyToIP[ringNode[h]]
	} else {
		return keyToIP[ringNode[0]]
	}
}

func (hr *HashRing) AddNode(node string, ip string) int {
	hr.Lock()
	defer hr.Unlock()

	key := crc32.ChecksumIEEE([]byte(ip))
	hashValue := getHashValue(key)

	if _, ok := hr.IpHash[ip]; ok {
		return hashValue
	}

	hr.IpHash[ip] = 1
	hr.NodeToIP[node] = ip
	hr.KeyToIP[hashValue] = ip

	insertPos := sort.Search(len(hr.Nodes), func(i int) bool {
		return hr.Nodes[i] >= hashValue
	})

	// Insert the element at the found position
	hr.Nodes = append(hr.Nodes[:insertPos], append([]int{hashValue}, hr.Nodes[insertPos:]...)...)
	fmt.Println(hr.Nodes, hr.IpHash, hr.KeyToIP, hr.NodeToIP)
	return hashValue
}

func (hr *HashRing) RemoveNode(node string) int {
	hr.Lock()
	defer hr.Unlock()
	key := crc32.ChecksumIEEE([]byte(hr.NodeToIP[node]))
	hashValue := getHashValue(key)

	if _, ok := hr.NodeToIP[node]; !ok {
		return hashValue
	}

	delete(hr.KeyToIP, hashValue)
	delete(hr.IpHash, hr.NodeToIP[node])
	delete(hr.NodeToIP, node)

	for i := len(hr.Nodes) - 1; i >= 0; i-- {
		if hr.Nodes[i] == hashValue {
			hr.Nodes = append(hr.Nodes[:i], hr.Nodes[i+1:]...)
		}
	}
	fmt.Println("deleting the pod")
	fmt.Println(hr.Nodes, hr.IpHash, hr.KeyToIP, hr.NodeToIP)
	return hashValue
}

// GetNode returns the node responsible for the given key
func (hr *HashRing) GetNode(key string) (string, int){
	hr.RLock()
	defer hr.RUnlock()

	// if node, ok := hr.KeyToNode[key]; ok {
	// 	fmt.Println("already exist , ", node)
	// 	return node
	// }

	hash := crc32.ChecksumIEEE([]byte(key))
	hashValue := getHashValue(hash)

	return getNodeIp(hr.Nodes, hr.KeyToIP, hashValue),hashValue
}
