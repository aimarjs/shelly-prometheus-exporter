package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	cmd := &cobra.Command{
		Use: "test",
		RunE: func(cmd *cobra.Command, args []string) error {
			viper.BindPFlags(cmd.Flags())
			
			devices := viper.GetStringSlice("shelly-devices")
			fmt.Printf("Devices: %v\n", devices)
			fmt.Printf("Devices count: %d\n", len(devices))
			
			return nil
		},
	}
	
	cmd.Flags().StringSlice("shelly-devices", []string{}, "List of Shelly device URLs")
	
	cmd.Execute()
}
