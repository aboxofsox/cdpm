package cdpm

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"log"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

// CDP holds the information we need from
// our CDP packet.
type CDP struct {
	NativeVlan byte   `json:"native-vlan"`
	DeviceId   string `json:"device-id"`
	PortId     string `json:"port-id"`
	VoiceVlan  byte   `json:"voice-vlan"`
}

// LLDP holds the information we need from
// the LLDP packet. LLDP is only sent if
// a switchport is trunked.
type LLDP struct {
	PortDescription   string `json:"port-id"`
	SystemName        string `json:"system-name"`
	SystemDescription string `json:"system-description"`
}

// cdpHandler handles our CDP packet.
// Returning the CDP struct.
func cdpHandler(layer gopacket.Layer) CDP {
	var cdp CDP
	cdpl, _ := layer.(*layers.CiscoDiscovery)

	for _, v := range cdpl.Values {
		switch v.Type.String() {
		case "Native VLAN":
			cdp.NativeVlan = v.Value[len(v.Value)-1]
		case "VoIP VLAN Reply":
			cdp.VoiceVlan = v.Value[len(v.Value)-1]
		case "Device ID":
			cdp.DeviceId = string(v.Value)
		case "Port ID":
			cdp.PortId = string(v.Value)
		}
	}

	return cdp
}

// lldpHandler handles the LLDP packet.
// Returning the LLDP struct.
func lldpHandler(layer gopacket.Layer) LLDP {
	var lldp LLDP

	lldpl, _ := layer.(*layers.LinkLayerDiscoveryInfo)

	lldp.PortDescription = lldpl.PortDescription
	lldp.SystemName = lldpl.SysName
	lldp.SystemDescription = lldpl.SysDescription

	return lldp
}

// pktHandler handles all the packets coming
// in from gopacket.NewPacketSource.Packets().
func pktHandler(packets chan gopacket.Packet) {
	var (
		cdp   CDP
		lldp  LLDP
		tw    = tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		since time.Duration
	)

	start := time.Now()
	for pkt := range packets {
		if cdpLayer := pkt.Layer(layers.LayerTypeCiscoDiscovery); cdpLayer != nil {
			fmt.Println("CDP Packet Received")
			cdp = cdpHandler(cdpLayer)
			printCdp(tw, cdp)
			since = time.Since(start)
			break
		} else if lldpLayer := pkt.Layer(layers.LayerTypeLinkLayerDiscovery); lldpLayer != nil {
			fmt.Println("LLDP Packet Received")
			lldp = lldpHandler(lldpLayer)
			printLldp(tw, lldp)
			since = time.Since(start)
			break
		} else {
			if timeout(since) {
				return
			}
		}
	}
}

// printCdp prints our CDP struct to the terminal.
func printCdp(tw *tabwriter.Writer, cdp CDP) {
	println()
	if _, err := fmt.Fprintf(
		tw,
		"Device Name\tPort\tVLAN\tVoice VLAN\n",
	); err != nil {
		log.Fatal(err)
	}

	if _, err := fmt.Fprintf(
		tw,
		"%s\t%s\t%s\t%s\n",
		strings.Repeat("-", len("Device Name")),
		strings.Repeat("-", len("Port")),
		strings.Repeat("-", len("VLAN")),
		strings.Repeat("-", len("Voice VLAN")),
	); err != nil {
		return
	}

	if _, err := fmt.Fprintf(
		tw,
		"%s\t%s\t%v\t%v\n\n",
		cdp.DeviceId,
		cdp.PortId,
		cdp.NativeVlan,
		cdp.VoiceVlan,
	); err != nil {
		log.Fatal(err)
	}

	if err := tw.Flush(); err != nil {
		return
	}
}

// printLldp prints the LLDP struct to the terminal.
func printLldp(tw *tabwriter.Writer, lldp LLDP) {
	println()
	if _, err := fmt.Fprintf(
		tw,
		"\nDevice Name\tPort\n",
	); err != nil {
		log.Fatal(err)
	}

	if _, err := fmt.Fprintf(
		tw,
		"%s\t%s\n\n",
		lldp.SystemName,
		lldp.PortDescription,
	); err != nil {
		log.Fatal(err)
	}

	if err := tw.Flush(); err != nil {
		return
	}
}
