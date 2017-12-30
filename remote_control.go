package LanGenius

import (
	"encoding/json"
	// "fmt"
	"net"
)

var (
	RemoteControlEnabled bool
)

func handleRemoteControlCmd(msg Msg) {
	b, _ := json.Marshal(msg)
	if RemoteControlEnabled {
		mEventHandler.OnRemoteControlCmdReceived(string(b))
	}
}
func SetRemoteControlStatus(b bool) {
	RemoteControlEnabled = b
	broadcastAddr, _ := net.ResolveUDPAddr("udp", "255.255.255.255"+DeamonPort)
	broData, _ := json.Marshal(Msg{Type: "LanGenius-Deamon", State: "Online", Port: mPort, Info: osInfo, RemoteControlStatus: RemoteControlEnabled})
	deamonConn.WriteToUDP(broData, broadcastAddr)
}
func SendRemoteControlCmd(data string) {
	msg := Msg{}
	json.Unmarshal([]byte(data), &msg)
	sendAddr, _ := net.ResolveUDPAddr("udp", msg.IP+DeamonPort)
	msg.Type = "LanGenius-RemoteControlCmd"
	b, _ := json.Marshal(msg)
	deamonConn.WriteToUDP(b, sendAddr)
}
