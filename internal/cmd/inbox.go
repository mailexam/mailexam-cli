package cmd

import (
	"github.com/mailexam/mailexam-cli/internal/api"
	"github.com/mailexam/mailexam-cli/internal/output"
	"github.com/spf13/cobra"
)

func newInboxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inbox",
		Short: "Manage inboxes",
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List inboxes for a project",
		RunE:  runInboxList,
	}
	listCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")

	getCmd := &cobra.Command{
		Use:   "get [uuid]",
		Short: "Get inbox by UUID",
		Args:  cobra.ExactArgs(1),
		RunE:  runInboxGet,
	}
	getCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")

	cmd.AddCommand(listCmd, getCmd)
	return cmd
}

func runInboxList(cmd *cobra.Command, args []string) error {
	client, cfg, err := loadClient()
	if err != nil {
		return err
	}
	if cfg.ProjectUUID == "" {
		return errProjectRequired()
	}

	inboxes, err := client.ListInboxes(cfg.ProjectUUID)
	if err != nil {
		return err
	}

	if jsonOutput {
		return output.PrintJSON(inboxes)
	}

	output.PrintInboxes(inboxes)
	return nil
}

func runInboxGet(cmd *cobra.Command, args []string) error {
	client, _, err := loadClient()
	if err != nil {
		return err
	}

	inbox, err := client.GetInbox(args[0])
	if err != nil {
		return err
	}

	if jsonOutput {
		return output.PrintJSON(inbox)
	}

	output.PrintInboxes([]api.Inbox{*inbox})
	return nil
}
