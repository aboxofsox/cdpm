package main

import (
	"cdpm/pkg/wininterface"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	_ "github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"net"
	"os/exec"
	"strings"
)

func networkInterfaces() []net.Interface {
	var iFaceSlice []net.Interface
	iFaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for _, iFace := range iFaces {
		if iFace.Name != "" {
			iFaceSlice = append(iFaceSlice, iFace)
		}
	}

	return iFaceSlice

}

func networkInterfaceByName(name string) net.Interface {
	_interfaces, err := net.Interfaces()
	if err != nil {
		return net.Interface{}
	}

	for _, _interface := range _interfaces {
		if _interface.Name == name {
			return _interface
		}
	}

	return net.Interface{}
}

const (
	snapLen = 262144
)

func start(device string) {
	handle, err := pcap.OpenLive(device, 262144, true, pcap.ErrNotActive)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	if err := handle.SetBPFFilter(""); err != nil {
		log.Fatal(err)
	}

	packets := gopacket.NewPacketSource(handle, handle.LinkType()).Packets()
	for pkt := range packets {
		if cdpLayer := pkt.Layer(layers.LayerTypeCiscoDiscovery); cdpLayer != nil {
			cdp, _ := cdpLayer.(*layers.CiscoDiscovery)
			cdpv := cdp.Values
			for _, v := range cdpv {
				v.Type.String()
			}
		}
	}
}

func getMac() string {
	cmd := exec.Command("getmac", "/FO", "list", "/V")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}

	return string(output)

}

func fmtTransportName(transportName string) {
	tpSplit := strings.Split(transportName, "\\")
	for _, tp := range tpSplit {
		fmt.Println(tp)
	}

}

func main() {
	tpName := getMac()
	macs := wininterface.Parse(tpName)
	tpn := wininterface.GetTransportByName("Ethernet", macs)
	fmt.Println(tpn)
	start(tpn)

}
