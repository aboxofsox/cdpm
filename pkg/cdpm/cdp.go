package cdpm

import (
	"encoding/json"
	"fmt"
	"github.com/google/gopacket/layers"
	"os"
	"text/tabwriter"
)

type Cdp struct {
	DeviceID string `json:"device-id,omitempty"`
	VLAN     uint16 `json:"vlan,omitempty"`
	PortID   string `json:"port-id,omitempty"`
}

func CdpHandler(cdp *layers.CiscoDiscoveryInfo) *Cdp {
	return &Cdp{
		DeviceID: cdp.DeviceID,
		VLAN:     cdp.NativeVLAN,
		PortID:   cdp.PortID,
	}
}

func (cdp *Cdp) Print(tw *tabwriter.Writer) {
	if _, err := fmt.Fprintf(tw, "%s\t%d\t%s\n", cdp.DeviceID, cdp.VLAN, cdp.PortID); err != nil {
		fmt.Println(err)
		return
	}
}

func (cdp *Cdp) WriteFile(path string) {
	b, err := json.MarshalIndent(cdp, "", " ")
	if err != nil {
		fmt.Println(err)
		return
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	if _, err := file.Write(b); err != nil {
		fmt.Println(err)
		return
	}
}
