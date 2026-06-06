package cmd

import (
	"github.com/mailexam/mailexam-cli/internal/api"
	"github.com/mailexam/mailexam-cli/internal/output"
	"github.com/spf13/cobra"
)

var jsonOutput bool

func newProjectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Manage projects",
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List projects",
		RunE:  runProjectList,
	}
	listCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")

	getCmd := &cobra.Command{
		Use:   "get [uuid]",
		Short: "Get project by UUID",
		Args:  cobra.ExactArgs(1),
		RunE:  runProjectGet,
	}
	getCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")

	cmd.AddCommand(listCmd, getCmd)
	return cmd
}

func runProjectList(cmd *cobra.Command, args []string) error {
	client, _, err := loadClient()
	if err != nil {
		return err
	}

	projects, err := client.ListProjects()
	if err != nil {
		return err
	}

	if jsonOutput {
		return output.PrintJSON(projects)
	}

	output.PrintProjects(projects)
	return nil
}

func runProjectGet(cmd *cobra.Command, args []string) error {
	client, _, err := loadClient()
	if err != nil {
		return err
	}

	project, err := client.GetProject(args[0])
	if err != nil {
		return err
	}

	if jsonOutput {
		return output.PrintJSON(project)
	}

	output.PrintProjects([]api.Project{*project})
	return nil
}
