package cmd

import (
	"os"

	"github.com/lngramos/chopper/internal/llm"
	"github.com/spf13/cobra"
)

var unsafeMode bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "chopper",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&unsafeMode, "unsafe", false, "Disable safe mode (no confirmation for tool execution)")

	client := llm.NewOllamaClient("http://localhost:11434")
	rootCmd.AddCommand(NewReplCommand(client))
	rootCmd.AddCommand(NewChatCommand(client))
}
