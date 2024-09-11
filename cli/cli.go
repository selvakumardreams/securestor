package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "bluenoise",
		Short: "bluenoise is a Vulnerability Scanning CLI",
		Long:  `bluenoise is a CLI tool that can be used to scan for vulnerabilities in your applications.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("bluenoise - Vulnerability Scanning!")
		},
	}

	var echoCmd = &cobra.Command{
		Use:   "echo [message]",
		Short: "Echo prints the provided message",
		Long:  `Echo prints the provided message to the console.`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(args[0])
		},
	}

	var scanCmd = &cobra.Command{
		Use:   "scan",
		Short: "Scan commands",
		Long:  `Scan commands for various scanning operations.`,
	}

	rootCmd.AddCommand(echoCmd)
	rootCmd.AddCommand(scanCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
