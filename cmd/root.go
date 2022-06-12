package cmd

import (
	"github.com/hongqchen/imagectl/utils/log"
	"github.com/spf13/cobra"
)

func Execute() error {
	enableDebug := "false"

	var cmdSync = &cobra.Command{
		Use:  "sync registry/[namespace]/image_name:tag",
		Args: cobra.MinimumNArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			enableDebug = cmd.Flag("debug").Value.String()
			log.InitLogger(enableDebug)
		},
		Run: func(cmd *cobra.Command, args []string) {
			Start(args)
		},
	}
	cmdSync.Flags().Bool("debug", false, "Enable debug mode")

	var rootCmd = &cobra.Command{
		Use: "imagectl",
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true},
		DisableSuggestions: true,
	}

	rootCmd.AddCommand(cmdSync)
	rootCmd.SetHelpTemplate(`
Usage:
  imagectl [command]

Basic Commands:
  sync	Manually synchronize external network images to Alibaba Cloud Mirror

Use "imagectl [command] --help" for more information about a command.
`)
	return rootCmd.Execute()
}
