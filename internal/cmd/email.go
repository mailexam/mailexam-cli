package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/mailexam/mailexam-cli/internal/api"
	"github.com/mailexam/mailexam-cli/internal/output"
	"github.com/spf13/cobra"
)

func newEmailCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "email",
		Short: "Work with emails",
	}

	cmd.AddCommand(newEmailListCmd())
	cmd.AddCommand(newEmailGetCmd())
	cmd.AddCommand(newEmailWaitCmd())
	cmd.AddCommand(newEmailAssertCmd())
	cmd.AddCommand(newEmailDeleteCmd())
	cmd.AddCommand(newAttachmentCmd())

	return cmd
}

func newEmailListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List emails in an inbox",
		RunE:  runEmailList,
	}
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
	return cmd
}

func runEmailList(cmd *cobra.Command, args []string) error {
	client, cfg, err := loadClient()
	if err != nil {
		return err
	}

	inboxUUID, err := client.ResolveInbox(cfg.ProjectUUID, cfg.InboxUUID)
	if err != nil {
		return err
	}

	emails, err := client.ListEmails(inboxUUID)
	if err != nil {
		return err
	}

	if jsonOutput {
		return output.PrintJSON(emails)
	}

	output.PrintEmails(emails)
	return nil
}

func newEmailGetCmd() *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:   "get [uuid]",
		Short: "Get email by UUID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, _, err := loadClient()
			if err != nil {
				return err
			}

			email, err := client.GetEmail(args[0])
			if err != nil {
				return err
			}

			if format == "json" || jsonOutput {
				return output.PrintJSON(email)
			}

			return output.PrintEmail(email, format)
		},
	}

	cmd.Flags().StringVar(&format, "format", "json", "Output format: json, text, html, raw, uuid")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON (same as --format json)")
	return cmd
}

func newEmailWaitCmd() *cobra.Command {
	var (
		subject      string
		to           string
		from         string
		subjectExact bool
		timeout      int
		interval     int
		format       string
	)

	cmd := &cobra.Command{
		Use:   "wait",
		Short: "Wait until an email matching filters arrives",
		Long:  "Polls the inbox until an email matches --subject, --to, or --from, then prints it.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if subject == "" && to == "" && from == "" {
				return fmt.Errorf("at least one filter is required: --subject, --to, or --from")
			}

			client, cfg, err := loadClient()
			if err != nil {
				return err
			}

			email, err := client.WaitForEmail(api.WaitOptions{
				Filter: api.EmailFilter{
					Subject:      subject,
					To:           to,
					From:         from,
					SubjectExact: subjectExact,
				},
				Project:  cfg.ProjectUUID,
				Inbox:    cfg.InboxUUID,
				Timeout:  time.Duration(timeout) * time.Second,
				Interval: time.Duration(interval) * time.Second,
			})
			if err != nil {
				return err
			}

			if format == "json" || jsonOutput {
				return output.PrintJSON(email)
			}

			return output.PrintEmail(email, format)
		},
	}

	cmd.Flags().StringVar(&subject, "subject", "", "Filter by subject (substring)")
	cmd.Flags().BoolVar(&subjectExact, "subject-exact", false, "Match subject exactly")
	cmd.Flags().StringVar(&to, "to", "", "Filter by recipient (substring)")
	cmd.Flags().StringVar(&from, "from", "", "Filter by sender (substring)")
	cmd.Flags().IntVar(&timeout, "timeout", 30, "Timeout in seconds")
	cmd.Flags().IntVar(&interval, "interval", 2, "Poll interval in seconds")
	cmd.Flags().StringVar(&format, "format", "json", "Output format: json, text, html, raw, uuid")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")

	return cmd
}

func newEmailAssertCmd() *cobra.Command {
	var (
		subject          string
		to               string
		from             string
		subjectExact     bool
		contains         []string
		matches          string
		attachmentCount  int
		requireAttach    bool
		attachmentCountSet bool
		timeout          int
		interval         int
	)

	cmd := &cobra.Command{
		Use:   "assert",
		Short: "Wait for an email and verify its content",
		Long:  "Exits 0 if all checks pass, 1 if assertion fails, 3 on timeout.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if subject == "" && to == "" && from == "" {
				return fmt.Errorf("at least one filter is required: --subject, --to, or --from")
			}

			client, cfg, err := loadClient()
			if err != nil {
				return err
			}

			email, err := client.WaitForEmail(api.WaitOptions{
				Filter: api.EmailFilter{
					Subject:      subject,
					To:           to,
					From:         from,
					SubjectExact: subjectExact,
				},
				Project:  cfg.ProjectUUID,
				Inbox:    cfg.InboxUUID,
				Timeout:  time.Duration(timeout) * time.Second,
				Interval: time.Duration(interval) * time.Second,
			})
			if err != nil {
				return err
			}

			body := email.TextBody()
			if body == "" {
				body = email.HTMLBody()
			}

			for _, needle := range contains {
				if !strings.Contains(body, needle) {
					return fmt.Errorf("body does not contain %q", needle)
				}
			}

			if matches != "" {
				re, err := regexp.Compile(matches)
				if err != nil {
					return fmt.Errorf("invalid --matches regex: %w", err)
				}
				if !re.MatchString(body) {
					return fmt.Errorf("body does not match regex %q", matches)
				}
			}

			if requireAttach && !email.Attachments {
				return fmt.Errorf("email has no attachments")
			}

			if attachmentCountSet && len(email.Body.Attachments) != attachmentCount {
				return fmt.Errorf("expected %d attachments, got %d", attachmentCount, len(email.Body.Attachments))
			}

			fmt.Println("OK")
			return nil
		},
	}

	cmd.Flags().StringVar(&subject, "subject", "", "Filter by subject (substring)")
	cmd.Flags().BoolVar(&subjectExact, "subject-exact", false, "Match subject exactly")
	cmd.Flags().StringVar(&to, "to", "", "Filter by recipient (substring)")
	cmd.Flags().StringVar(&from, "from", "", "Filter by sender (substring)")
	cmd.Flags().StringArrayVar(&contains, "contains", nil, "Body must contain text (repeatable)")
	cmd.Flags().StringVar(&matches, "matches", "", "Body must match regular expression")
	cmd.Flags().BoolVar(&requireAttach, "require-attachments", false, "Email must have attachments")
	cmd.Flags().IntVar(&attachmentCount, "attachment-count", 0, "Expected number of attachments")
	cmd.Flags().IntVar(&timeout, "timeout", 30, "Timeout in seconds")
	cmd.Flags().IntVar(&interval, "interval", 2, "Poll interval in seconds")

	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		attachmentCountSet = cmd.Flags().Changed("attachment-count")
		return nil
	}

	return cmd
}

func newEmailDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete [uuid]",
		Short: "Delete email by UUID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, _, err := loadClient()
			if err != nil {
				return err
			}
			return client.DeleteEmail(args[0])
		},
	}
}

func newAttachmentCmd() *cobra.Command {
	var outputPath string

	cmd := &cobra.Command{
		Use:   "attachment",
		Short: "Download email attachments",
	}

	downloadCmd := &cobra.Command{
		Use:   "download [email-uuid]",
		Short: "Download attachment by CID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cid, err := cmd.Flags().GetString("cid")
			if err != nil {
				return err
			}
			if cid == "" {
				return fmt.Errorf("--cid is required")
			}

			client, _, err := loadClient()
			if err != nil {
				return err
			}

			data, err := client.DownloadAttachment(args[0], cid)
			if err != nil {
				return err
			}

			if outputPath == "" || outputPath == "-" {
				_, err = os.Stdout.Write(data)
				return err
			}

			return os.WriteFile(outputPath, data, 0o644)
		},
	}

	downloadCmd.Flags().String("cid", "", "Attachment CID (required)")
	downloadCmd.Flags().StringVarP(&outputPath, "output", "o", "-", "Output file path (- for stdout)")
	_ = downloadCmd.MarkFlagRequired("cid")

	cmd.AddCommand(downloadCmd)
	return cmd
}

func errProjectRequired() error {
	return fmt.Errorf("project is required (flag --project or env MAILEXAM_PROJECT_UUID)")
}
