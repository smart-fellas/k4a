package unit

import (
	"testing"

	"github.com/smart-fellas/k4a/internal/utils"
)

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name  string
		bytes int64
		want  string
	}{
		{
			name:  "bytes",
			bytes: 512,
			want:  "512 B",
		},
		{
			name:  "kilobytes",
			bytes: 1024,
			want:  "1.0 KiB",
		},
		{
			name:  "megabytes",
			bytes: 1048576,
			want:  "1.0 MiB",
		},
		{
			name:  "gigabytes",
			bytes: 1073741824,
			want:  "1.0 GiB",
		},
		{
			name:  "complex size",
			bytes: 1536,
			want:  "1.5 KiB",
		},
		{
			name:  "zero bytes",
			bytes: 0,
			want:  "0 B",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.FormatBytes(tt.bytes)
			if got != tt.want {
				t.Errorf("FormatBytes(%d) = %v, want %v", tt.bytes, got, tt.want)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name string
		ms   int64
		want string
	}{
		{
			name: "milliseconds",
			ms:   500,
			want: "500ms",
		},
		{
			name: "seconds",
			ms:   5000,
			want: "5.0s",
		},
		{
			name: "minutes",
			ms:   120000,
			want: "2.0m",
		},
		{
			name: "hours",
			ms:   7200000,
			want: "2.0h",
		},
		{
			name: "days",
			ms:   172800000,
			want: "2.0d",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.FormatDuration(tt.ms)
			if got != tt.want {
				t.Errorf("FormatDuration(%d) = %v, want %v", tt.ms, got, tt.want)
			}
		})
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		want   string
	}{
		{
			name:   "no truncation needed",
			input:  "short",
			maxLen: 10,
			want:   "short",
		},
		{
			name:   "exact length",
			input:  "exact",
			maxLen: 5,
			want:   "exact",
		},
		{
			name:   "needs truncation",
			input:  "this is a very long string",
			maxLen: 10,
			want:   "this is...",
		},
		{
			name:   "very short max length",
			input:  "hello",
			maxLen: 3,
			want:   "hel",
		},
		{
			name:   "empty string",
			input:  "",
			maxLen: 5,
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.TruncateString(tt.input, tt.maxLen)
			if got != tt.want {
				t.Errorf("TruncateString(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
			}
		})
	}
}

func TestPadRight(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		length int
		want   string
	}{
		{
			name:   "pad needed",
			input:  "test",
			length: 10,
			want:   "test      ",
		},
		{
			name:   "no padding needed",
			input:  "exact",
			length: 5,
			want:   "exact",
		},
		{
			name:   "longer than length",
			input:  "toolong",
			length: 5,
			want:   "toolong",
		},
		{
			name:   "empty string",
			input:  "",
			length: 5,
			want:   "     ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.PadRight(tt.input, tt.length)
			if got != tt.want {
				t.Errorf("PadRight(%q, %d) = %q, want %q", tt.input, tt.length, got, tt.want)
			}
		})
	}
}
