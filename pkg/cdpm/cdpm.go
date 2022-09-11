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
	"text/tabwriter"
	"time"
)

const (
	Max      = 100              // maximum number of packets to get
	Timeout  = time.Second * 10 // time before timeout (10 seconds)
	PcapSize = 262144           // pcap size
)

var (
	current = 0 // current packet count
)

// Start starts listening for packets on a given interface
func Start(device string) {
	handle, err := pcap.OpenLive(device, PcapSize, true, Timeout)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		handle.Close()
	}()

	if err := handle.SetBPFFilter("port 443"); err != nil {
		log.Fatal(err)
	}

	packets := gopacket.NewPacketSource(handle, handle.LinkType()).Packets()

	cdpProcessed := make([]*Cdp, Max)

	bar := progressbar.Default(int64(Max))

	start := time.Now()
	for pkt := range packets {
		if cdpLayer := pkt.Layer(layers.LayerTypeCiscoDiscovery); cdpLayer != nil {
			cdp, _ := cdpLayer.(*layers.CiscoDiscoveryInfo)
			if b := increment(bar); b {
				break
			}
			cdpProcessed = append(cdpProcessed, CdpHandler(cdp))
		} else {
			since := time.Since(start)
			if tmo := timeout(since); tmo {
				fmt.Println(" no CDP packets were received")
				os.Exit(1)
			}
		}
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	if _, err := fmt.Fprintf(tw, "Device ID\tVLAN\tPort ID\n"); err != nil {
		fmt.Println(err)
		return
	}

	for _, cdp := range cdpProcessed {
		cdp.Print(tw)
	}
}

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
