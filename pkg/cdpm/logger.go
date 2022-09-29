package cdpm

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

func Log(device string, duration int) {
	if strings.ToLower(device) == "media disconnected" {
		fmt.Println(ErrMediaDisconnected)
		os.Exit(1)
	}

	var (
		cdpCount  int
		lldpCount int
		since     time.Duration
	)

	file, err := os.OpenFile("cdp_packets.log", os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	tw := tabwriter.NewWriter(file, 0, 0, 3, ' ', 0)

	handle, err := pcap.OpenLive(device, PcapSize, true, pcap.ErrNotActive)
	if err != nil {
		log.Fatal(err)
	}

	if err := handle.SetBPFFilter(""); err != nil {
		log.Fatal(err)
	}

	packets := gopacket.NewPacketSource(handle, handle.LinkType()).Packets()
	start := time.Now()
	end := start.Add(time.Second * time.Duration(duration))
	for pkt := range packets {
		if time.Now().UnixNano() >= end.UnixNano() {
			fmt.Println("timeout reached")
			return
		}

		if cdpLayer := pkt.Layer(layers.LayerTypeCiscoDiscovery); cdpLayer != nil {
			cdpCount++
			fmt.Printf("CDP: %d\n", cdpCount)
			cdp := cdpHandler(cdpLayer)
			printCdp(tw, cdp, time.Now())
			since = time.Since(start)
		} else if lldpLayer := pkt.Layer(layers.LayerTypeLinkLayerDiscovery); lldpLayer != nil {
			lldpCount++
			fmt.Printf("LLDP: %d\n", lldpCount)
			lldp := lldpHandler(lldpLayer)
			printLldp(tw, lldp, time.Now())
			since = time.Since(start)

		}

		if timeout(since) {
			return
		}

	}

	handle.Close()
}
