package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/mailexam/mailexam-cli/internal/api"
	"github.com/mailexam/mailexam-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	flagBase    string
	flagToken   string
	flagProject string
	flagInbox   string
)

var rootCmd = &cobra.Command{
	Use:   "mailexam",
	Short: "CLI for Mailexam REST API",
	Long:  "Command-line tool for Mailexam test mail sandbox: projects, inboxes, emails, and CI/CD assertions.",
	SilenceUsage: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(exitCode(err))
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagBase, "base", "", "API base URL (env MAILEXAM_API_BASE)")
	rootCmd.PersistentFlags().StringVar(&flagToken, "token", "", "API token (env MAILEXAM_API_TOKEN)")
	rootCmd.PersistentFlags().StringVar(&flagProject, "project", "", "Project UUID (env MAILEXAM_PROJECT_UUID)")
	rootCmd.PersistentFlags().StringVar(&flagInbox, "inbox", "", "Inbox UUID (env MAILEXAM_INBOX_UUID)")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(newProjectCmd())
	rootCmd.AddCommand(newInboxCmd())
	rootCmd.AddCommand(newEmailCmd())
}

func loadClient() (*api.Client, config.Config, error) {
	cfg, err := config.Load(flagBase, flagToken, flagProject, flagInbox)
	if err != nil {
		return nil, cfg, err
	}
	return api.NewClient(cfg.BaseURL, cfg.Token), cfg, nil
}

func exitCode(err error) int {
	if err == nil {
		return 0
	}

	switch {
	case isAuthError(err):
		return 2
	case isTimeoutError(err):
		return 3
	default:
		return 1
	}
}

func isAuthError(err error) bool {
	if apiErr, ok := err.(*api.APIError); ok && apiErr.StatusCode == 403 {
		return true
	}
	return false
}

func isTimeoutError(err error) bool {
	return strings.Contains(err.Error(), "not found within")
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("mailexam-cli dev")
	},
}
