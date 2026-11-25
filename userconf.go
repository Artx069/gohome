package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	defaultWorkDays   = 5
	defaultTargetTime = 8 * time.Hour
	// GetLeaveTime rejects targets above the legal daily maximum.
	maxDailyWorkTime = 10 * time.Hour
)

type UserConfig struct {
	TargetTimeStr string  `json:"TargetTime"`
	WeeklyHours   float64 `json:"WeeklyHours,omitempty"`
	WorkDays      int     `json:"WorkDays,omitempty"`
}

// DailyTarget derives the daily target work time from the weekly hours, rounded to full minutes.
func (u UserConfig) DailyTarget() (time.Duration, bool) {
	if u.WeeklyHours <= 0 {
		return 0, false
	}
	weekly := time.Duration(u.WeeklyHours * float64(time.Hour))
	return (weekly / time.Duration(u.EffectiveWorkDays())).Round(time.Minute), true
}

func (u UserConfig) EffectiveWorkDays() int {
	if u.WorkDays <= 0 {
		return defaultWorkDays
	}
	return u.WorkDays
}

// ApplySetting updates the config from a "key=value" pair. An empty value resets the key.
func (u *UserConfig) ApplySetting(setting string) error {
	key, value, ok := strings.Cut(setting, "=")
	if !ok {
		return fmt.Errorf("invalid setting %q, expected key=value", setting)
	}
	key = strings.ToLower(strings.TrimSpace(key))
	value = strings.TrimSpace(value)

	switch key {
	case "workhours":
		if value == "" {
			u.WeeklyHours = 0
			return nil
		}
		hours, err := strconv.ParseFloat(strings.Replace(value, ",", ".", 1), 64)
		if err != nil || hours <= 0 {
			return fmt.Errorf("invalid workhours %q, expected a positive number like 36 or 38.5", value)
		}
		u.WeeklyHours = hours
	case "workdays":
		if value == "" {
			u.WorkDays = 0
			return nil
		}
		days, err := strconv.Atoi(value)
		if err != nil || days < 1 || days > 7 {
			return fmt.Errorf("invalid workdays %q, expected a number between 1 and 7", value)
		}
		u.WorkDays = days
	default:
		return fmt.Errorf("unknown setting %q, supported: workhours, workdays", key)
	}
	return nil
}

// Validate ensures the derived daily target can be used for leave time computation.
func (u UserConfig) Validate() error {
	daily, ok := u.DailyTarget()
	if !ok {
		return nil
	}
	if daily <= 0 || daily > maxDailyWorkTime {
		return fmt.Errorf("daily target of %s exceeds the allowed maximum of %s", formatDurationMinutes(daily), formatDurationMinutes(maxDailyWorkTime))
	}
	return nil
}

// resolveTargetTime picks the daily target: command line flag, then weekly hours, then the legacy target time.
func resolveTargetTime(flagValue string, usrConf UserConfig) (time.Duration, error) {
	if flagValue != "" {
		return parseClockDuration(flagValue)
	}
	if daily, ok := usrConf.DailyTarget(); ok {
		return daily, nil
	}
	if usrConf.TargetTimeStr != "" {
		return parseClockDuration(usrConf.TargetTimeStr)
	}
	return defaultTargetTime, nil
}

func parseClockDuration(value string) (time.Duration, error) {
	t, err := time.Parse("15:04", value)
	if err != nil {
		return 0, fmt.Errorf("failed to parse target time: %w", err)
	}
	return time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute, nil
}

func ReadUserConfig() (UserConfig, error) {
	configDir := getConfigDir()
	matrixOutputFileDir = configDir
	userConfFile := filepath.Join(configDir, "userconfig.json")

	data, err := os.ReadFile(userConfFile)
	if err != nil {
		if os.IsNotExist(err) {
			return UserConfig{}, nil
		}
		return UserConfig{}, err
	}

	var usrConf UserConfig
	if err := json.Unmarshal(data, &usrConf); err != nil {
		return UserConfig{}, err
	}

	return usrConf, nil
}

func WriteUserConfig(usrConf UserConfig) error {
	configDir := getConfigDir()
	if err := os.MkdirAll(configDir, 0o700); err != nil {
		return err
	}
	userConfFile := filepath.Join(configDir, "userconfig.json")

	data, err := json.Marshal(&usrConf)
	if err != nil {
		return err
	}

	return os.WriteFile(userConfFile, data, os.ModePerm)
}
