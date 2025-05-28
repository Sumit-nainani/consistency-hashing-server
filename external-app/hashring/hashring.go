package hashring

import (
	"hash/crc32"
	"hashing/utility"
	"sort"
	"sync"
)

// consistency hash ring size adjusted according to screen.
const (
	RING_SIZE = 100
)

var (
	once     sync.Once
	instance *HashRing
)

// Metadata for sending to UI for a particular request ip address , same as proto buff structure.
type RequestMetaData struct {
	RequestHash      int32
	AssignedNodeHash int32
	AssignedNodeIP   string
	AssignedNodeName string
}

// Metadata for sending to UI for a particular node ip address , same as proto buff structure.
type NodeMetaData struct {
	NodeHash int32
	NodeIP   string
	NodeName string
}

// Consistency hash ring component.
type HashRing struct {
	sync.RWMutex
	NodeMetaData           []NodeMetaData
	NodeNameToNodeMetaData map[string]NodeMetaData
	RequestIpToMetaData    map[string]RequestMetaData
	RequestIpCollection    map[string][]string // This filed is required for rebalancing the ring , while rebalancing the request ip to serve from different nodes.
}

// SingleTon design pattern.
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

// Utility function for finding the correct node for a particular client , using binary search algorithm.
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

// This function is used for modifying all metadatas and doing rebalancing after node upscaling.
// This function uses the name of node and the ip address , which is added in the cluster.
// The name and ip is coming from kube watcher.
func (hr *HashRing) AddNode(nodeName string, nodeIp string) {
	hr.Lock()
	defer hr.Unlock()

	// This is special case everytime if a node restarts inside the cluster then we do
	// not want to insert again the same node metadata
	if _, isAlreadyNodePresent := hr.NodeNameToNodeMetaData[nodeName]; isAlreadyNodePresent {
		return
	}

	// Calculating hash value of the ip address of node using crs32 check sum algorithm.
	// And finding the insert index for keeping the slice sorted.
	nodeHashValue := int32(crc32.ChecksumIEEE([]byte(nodeIp))%RING_SIZE) + 1

	insertPos := sort.Search(len(hr.NodeMetaData), func(i int) bool {
		return hr.NodeMetaData[i].NodeHash >= nodeHashValue
	})

	// modifying the metadas for added node.
	hr.NodeNameToNodeMetaData[nodeName] = NodeMetaData{NodeHash: nodeHashValue, NodeName: nodeName, NodeIP: nodeIp}

	hr.NodeMetaData = append(hr.NodeMetaData[:insertPos], append([]NodeMetaData{{
		NodeHash: nodeHashValue,
		NodeIP:   nodeIp,
		NodeName: nodeName,
	}}, hr.NodeMetaData[insertPos:]...)...)

	// Again when a new node comes then rebalancing must be done to reduce the load of the existing nodes.
	// Each time when rebalancing is done then a websocket event will be broadcasted for real time update in UI.
	for requestIP, requestIPMetaData := range hr.RequestIpToMetaData {
		nodeMetaData := hr.getNodeMetaData(requestIPMetaData.RequestHash)

		// if a old client should be served by new node from now then do rebalancing for that client.
		if nodeMetaData.NodeName == nodeName {
			oldNodeIp := requestIPMetaData.AssignedNodeIP
			requestIPCollection := hr.RequestIpCollection[oldNodeIp]

			// Removing clients from old nodes.
			for index, rerequestIPs := range requestIPCollection {
				if rerequestIPs == requestIP {
					hr.RequestIpCollection[oldNodeIp] = append(requestIPCollection[:index], requestIPCollection[index+1:]...)

					if len(hr.RequestIpCollection[oldNodeIp]) == 0 {
						delete(hr.RequestIpCollection, oldNodeIp)
					}

					break
				}
			}

			// Inserting new assigned node metadata.
			requestIPMetaData.AssignedNodeHash = nodeMetaData.NodeHash
			requestIPMetaData.AssignedNodeIP = nodeMetaData.NodeIP
			requestIPMetaData.AssignedNodeName = nodeMetaData.NodeName
			hr.RequestIpToMetaData[requestIP] = requestIPMetaData
			hr.RequestIpCollection[nodeIp] = append(hr.RequestIpCollection[nodeIp], requestIP)

			// Broadcast websocket event of rebalancing.
			utility.BroadcastRequestIPMetaData(requestIPMetaData.RequestHash, requestIPMetaData.AssignedNodeName, requestIPMetaData.AssignedNodeIP, requestIPMetaData.AssignedNodeHash)
		}
	}
}

// This function is used for modifying all metadatas and doing rebalancing after node downscaling.
// This function uses the name of node which is removed from the cluster. The name is coming from kube watcher.
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

	// Removing node meta data from hash map , which does not exist anymore.
	hr.NodeMetaData = append(hr.NodeMetaData[:indexOfRemovableNode], hr.NodeMetaData[indexOfRemovableNode+1:]...)
	delete(hr.NodeNameToNodeMetaData, nodeNameOfRemovableNode)

	// We have to find all request ip addresses which were served by removable node , because now rebalancing will work.
	requestIpCollection, isNode := hr.RequestIpCollection[ipOfRemovableNode]

	// If the deleted node was serving more than client then we have to find new assigned nodes using binary search algorithm each
	// time. And if there is no node available right now then we will send dummy data to the UI.
	// Each time when rebalancing is done then websocket event will be broadcasted to UI for real time UI change.
	if isNode {
		newAssignedNodeMetaData := hr.getNodeMetaData(hr.RequestIpToMetaData[requestIpCollection[0]].RequestHash)
		for _, requestIp := range requestIpCollection {
			metaData := hr.RequestIpToMetaData[requestIp]
			metaData.AssignedNodeHash = newAssignedNodeMetaData.NodeHash
			metaData.AssignedNodeIP = newAssignedNodeMetaData.NodeIP
			metaData.AssignedNodeName = newAssignedNodeMetaData.NodeName

			hr.RequestIpToMetaData[requestIp] = metaData
			hr.RequestIpCollection[newAssignedNodeMetaData.NodeIP] = append(hr.RequestIpCollection[newAssignedNodeMetaData.NodeIP], requestIp)

			// If no server/node is available then sending dummy data to UI.
			if !utility.IsNodeAvailable(metaData.AssignedNodeIP, metaData.AssignedNodeName, metaData.AssignedNodeHash) {
				metaData.AssignedNodeHash = -1
				metaData.AssignedNodeIP = "0.0.0.0"
				metaData.AssignedNodeName = "NA"
			}

			// Broadcast the rebalance event.
			utility.BroadcastRequestIPMetaData(metaData.RequestHash, metaData.AssignedNodeName, metaData.AssignedNodeIP, metaData.AssignedNodeHash)

		}
		delete(hr.RequestIpCollection, ipOfRemovableNode)
	}
}

// This function is used to assign a node to a particular ip address
// by binary search algorithm using lower bound approch.
// We are calculation ip hashes using crc32 checksum and finding lower bound
// for that request ip hash value.
func (hr *HashRing) GetNode(requestIP string) bool {
	hr.RLock()
	defer hr.RUnlock()

	// If a node is assigned to a request first time then for next requests we dont need to compute again.
	if _, isAlreadyAssignedNode := hr.RequestIpToMetaData[requestIP]; isAlreadyAssignedNode {
		return false
	}

	// Finding hash value and node metadata for a particulart client.
	RequestIPHashValue := int32(crc32.ChecksumIEEE([]byte(requestIP))%RING_SIZE) + 1
	AssignedNodeMetaData := hr.getNodeMetaData(RequestIPHashValue)

	hr.RequestIpToMetaData[requestIP] = RequestMetaData{RequestHash: RequestIPHashValue, AssignedNodeHash: AssignedNodeMetaData.NodeHash, AssignedNodeName: AssignedNodeMetaData.NodeName, AssignedNodeIP: AssignedNodeMetaData.NodeIP}

	// Record of all requests which is served by a particular node. It will be modified on each rebalance.
	hr.RequestIpCollection[AssignedNodeMetaData.NodeIP] = append(hr.RequestIpCollection[AssignedNodeMetaData.NodeIP], requestIP)

	return true
}
