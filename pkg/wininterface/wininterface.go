package wininterface

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"
)

// Mac represents the output from getmac
type Mac struct {
	ConnectionName  string
	NetworkAdapter  string
	PhysicalAddress string
	TransportName   string
}

// Because Windows
const (
	CR = "\r" // carriage return
	LF = "\n" // line feed
)

// GetMac runs the getmac Windows command and returns its output.
func GetMac() string {
	cmd := exec.Command("getmac", "/FO", "list", "/V")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}

	return string(output)
}

// chunkSlice splits a slice by "chunk" given a chunk size.
func chunkSlice(slice []string, size int) [][]string {
	var chunks [][]string

	for i := 0; i < len(slice); i += size {
		end := i + size

		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

// chunkMap splits a map by "chunk" given a chunk size
func chunkMap(mp []map[string]string, size int) [][]map[string]string {
	var chunks [][]map[string]string
	for i := 0; i < len(mp); i += size {
		end := i + size
		if end > len(mp) {
			end = len(mp)
		}

		chunks = append(chunks, mp[i:end])
	}

	return chunks
}

// Parse takes the output of GetMac as a list and parses it for use, because Windows
func Parse(list string) []Mac {
	var amp []map[string]string

	listSplit := strings.Split(list, CR+LF)
	listChunks := chunkSlice(listSplit, 5)

	for i := range listChunks {
		mp := make(map[string]string)
		for j := range listChunks[i] {
			v := listChunks[i][j]

			vs := strings.Split(v, ":")
			if len(vs) == 2 {
				mp[vs[0]] = strings.Trim(vs[1], " ")
			}
		}
		amp = append(amp, mp)
	}

	dataChunks := chunkMap(amp, 4)

	var macs []Mac

	for i := range dataChunks {
		var mac Mac

		for _, chunk := range dataChunks[i] {

			if chunk["Connection Name"] != "" {
				mac.ConnectionName = chunk["Connection Name"]
				mac.NetworkAdapter = chunk["Network Adapter"]
				mac.PhysicalAddress = chunk["Physical Address"]

				tp := strings.ReplaceAll(chunk["Transport Name"], "Tcpip", "NPF")

				mac.TransportName = tp

				macs = append(macs, mac)
			}

		}
	}

	return macs

}

// Print pretty=prints the results []Mac
func Print(macs []Mac) {
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	_, err := fmt.Fprintf(tw, "Connection Name\tNetwork Adapter\tPhysical Address\tTransport Name\n")
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, m := range macs {
		_, err := fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", m.ConnectionName, m.NetworkAdapter, m.PhysicalAddress, m.TransportName)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if err := tw.Flush(); err != nil {
		log.Fatal(err)
	}
}

// GetTransportByName filters []Mac by name.
func GetTransportByName(name string, macs []Mac) string {
	for _, m := range macs {
		if m.ConnectionName == name {
			return m.TransportName
		}
	}

	return ""
}
