package config

import (
	"fmt"
	"os"
)

const (
	EnvToken       = "MAILEXAM_API_TOKEN"
	EnvBase        = "MAILEXAM_API_BASE"
	EnvProjectUUID = "MAILEXAM_PROJECT_UUID"
	EnvInboxUUID   = "MAILEXAM_INBOX_UUID"

	DefaultBase = "https://mailexam.ru/api/v1"
)

type Config struct {
	BaseURL     string
	Token       string
	ProjectUUID string
	InboxUUID   string
}

func Load(baseURL, token, projectUUID, inboxUUID string) (Config, error) {
	cfg := Config{
		BaseURL:     firstNonEmpty(baseURL, os.Getenv(EnvBase), DefaultBase),
		Token:       firstNonEmpty(token, os.Getenv(EnvToken)),
		ProjectUUID: firstNonEmpty(projectUUID, os.Getenv(EnvProjectUUID)),
		InboxUUID:   firstNonEmpty(inboxUUID, os.Getenv(EnvInboxUUID)),
	}

	if cfg.Token == "" {
		return cfg, fmt.Errorf("api token is required (flag --token or env %s)", EnvToken)
	}

	return cfg, nil
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
