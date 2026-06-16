package models

import (
	"testing"
	"time"
)

func TestCalculateAge(t *testing.T) {
	tests := []struct {
		name     string
		dob      time.Time
		wantAge  int
		wantSign int // 1 = positive, 0 = zero, -1 = could be negative (future)
	}{
		{
			name:    "born 30 years ago",
			dob:     time.Now().AddDate(-30, 0, 0),
			wantAge: 30,
		},
		{
			name:    "born 1 year ago",
			dob:     time.Now().AddDate(-1, 0, 0),
			wantAge: 1,
		},
		{
			name:    "born today",
			dob:     time.Now(),
			wantAge: 0,
		},
		{
			name:    "born 25 years and 6 months ago",
			dob:     time.Now().AddDate(-25, -6, 0),
			wantAge: 25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateAge(tt.dob)
			if got != tt.wantAge {
				t.Errorf("CalculateAge() = %d, want %d (dob: %s)", got, tt.wantAge, tt.dob.Format("2006-01-02"))
			}
		})
	}
}

func TestCalculateAge_BirthdayNotYetThisYear(t *testing.T) {
	// Set DOB to 30 years ago, but 6 months in the future from today's month.
	now := time.Now()
	dob := time.Date(now.Year()-30, now.Month()+6, now.Day(), 0, 0, 0, 0, time.UTC)

	// If adding 6 months goes past December, the year rolls over,
	// so adjust for that edge case.
	if dob.After(now) || dob.Equal(now) {
		// Birthday hasn't happened yet this year.
		age := CalculateAge(dob)
		if age != 29 {
			t.Errorf("CalculateAge() = %d, want 29 (birthday hasn't passed yet)", age)
		}
	}
}

func TestCalculateAge_LeapYearBirthday(t *testing.T) {
	// Feb 29 birthday — non-leap-year handling.
	dob := time.Date(1996, time.February, 29, 0, 0, 0, 0, time.UTC)
	age := CalculateAge(dob)

	now := time.Now()
	expectedAge := now.Year() - 1996
	if now.YearDay() < dob.YearDay() {
		expectedAge--
	}

	if age != expectedAge {
		t.Errorf("CalculateAge() = %d, want %d (leap year birthday)", age, expectedAge)
	}
}
