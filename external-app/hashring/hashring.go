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
	RequestIpCollection    map[string][]string
}

func GetRingInstance() *HashRing {
	once.Do(func() {
		instance = &HashRing{
			NodeNameToNodeMetaData: make(map[string]NodeMetaData),
			RequestIpToMetaData:    make(map[string]RequestMetaData),
			RequestIpCollection:    make(map[string][]string),
		}
	})
	return instance
}

func (hr *HashRing) getNodeMetaData(RequestIPHashValue int32) NodeMetaData {

	if len(hr.NodeMetaData) == 0 {
		return NodeMetaData{}
	}

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

	if _, isAlreadyNodePresent := hr.NodeNameToNodeMetaData[nodeName]; isAlreadyNodePresent {
		return
	}

	nodeHashValue := int32(crc32.ChecksumIEEE([]byte(nodeIp))%RING_SIZE) + 1

	insertPos := sort.Search(len(hr.NodeMetaData), func(i int) bool {
		return hr.NodeMetaData[i].NodeHash >= nodeHashValue
	})

	hr.NodeNameToNodeMetaData[nodeName] = NodeMetaData{NodeHash: nodeHashValue, NodeName: nodeName, NodeIP: nodeIp}

	hr.NodeMetaData = append(hr.NodeMetaData[:insertPos], append([]NodeMetaData{{
		NodeHash: nodeHashValue,
		NodeIP:   nodeIp,
		NodeName: nodeName,
	}}, hr.NodeMetaData[insertPos:]...)...)
	fmt.Println(hr.NodeMetaData)

	for requestIP, requestIPMetaData := range hr.RequestIpToMetaData {
		nodeMetaData := hr.getNodeMetaData(requestIPMetaData.RequestHash)

		if nodeMetaData.NodeName == nodeName {
			oldNodeIp := requestIPMetaData.AssignedNodeIP
			requestIPCollection := hr.RequestIpCollection[oldNodeIp]

			for index, rerequestIPs := range requestIPCollection {
				if rerequestIPs == requestIP {
					hr.RequestIpCollection[oldNodeIp] = append(requestIPCollection[:index], requestIPCollection[index+1:]...)
					break
				}
			}

			requestIPMetaData.AssignedNodeHash = nodeMetaData.NodeHash
			requestIPMetaData.AssignedNodeIP = nodeMetaData.NodeIP
			requestIPMetaData.AssignedNodeName = nodeMetaData.NodeName
			hr.RequestIpToMetaData[requestIP] = requestIPMetaData
			hr.RequestIpCollection[nodeIp] = append(hr.RequestIpCollection[nodeIp], requestIP)
		}
	}
}

func (hr *HashRing) RemoveNode(nodeNameOfRemovableNode string) {
	hr.Lock()
	defer hr.Unlock()
	var indexOfRemovableNode int = -1
	var ipOfRemovableNode string
	for i := 0; i < len(hr.NodeMetaData); i++ {
		if hr.NodeMetaData[i].NodeName == nodeNameOfRemovableNode {
			indexOfRemovableNode = i
			ipOfRemovableNode = hr.NodeMetaData[i].NodeIP
			break
		}
	}

	hr.NodeMetaData = append(hr.NodeMetaData[:indexOfRemovableNode], hr.NodeMetaData[indexOfRemovableNode+1:]...)
	delete(hr.NodeNameToNodeMetaData, nodeNameOfRemovableNode)
	fmt.Println(hr.NodeMetaData, "after delete")
	requestIpCollection, isNode := hr.RequestIpCollection[ipOfRemovableNode]
	if isNode {
		newAssignedNodeMetaData := hr.getNodeMetaData(hr.RequestIpToMetaData[requestIpCollection[0]].RequestHash)
		for _, requestIp := range requestIpCollection {
			metaData := hr.RequestIpToMetaData[requestIp]
			metaData.AssignedNodeHash = newAssignedNodeMetaData.NodeHash

			if metaData.AssignedNodeHash == 0 {
				metaData.AssignedNodeHash = -1
			}

			metaData.AssignedNodeIP = newAssignedNodeMetaData.NodeIP
			if len(metaData.AssignedNodeIP) == 0 {
				metaData.AssignedNodeIP = "NA"
			}

			metaData.AssignedNodeName = newAssignedNodeMetaData.NodeName
			if len(metaData.AssignedNodeName) == 0 {
				metaData.AssignedNodeName = "NA"
			}

			hr.RequestIpToMetaData[requestIp] = metaData
			hr.RequestIpCollection[newAssignedNodeMetaData.NodeIP] = append(hr.RequestIpCollection[newAssignedNodeMetaData.NodeIP], requestIp)
		}
		delete(hr.RequestIpCollection, ipOfRemovableNode)
	}
	fmt.Println(hr.RequestIpToMetaData, "in remove node")
}

func (hr *HashRing) GetNode(requestIP string) bool {
	hr.RLock()
	defer hr.RUnlock()

	RequestIPHashValue := int32(crc32.ChecksumIEEE([]byte(requestIP))%RING_SIZE) + 1
	AssignedNodeMetaData := hr.getNodeMetaData(RequestIPHashValue)

	if AssignedNodeMetaData.NodeHash == 0 && len(AssignedNodeMetaData.NodeIP) == 0 && len(AssignedNodeMetaData.NodeName) == 0 {
		return false
	}
	hr.RequestIpToMetaData[requestIP] = RequestMetaData{RequestHash: RequestIPHashValue, AssignedNodeHash: AssignedNodeMetaData.NodeHash, AssignedNodeName: AssignedNodeMetaData.NodeName, AssignedNodeIP: AssignedNodeMetaData.NodeIP}
	fmt.Println(hr.RequestIpToMetaData, "in getnode")
	hr.RequestIpCollection[AssignedNodeMetaData.NodeIP] = append(hr.RequestIpCollection[AssignedNodeMetaData.NodeIP], requestIP)
	fmt.Println(hr.RequestIpCollection, "collection")
	return true
}
