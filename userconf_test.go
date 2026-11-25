package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDailyTarget(t *testing.T) {
	testCases := []struct {
		name     string
		conf     UserConfig
		expected time.Duration
		ok       bool
	}{
		{"not set", UserConfig{}, 0, false},
		{"36h default 5 days", UserConfig{WeeklyHours: 36}, dur(7, 12), true},
		{"40h default 5 days", UserConfig{WeeklyHours: 40}, dur(8, 0), true},
		{"38.5h default 5 days", UserConfig{WeeklyHours: 38.5}, dur(7, 42), true},
		{"30h 4 days", UserConfig{WeeklyHours: 30, WorkDays: 4}, dur(7, 30), true},
		{"rounded to full minutes", UserConfig{WeeklyHours: 37, WorkDays: 6}, dur(6, 10), true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			daily, ok := tc.conf.DailyTarget()
			assert.Equal(t, tc.ok, ok)
			assert.Equal(t, tc.expected, daily)
		})
	}
}

func TestApplySetting(t *testing.T) {
	t.Run("workhours and workdays", func(t *testing.T) {
		var conf UserConfig
		require.NoError(t, conf.ApplySetting("workhours=36"))
		require.NoError(t, conf.ApplySetting("workdays=4"))
		assert.Equal(t, UserConfig{WeeklyHours: 36, WorkDays: 4}, conf)
	})

	t.Run("decimal comma and case insensitive key", func(t *testing.T) {
		var conf UserConfig
		require.NoError(t, conf.ApplySetting(" WorkHours = 38,5 "))
		assert.Equal(t, 38.5, conf.WeeklyHours)
	})

	t.Run("empty value resets", func(t *testing.T) {
		conf := UserConfig{WeeklyHours: 36, WorkDays: 4}
		require.NoError(t, conf.ApplySetting("workhours="))
		require.NoError(t, conf.ApplySetting("workdays="))
		assert.Equal(t, UserConfig{}, conf)
	})

	invalid := []string{"workhours", "workhours=abc", "workhours=0", "workhours=-5", "workdays=0", "workdays=8", "workdays=2.5", "unknown=1"}
	for _, setting := range invalid {
		t.Run("invalid "+setting, func(t *testing.T) {
			conf := UserConfig{WeeklyHours: 36}
			assert.Error(t, conf.ApplySetting(setting))
			assert.Equal(t, UserConfig{WeeklyHours: 36}, conf, "config must stay unchanged on error")
		})
	}
}

func TestValidate(t *testing.T) {
	assert.NoError(t, UserConfig{}.Validate())
	assert.NoError(t, UserConfig{WeeklyHours: 50}.Validate())
	assert.Error(t, UserConfig{WeeklyHours: 51}.Validate())
	assert.Error(t, UserConfig{WeeklyHours: 36, WorkDays: 3}.Validate())
}

func TestResolveTargetTime(t *testing.T) {
	testCases := []struct {
		name     string
		flag     string
		conf     UserConfig
		expected time.Duration
	}{
		{"default", "", UserConfig{}, dur(8, 0)},
		{"legacy target time", "", UserConfig{TargetTimeStr: "07:30"}, dur(7, 30)},
		{"workhours beat legacy target time", "", UserConfig{TargetTimeStr: "07:30", WeeklyHours: 36}, dur(7, 12)},
		{"flag beats everything", "06:45", UserConfig{TargetTimeStr: "07:30", WeeklyHours: 36}, dur(6, 45)},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			target, err := resolveTargetTime(tc.flag, tc.conf)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, target)
		})
	}

	t.Run("invalid flag", func(t *testing.T) {
		_, err := resolveTargetTime("abc", UserConfig{})
		assert.Error(t, err)
	})
}
