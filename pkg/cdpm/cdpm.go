package cdpm

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	_ "github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/schollz/progressbar/v3"
	"log"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

type Cdp struct {
	NativeVlan byte   `json:"native-vlan"`
	DeviceId   string `json:"device-id"`
	PortId     string `json:"port-id"`
	VoiceVlan  byte   `json:"voice-vlan"`
}

const (
	Max      = 1                // maximum number of packets to get
	Timeout  = time.Second * 10 // time before timeout (10 seconds)
	PcapSize = 262144           // pcap size
)

var (
	current = 0 // current packet count
)

// Start starts listening for packets on a given interface
func Start(device string) {
	open(device)
}

func open(device string) {
	if device == "Media disconnected" {
		fmt.Println("Media disconnected. Exiting.")
		os.Exit(1)
	}

	var __cdp Cdp

	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	handle, err := pcap.OpenLive(device, PcapSize, true, pcap.ErrNotActive)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		handle.Close()
	}()

	if err := handle.SetBPFFilter(""); err != nil {
		log.Fatal(err)
	}

	packets := gopacket.NewPacketSource(handle, handle.LinkType()).Packets()

	// NOTE: it can take up 60 seconds or more to receive a CDP packet from a switch

	var since time.Duration
	var done bool

	start := time.Now()

	fmt.Println("Waiting for a CDP packet to arrive...")
	for pkt := range packets {
		if contains(pkt.Layers()) {
			since = time.Since(start)

			if cdpLayer := pkt.Layer(layers.LayerTypeCiscoDiscovery); cdpLayer != nil {
				cdp, _ := cdpLayer.(*layers.CiscoDiscovery)
				for _, v := range cdp.Values {
					switch v.Type.String() {
					case "Native VLAN":
						__cdp.NativeVlan = v.Value[len(v.Value)-1]
					case "VoIP VLAN Reply":
						__cdp.VoiceVlan = v.Value[len(v.Value)-1]
					case "Device ID":
						__cdp.DeviceId = string(v.Value)
					case "Port ID":
						__cdp.PortId = string(v.Value)
					}
					since = time.Since(start)
				}
				break
			}
		}
		if done {
			break
		}
	}

	fmt.Fprintf(tw, "Switch Name\tPort\tVLAN\tVoice VLAN\n")
	fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n",
		strings.Repeat("-", len("Switch Name")),
		strings.Repeat("-", len("Port")),
		strings.Repeat("-", len("VLAN")),
		strings.Repeat("-", len("Voice VLAN")),
	)
	fmt.Fprintf(tw, "%s\t%s\t%v\t%v\n\n", __cdp.DeviceId, __cdp.PortId, __cdp.NativeVlan, __cdp.VoiceVlan)
	fmt.Printf("took %.2f \n seconds", since.Seconds())

	tw.Flush()

	return

}

func contains(layers []gopacket.Layer) bool {
	for _, l := range layers {
		if strings.Contains(l.LayerType().String(), "Cisco") {
			return true
		}
	}

	return false
}

// DEAD:

// increment adds 1 to bar, increments current, and returns true/false if/when current == Max
func increment(bar *progressbar.ProgressBar) bool {
	bar.Add(1)
	current += 1
	if current == Max {
		return true
	}

	return false
}

// TODO: refactor into a goroutine and use CTX to terminate the Start function.

// returns true/false if the time elapsed is greater than or equal to Timeout
func timeout(since time.Duration) bool { return since >= Timeout }
