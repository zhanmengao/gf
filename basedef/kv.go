package basedef

import (
	"bytes"
	"fmt"
)

type Event int

const (
	EventNone Event = iota
	EventCreate
	EventPut
	EventDelete
	EventClose
	EventConn
)

type Node struct {
	Key   string
	Value *NodeInfo
}

type RouteEvent struct {
	Type   Event
	KV     *Node
	PreKV  *Node
	Prefix string
}

func (ke *RouteEvent) String() string {
	buf := bytes.Buffer{}
	fmt.Fprintf(&buf, "%s event = %d\t", ke.Prefix, ke.Type)
	if ke.KV != nil {
		fmt.Fprintf(&buf, "current %s:%v  \t", ke.KV.Key, ke.KV.Value)
	}
	if ke.PreKV != nil {
		fmt.Fprintf(&buf, "preview %s:%v  \t", ke.PreKV.Key, ke.PreKV.Value)
	}
	return buf.String()
}

type NodeInfo struct {
	TCPAddr     string `json:"tcp_addr,omitempty"`
	SrvName     string `json:"srv_name,omitempty"`
	Status      int    `json:"status,omitempty"`
	Sid         int    `json:"sid,omitempty"`
	Version     string `json:"version,omitempty"`
	ListenMode  int    `json:"listen_mode,omitempty"`
	HostName    string `json:"host_name,omitempty"`
	GRPCAddr    string `json:"grpc_addr,omitempty"`
	Pid         int    `json:"pid,omitempty"`
	WsAddr      string `json:"ws_addr,omitempty"`
	HttpAddr    string `json:"http_addr"`
	BaseVersion string `json:"base_version"`
	SrvVersion  string `json:"srv_version"`
}
