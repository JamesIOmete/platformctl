package fleet

import (
	"fmt"
	"strings"
)

// Device represents a mock fleet device.
type Device struct {
	ID        string `json:"id"`
	Model     string `json:"model"`
	State     string `json:"state"`
	IP        string `json:"ip"`
	Battery   int    `json:"battery"`
	Firmware  string `json:"firmware"`
	LastSeen  string `json:"last_seen"` // ISO8601 string
}

// ListDevices returns mocked devices.
func ListDevices() []Device {
	return []Device{
		{ID: "robot-001", Model: "digit", State: "online"},
		{ID: "robot-002", Model: "digit", State: "maintenance"},
		{ID: "edge-101", Model: "edge-node", State: "offline"},
	}
}

// FormatDevices renders a simple table.
func FormatDevices(devices []Device) string {
	if len(devices) == 0 {
		return "No devices found\n"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "ID          \tMODEL   \tSTATE\n")
	fmt.Fprintf(&b, "------------\t--------\t-----------\n")
	for _, d := range devices {
		fmt.Fprintf(&b, "%s\t%s\t%s\n", d.ID, d.Model, d.State)
	}
	return b.String()
}
