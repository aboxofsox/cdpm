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

func Log(device string, elapse int) {
	if strings.ToLower(device) == "media disconnected" {
		fmt.Println(ErrMediaDisconnected)
		os.Exit(1)
	}

	var (
		i         int
		cdpCount  int
		lldpCount int
		since     time.Duration
	)

	file, err := os.OpenFile("cdp_packets.log", os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	tw := tabwriter.NewWriter(file, 0, 0, 3, ' ', 0)

	for i < elapse {
		handle, err := pcap.OpenLive(device, PcapSize, true, pcap.ErrNotActive)
		if err != nil {
			log.Fatal(err)
		}

		if err := handle.SetBPFFilter(""); err != nil {
			log.Fatal(err)
		}

		packets := gopacket.NewPacketSource(handle, handle.LinkType()).Packets()

		start := time.Now()
		for pkt := range packets {
			if cdpLayer := pkt.Layer(layers.LayerTypeCiscoDiscovery); cdpLayer != nil {
				cdpCount++
				fmt.Printf("CDP: %d\n", cdpCount)
				cdp := cdpHandler(cdpLayer)
				printCdp(tw, cdp)
				since = time.Since(start)
				break
			} else if lldpLayer := pkt.Layer(layers.LayerTypeLinkLayerDiscovery); lldpLayer != nil {
				lldpCount++
				fmt.Printf("LLDP: %d\n", lldpCount)
				lldp := lldpHandler(lldpLayer)
				printLldp(tw, lldp)
				since = time.Since(start)
				break
			} else {
				if timeout(since) {
					return
				}

			}
		}

		handle.Close()
	}
}
