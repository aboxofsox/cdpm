package cmd

import (
	"cdpm/pkg/cdpm"
	"cdpm/pkg/wininterface"
	"github.com/spf13/cobra"
)

var (
	outfile      string
	netInterface string
)

func init() {
	rootCmd.AddCommand(start)
	rootCmd.PersistentFlags().StringVarP(&outfile, "outfile", "o", "results.md", "Output results to file.")
	rootCmd.PersistentFlags().StringVarP(&netInterface, "interface", "i", "", "Defines the interface to listen to.")
}

var start = &cobra.Command{
	Use:   "start",
	Short: "start cdpm",
	Run: func(cmd *cobra.Command, args []string) {
		winint := wininterface.GetMac()
		name := wininterface.Parse(winint)
		for _, n := range name {
			if netInterface == "" {
				netInterface = "Ethernet"
			}
			if netInterface == n.ConnectionName {
				cdpm.Start(n.TransportName)
			}
		}
	},
}
