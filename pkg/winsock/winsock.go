package winsock

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	_ "github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
)

func Start(device string) {
	handle, err := pcap.OpenLive(device, 262144, true, pcap.ErrNotActive)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	if err := handle.SetBPFFilter(""); err != nil {
		log.Fatal(err)
	}

	var (
		eth layers.Ethernet
		ip4 layers.IPv4
		ip6 layers.IPv6
		tpc layers.TCP
		cdp layers.CiscoDiscovery
	)

	parser := gopacket.NewDecodingLayerParser(layers.LayerTypeCiscoDiscovery, &cdp)

	//packets := gopacket.NewPacketSource(handle, handle.LinkType()).Packets()
	//for pkt := range packets {
	//	if cdpLayer := pkt.Layer(layers.LayerTypeCiscoDiscovery); cdpLayer != nil {
	//		cdp, _ := cdpLayer.(*layers.CiscoDiscovery)
	//
	//	}
	//}
}
