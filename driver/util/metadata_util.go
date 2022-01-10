package util

import "os"

const (
	NodeName = "MY_NODE_NAME"
	NodeZone = "MY_NODE_ZONE"
)

func GetNodeName() (nodeName, nodeZone string) {
	nodeName = os.Getenv(NodeName)
	nodeZone = os.Getenv(NodeZone)
	return
}

