package api

import "testing"

func TestEmailFilterMatch(t *testing.T) {
	email := EmailShort{
		Subject: "CI check",
		To:      "user@example.com",
		From:    "noreply@example.com",
	}

	tests := []struct {
		name   string
		filter EmailFilter
		want   bool
	}{
		{
			name:   "subject substring",
			filter: EmailFilter{Subject: "CI"},
			want:   true,
		},
		{
			name:   "subject exact match",
			filter: EmailFilter{Subject: "CI check", SubjectExact: true},
			want:   true,
		},
		{
			name:   "subject exact fail",
			filter: EmailFilter{Subject: "CI", SubjectExact: true},
			want:   false,
		},
		{
			name:   "to filter",
			filter: EmailFilter{To: "user@"},
			want:   true,
		},
		{
			name:   "empty filter",
			filter: EmailFilter{},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.filter.Match(email); got != tt.want {
				t.Fatalf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}
