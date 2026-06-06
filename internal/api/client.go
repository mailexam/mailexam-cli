package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultTimeout = 30 * time.Second

type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("api error %d: %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("api error %d", e.StatusCode)
}

func NewClient(baseURL, token string) *Client {
	baseURL = strings.TrimRight(baseURL, "/")
	return &Client{
		BaseURL: baseURL,
		Token:   token,
		HTTPClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

func (c *Client) Get(path string, query url.Values) ([]byte, error) {
	if query == nil {
		query = url.Values{}
	}
	query.Set("token", c.Token)

	u, err := url.Parse(c.BaseURL + path)
	if err != nil {
		return nil, err
	}
	u.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    strings.TrimSpace(string(body)),
		}
	}

	return body, nil
}

func (c *Client) Delete(path string) error {
	query := url.Values{"token": {c.Token}}
	u, err := url.Parse(c.BaseURL + path)
	if err != nil {
		return err
	}
	u.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodDelete, u.String(), nil)
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    strings.TrimSpace(string(body)),
		}
	}

	return nil
}

func decode[T any](data []byte, out *T) error {
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, out)
}

func (c *Client) ListProjects() ([]Project, error) {
	data, err := c.Get("/project", nil)
	if err != nil {
		return nil, err
	}
	var projects []Project
	return projects, decode(data, &projects)
}

func (c *Client) GetProject(uuid string) (*Project, error) {
	data, err := c.Get("/project/"+uuid, nil)
	if err != nil {
		return nil, err
	}
	var project Project
	if err := decode(data, &project); err != nil {
		return nil, err
	}
	return &project, nil
}

func (c *Client) ListInboxes(projectUUID string) ([]Inbox, error) {
	query := url.Values{"project": {projectUUID}}
	data, err := c.Get("/inbox", query)
	if err != nil {
		return nil, err
	}
	var inboxes []Inbox
	return inboxes, decode(data, &inboxes)
}

func (c *Client) GetInbox(uuid string) (*Inbox, error) {
	data, err := c.Get("/inbox/"+uuid, nil)
	if err != nil {
		return nil, err
	}
	var inbox Inbox
	if err := decode(data, &inbox); err != nil {
		return nil, err
	}
	return &inbox, nil
}

func (c *Client) ListEmails(inboxUUID string) ([]EmailShort, error) {
	query := url.Values{"inbox": {inboxUUID}}
	data, err := c.Get("/email", query)
	if err != nil {
		return nil, err
	}
	var emails []EmailShort
	return emails, decode(data, &emails)
}

func (c *Client) GetEmail(uuid string) (*EmailDetailed, error) {
	data, err := c.Get("/email/"+uuid, nil)
	if err != nil {
		return nil, err
	}
	var email EmailDetailed
	if err := decode(data, &email); err != nil {
		return nil, err
	}
	return &email, nil
}

func (c *Client) DeleteEmail(uuid string) error {
	return c.Delete("/email/" + uuid + "/delete")
}

func (c *Client) DownloadAttachment(emailUUID, cid string) ([]byte, error) {
	query := url.Values{"cid": {cid}}
	return c.Get("/email/"+emailUUID+"/attachment", query)
}

func (c *Client) ResolveInbox(projectUUID, inboxUUID string) (string, error) {
	if inboxUUID != "" {
		return inboxUUID, nil
	}
	if projectUUID == "" {
		return "", fmt.Errorf("project or inbox is required")
	}

	inboxes, err := c.ListInboxes(projectUUID)
	if err != nil {
		return "", err
	}
	if len(inboxes) == 0 {
		return "", fmt.Errorf("no inboxes found for project %s", projectUUID)
	}

	for _, inbox := range inboxes {
		if inbox.IsDefault {
			return inbox.UUID, nil
		}
	}

	return inboxes[0].UUID, nil
}

type EmailFilter struct {
	Subject      string
	To           string
	From         string
	SubjectExact bool
}

func (f EmailFilter) Match(item EmailShort) bool {
	if f.Subject != "" {
		if f.SubjectExact {
			if item.Subject != f.Subject {
				return false
			}
		} else if !strings.Contains(item.Subject, f.Subject) {
			return false
		}
	}
	if f.To != "" && !strings.Contains(item.To, f.To) {
		return false
	}
	if f.From != "" && !strings.Contains(item.From, f.From) {
		return false
	}
	return f.Subject != "" || f.To != "" || f.From != ""
}

type WaitOptions struct {
	Filter   EmailFilter
	Project  string
	Inbox    string
	Timeout  time.Duration
	Interval time.Duration
}

func (c *Client) WaitForEmail(opts WaitOptions) (*EmailDetailed, error) {
	inboxUUID, err := c.ResolveInbox(opts.Project, opts.Inbox)
	if err != nil {
		return nil, err
	}

	if opts.Timeout <= 0 {
		opts.Timeout = 30 * time.Second
	}
	if opts.Interval <= 0 {
		opts.Interval = 2 * time.Second
	}

	deadline := time.Now().Add(opts.Timeout)
	for {
		emails, err := c.ListEmails(inboxUUID)
		if err != nil {
			return nil, err
		}

		for _, item := range emails {
			if opts.Filter.Match(item) {
				return c.GetEmail(item.UUID)
			}
		}

		if time.Now().After(deadline) {
			return nil, fmt.Errorf("email not found within %s", opts.Timeout)
		}

		time.Sleep(opts.Interval)
	}
}
