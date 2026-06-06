package api

import "time"

type User struct {
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Project struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Description *string `json:"description"`
	User        User   `json:"user"`
}

type Inbox struct {
	UUID         string  `json:"uuid"`
	ProjectUUID  string  `json:"project_uuid"`
	Name         string  `json:"name"`
	IsDefault    bool    `json:"is_default"`
	HasAutorelay bool    `json:"has_autorelay"`
	RelayHost    *string `json:"relay_host"`
	RelayPort    *string `json:"relay_port"`
	RelayUser    *string `json:"relay_user"`
	RelayPass    *string `json:"relay_pass"`
}

type EmailShort struct {
	UUID        string `json:"uuid"`
	From        string `json:"from"`
	To          string `json:"to"`
	CC          *string `json:"cc"`
	Subject     string `json:"subject"`
	HumanSize   string `json:"human_size"`
	Date        string `json:"date"`
	Attachments bool   `json:"attachments"`
	Unread      bool   `json:"unread"`
	Like        bool   `json:"like"`
	Dislike     bool   `json:"dislike"`
}

type EmailBody struct {
	HTML        *string           `json:"html"`
	Text        *string           `json:"text"`
	Raw         *string           `json:"raw"`
	Attachments []EmailAttachment `json:"attachments"`
}

type EmailAttachment struct {
	CID      string `json:"cid"`
	Filename string `json:"filename"`
	MimeType string `json:"mime_type"`
	Size     int    `json:"size"`
}

type EmailDetailed struct {
	EmailShort
	Size    int                    `json:"size"`
	Body    EmailBody              `json:"body"`
	Headers []map[string]any       `json:"headers"`
	Relay   []map[string]any       `json:"relay"`
}

func (e EmailShort) ParsedDate() time.Time {
	t, err := time.Parse(time.RFC3339Nano, e.Date)
	if err != nil {
		t, _ = time.Parse(time.RFC3339, e.Date)
	}
	return t
}

func (e EmailDetailed) TextBody() string {
	if e.Body.Text != nil {
		return *e.Body.Text
	}
	return ""
}

func (e EmailDetailed) HTMLBody() string {
	if e.Body.HTML != nil {
		return *e.Body.HTML
	}
	return ""
}

func (e EmailDetailed) RawBody() string {
	if e.Body.Raw != nil {
		return *e.Body.Raw
	}
	return ""
}
