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

type RequestMetaData struct {
	RequestHash      int32
	AssignedNodeHash int32
	AssignedNodeIP   string
	AssignedNodeName string
}

type NodeMetaData struct {
	NodeHash int32
	NodeIP   string
	NodeName string
}

type HashRing struct {
	sync.RWMutex
	NodeMetaData           []NodeMetaData
	NodeNameToNodeMetaData map[string]NodeMetaData
	RequestIpToMetaData    map[string]RequestMetaData
}

func GetRingInstance() *HashRing {
	once.Do(func() {
		instance = &HashRing{
			NodeNameToNodeMetaData: make(map[string]NodeMetaData),
			RequestIpToMetaData:    make(map[string]RequestMetaData),
		}
	})
	return instance
}

func (hr *HashRing) getNodeIp(RequestIPHashValue int32) NodeMetaData {
	l := 0
	h := len(hr.NodeMetaData) - 1

	for l < h {
		m := (l + h) / 2
		if hr.NodeMetaData[m].NodeHash >= RequestIPHashValue {
			h = m
		} else {
			l = m + 1
		}
	}
	if hr.NodeMetaData[l].NodeHash >= RequestIPHashValue {
		return hr.NodeMetaData[l]
	} else if hr.NodeMetaData[h].NodeHash >= RequestIPHashValue {
		return hr.NodeMetaData[h]
	} else {
		return hr.NodeMetaData[0]
	}
}

func (hr *HashRing) AddNode(nodeName string, nodeIp string) {
	hr.Lock()
	defer hr.Unlock()

	nodeHashValue := int32(crc32.ChecksumIEEE([]byte(nodeIp)) % RING_SIZE)

	insertPos := sort.Search(len(hr.NodeMetaData), func(i int) bool {
		return hr.NodeMetaData[i].NodeHash >= nodeHashValue
	})

	hr.NodeNameToNodeMetaData[nodeName] = NodeMetaData{NodeHash: nodeHashValue, NodeName: nodeName, NodeIP: nodeIp}
	hr.NodeMetaData = append(hr.NodeMetaData[:insertPos], append([]NodeMetaData{{
		NodeHash: nodeHashValue,
		NodeIP:   nodeIp,
		NodeName: nodeName,
	}}, hr.NodeMetaData[insertPos:]...)...)
	fmt.Println(hr.NodeMetaData, hr.NodeNameToNodeMetaData)
}

func (hr *HashRing) RemoveNode(nodeName string) {
	hr.Lock()
	defer hr.Unlock()

	for i := len(hr.NodeMetaData) - 1; i >= 0; i-- {
		if hr.NodeMetaData[i].NodeName == nodeName {
			name := hr.NodeMetaData[i].NodeName
			hr.NodeMetaData = append(hr.NodeMetaData[:i], hr.NodeMetaData[i+1:]...)
			delete(hr.NodeNameToNodeMetaData, name)
		}
	}
}

func (hr *HashRing) GetNode(requestIP string) {
	hr.RLock()
	defer hr.RUnlock()

	RequestIPHashValue := int32(crc32.ChecksumIEEE([]byte(requestIP)) % RING_SIZE)
	AssignedNodeMetaData := hr.getNodeIp(RequestIPHashValue)

	hr.RequestIpToMetaData[requestIP] = RequestMetaData{RequestHash: RequestIPHashValue, AssignedNodeHash: AssignedNodeMetaData.NodeHash, AssignedNodeName: AssignedNodeMetaData.NodeName, AssignedNodeIP: AssignedNodeMetaData.NodeIP}
}
