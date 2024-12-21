package main

import (
	"example/cmd"
	"example/config"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var cobraCmd = &cobra.Command{
	Use:   "tss-example",
	Short: `This is a TSS example`,
	PersistentPreRunE: func(c *cobra.Command, args []string) error {
		if err := viper.BindPFlags(c.Flags()); err != nil {
			return err
		}
		config.Initialization()
		return cmd.MonitorP2p()
	},
}

func init() {
	cobraCmd.PersistentFlags().String("config", "", "config file path")
	cobraCmd.AddCommand(cmd.DkgCmd)
	cobraCmd.AddCommand(cmd.SignerCmd)
	cobraCmd.AddCommand(cmd.ReshareCmd)
	cobraCmd.AddCommand(cmd.ServerCmd)
}

func main() {
	if err := cobraCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
