package output

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/mailexam/mailexam-cli/internal/api"
)

func PrintJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func PrintProjects(projects []api.Project) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "UUID\tNAME")
	for _, p := range projects {
		fmt.Fprintf(w, "%s\t%s\n", p.UUID, p.Name)
	}
	w.Flush()
}

func PrintInboxes(inboxes []api.Inbox) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "UUID\tNAME\tDEFAULT\tPROJECT")
	for _, i := range inboxes {
		fmt.Fprintf(w, "%s\t%s\t%t\t%s\n", i.UUID, i.Name, i.IsDefault, i.ProjectUUID)
	}
	w.Flush()
}

func PrintEmails(emails []api.EmailShort) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "UUID\tDATE\tFROM\tSUBJECT")
	for _, e := range emails {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", e.UUID, e.Date, e.From, e.Subject)
	}
	w.Flush()
}

func PrintEmail(email *api.EmailDetailed, format string) error {
	switch format {
	case "json":
		return PrintJSON(email)
	case "text":
		fmt.Println(email.TextBody())
	case "html":
		fmt.Println(email.HTMLBody())
	case "raw":
		fmt.Println(email.RawBody())
	case "uuid":
		fmt.Println(email.UUID)
	default:
		return fmt.Errorf("unknown format %q (use json, text, html, raw, uuid)", format)
	}
	return nil
}
