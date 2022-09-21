package cmd

import (
	"cdpm/pkg/cdpm"
	"github.com/aboxofsox/wininterface"
	"github.com/spf13/cobra"
)

var (
	outfile      string
	netInterface string
	timeout      int
	log          bool
)

func init() {
	rootCmd.AddCommand(start)
	rootCmd.PersistentFlags().StringVarP(&outfile, "outfile", "o", "results.md", "Output results to file.")
	rootCmd.PersistentFlags().StringVarP(&netInterface, "interface", "i", "", "Defines the interface to listen to.")
	rootCmd.PersistentFlags().BoolVarP(&log, "log", "l", false, "Log received packets to file.")
	rootCmd.PersistentFlags().IntVarP(&timeout, "duration", "d", 30, "How long to listen for packets.")
}

var start = &cobra.Command{
	Use:   "start",
	Short: "start cdpm",
	Run: func(cmd *cobra.Command, args []string) {
		var tpName string
		win := wininterface.GetMac()
		names := win.Parse()
		for _, n := range names {
			if netInterface == "" {
				netInterface = "Ethernet"
			}
			if netInterface == n.ConnectionName {
				tpName = n.TransportName
			}

			if log {
				cdpm.Log(tpName, 0)
			} else {
				cdpm.Start(tpName)
			}
		}
	},
}
