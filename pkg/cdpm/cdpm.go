package cdpm

import (
	"fmt"
	"github.com/google/gopacket"
	_ "github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"os"
	"strings"
	"time"
)

const (
	Timeout  = time.Second * 10 // time before timeout (10 seconds)
	PcapSize = 262144           // pcap size
)

// Start starts listening for packets on a given interface
func Start(device string) {
	if strings.ToLower(device) == "media disconnected" {
		fmt.Println("Media disconnected. Exiting.")
		os.Exit(1)
	}

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

	fmt.Println("Waiting for a CDP packet to arrive...")

	pktHandler(packets)

	return
}

// timeout returns true/false if the time elapsed is greater than or equal to Timeout
func timeout(since time.Duration) bool { return since >= Timeout }
