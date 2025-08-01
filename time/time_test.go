package time

import (
	"testing"
)

func TestCombine(t *testing.T) {
	tests := []struct {
		name      string
		date      string
		tm        string
		expected  string
		expectErr bool
	}{
		{
			name:      "Valid combination",
			date:      "2023-10-26",
			tm:        "15:30:00",
			expected:  "2023-10-26 15:30:00",
			expectErr: false,
		},
		{
			name:      "Invalid date format",
			date:      "2023/10/26",
			tm:        "15:30:00",
			expected:  "",
			expectErr: true,
		},
		{
			name:      "Invalid time format",
			date:      "2023-10-26",
			tm:        "15-30-00",
			expected:  "",
			expectErr: true,
		},
		{
			name:      "Empty date",
			date:      "",
			tm:        "15:30:00",
			expected:  "",
			expectErr: true,
		},
		{
			name:      "Empty time",
			date:      "2023-10-26",
			tm:        "",
			expected:  "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := Combine(tt.date, tt.tm)
			if (err != nil) != tt.expectErr {
				t.Errorf("Combine(%q, %q) error = %v, expectErr %v", tt.date, tt.tm, err, tt.expectErr)
				return
			}
			if !tt.expectErr && actual != tt.expected {
				t.Errorf("Combine(%q, %q) = %q, expected %q", tt.date, tt.tm, actual, tt.expected)
			}
		})
	}
}

func TestSeparate(t *testing.T) {
	tests := []struct {
		name         string
		datetime     string
		expectedDate string
		expectedTime string
		expectErr    bool
	}{
		{
			name:         "Valid datetime",
			datetime:     "2023-10-26 15:30:00",
			expectedDate: "2023-10-26",
			expectedTime: "15:30:00",
			expectErr:    false,
		},
		{
			name:         "Invalid datetime format",
			datetime:     "2023/10/26 15:30:00",
			expectedDate: "",
			expectedTime: "",
			expectErr:    true,
		},
		{
			name:         "Empty datetime",
			datetime:     "",
			expectedDate: "",
			expectedTime: "",
			expectErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualDate, actualTime, err := Separate(tt.datetime)
			if (err != nil) != tt.expectErr {
				t.Errorf("Separate(%q) error = %v, expectErr %v", tt.datetime, err, tt.expectErr)
				return
			}
			if !tt.expectErr {
				if actualDate != tt.expectedDate {
					t.Errorf("Separate(%q) date = %q, expected %q", tt.datetime, actualDate, tt.expectedDate)
				}
				if actualTime != tt.expectedTime {
					t.Errorf("Separate(%q) time = %q, expected %q", tt.datetime, actualTime, tt.expectedTime)
				}
			}
		})
	}
}
