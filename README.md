# mailexam-cli

Command-line tool for [Mailexam](https://mailexam.io/) REST API — projects, inboxes, emails, and CI/CD assertions.

## Install

```bash
go install github.com/mailexam/mailexam-cli/cmd/mailexam@latest
```

Or build from source:

```bash
git clone git@github.com:mailexam/mailexam-cli.git
cd mailexam-cli
go build -o mailexam ./cmd/mailexam
```

## Configuration

| Variable | Description |
|----------|-------------|
| `MAILEXAM_API_TOKEN` | API token from [dashboard](https://mailexam.io/login) |
| `MAILEXAM_API_BASE` | API base URL (default: `https://mailexam.io/api/v1`) |
| `MAILEXAM_PROJECT_UUID` | Default project UUID |
| `MAILEXAM_INBOX_UUID` | Default inbox UUID |

All values can be overridden with flags: `--token`, `--base`, `--project`, `--inbox`.

Regional endpoints:

- `https://mailexam.io/api/v1` (default)
- `https://mailexam.cn/api/v1`
- `https://mailexam.io/api/v1`

## Usage

```bash
export MAILEXAM_API_TOKEN="your_token"
export MAILEXAM_PROJECT_UUID="536a47df-5aad-44d0-8163-a39bb55abe0b"

# Projects
mailexam project list
mailexam project get 536a47df-5aad-44d0-8163-a39bb55abe0b

# Inboxes
mailexam inbox list --project $MAILEXAM_PROJECT_UUID
mailexam inbox get deab7974-a252-4412-9169-b965116b63cf

# Emails
mailexam email list
mailexam email get e2f9a506-d766-4935-be11-c413384de020 --format text
mailexam email wait --subject "CI check" --timeout 30
mailexam email assert --subject "CI check" --contains "Hello"
mailexam email delete e2f9a506-d766-4935-be11-c413384de020

# Attachments
mailexam email attachment download EMAIL_UUID --cid ATTACHMENT_CID -o file.pdf
```

## CI/CD example

```yaml
integration_test:
  stage: test
  variables:
    MAILEXAM_API_BASE: "https://mailexam.io/api/v1"
    MAILEXAM_PROJECT_UUID: "536a47df-5aad-44d0-8163-a39bb55abe0b"
  script:
    - npm run send-test-email
    - mailexam email assert --subject "CI check" --contains "Hello"
  # MAILEXAM_API_TOKEN - CI/CD secret
```

## Exit codes

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | Error or assertion failed |
| `2` | Authentication error (403) |
| `3` | Timeout waiting for email |

## Documentation

- [REST API](https://mailexam.io/api)
- [CLI guide (wiki)](https://wiki.mailexam.ru/en/cli/)
- [Knowledge base](https://wiki.mailexam.ru/en/)

## License

Apache 2.0
