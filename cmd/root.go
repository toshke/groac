package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type Config struct {
	Platform string
}

var (
	config  Config
	rootCmd = &cobra.Command{
		Use:   "groac",
		Short: "Gitlab runner manager for 202x ",
		Long: `Gitlab runner manager for major cloud providers
			Complete documentation is available at http://github.com/toshke/groac`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Display groac version",
		Run: func(cmd *cobra.Command, args []string) {
			version, _ := os.LookupEnv("GROAC_VERSION")
			fmt.Printf("groac version %v\n", version)
		},
	}
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "config stage custom gitlab runner executor",
		Run: func(cmd *cobra.Command, args []string) {
			groacCleanup()
		},
	}
	prepareCmd = &cobra.Command{
		Use:   "prepare",
		Short: "prepare stage custom gitlab runner executor",
		Run: func(cmd *cobra.Command, args []string) {
			groacPrepare()
		},
	}
	stepCmd = &cobra.Command{
		Use:   "step",
		Short: "step/exec stage custom gitlab runner executor",
		Run: func(cmd *cobra.Command, args []string) {
			groacStep()
		},
	}
	cleanupCmd = &cobra.Command{
		Use:   "cleanup",
		Short: "cleanup stage custom gitlab runner executor",
		Run: func(cmd *cobra.Command, args []string) {
			groacStep()
		},
	}
)

func initConfig() {
	fmt.Printf("Running on platform %v\n", config.Platform)
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&config.Platform, "platform", "p", "aws", "Runner platform aws|gcp|azure")
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(prepareCmd)
	rootCmd.AddCommand(stepCmd)
	rootCmd.AddCommand(cleanupCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
