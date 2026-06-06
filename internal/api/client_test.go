package api

import "testing"

func TestEmailFilterMatch(t *testing.T) {
	email := EmailShort{
		Subject: "Проверка CI",
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
			filter: EmailFilter{Subject: "Проверка"},
			want:   true,
		},
		{
			name:   "subject exact miss",
			filter: EmailFilter{Subject: "Проверка CI", SubjectExact: true},
			want:   true,
		},
		{
			name:   "subject exact fail",
			filter: EmailFilter{Subject: "Проверка", SubjectExact: true},
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
