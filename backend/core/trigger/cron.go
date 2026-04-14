package trigger

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// CronExpression represents a parsed cron expression (minute hour dom month dow).
type CronExpression struct {
	Minutes    []int
	Hours      []int
	DaysOfMonth []int
	Months     []int
	DaysOfWeek []int
}

// ParseCron parses a standard 5-field cron expression.
// Format: minute hour day-of-month month day-of-week
// Supports: *, specific values, ranges (1-5), steps (*/5), lists (1,3,5)
func ParseCron(expr string) (*CronExpression, error) {
	fields := strings.Fields(strings.TrimSpace(expr))
	if len(fields) != 5 {
		return nil, fmt.Errorf("invalid cron expression: expected 5 fields, got %d", len(fields))
	}

	minutes, err := parseField(fields[0], 0, 59)
	if err != nil {
		return nil, fmt.Errorf("invalid minute field: %w", err)
	}
	hours, err := parseField(fields[1], 0, 23)
	if err != nil {
		return nil, fmt.Errorf("invalid hour field: %w", err)
	}
	dom, err := parseField(fields[2], 1, 31)
	if err != nil {
		return nil, fmt.Errorf("invalid day-of-month field: %w", err)
	}
	months, err := parseField(fields[3], 1, 12)
	if err != nil {
		return nil, fmt.Errorf("invalid month field: %w", err)
	}
	dow, err := parseField(fields[4], 0, 6)
	if err != nil {
		return nil, fmt.Errorf("invalid day-of-week field: %w", err)
	}

	return &CronExpression{
		Minutes:     minutes,
		Hours:       hours,
		DaysOfMonth: dom,
		Months:      months,
		DaysOfWeek:  dow,
	}, nil
}

// NextAfter returns the next time the cron expression matches after the given time.
func (c *CronExpression) NextAfter(after time.Time) time.Time {
	// Start from the next minute
	t := after.Truncate(time.Minute).Add(time.Minute)

	// Search up to 4 years ahead to find a match
	maxTime := after.Add(4 * 365 * 24 * time.Hour)

	for t.Before(maxTime) {
		if !contains(c.Months, int(t.Month())) {
			// Skip to next month
			t = time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
			continue
		}
		if !contains(c.DaysOfMonth, t.Day()) || !contains(c.DaysOfWeek, int(t.Weekday())) {
			// Skip to next day
			t = time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, 0, t.Location())
			continue
		}
		if !contains(c.Hours, t.Hour()) {
			// Skip to next hour
			t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour()+1, 0, 0, 0, t.Location())
			continue
		}
		if !contains(c.Minutes, t.Minute()) {
			t = t.Add(time.Minute)
			continue
		}
		return t
	}

	// Fallback: return 24h from now if no match found
	return after.Add(24 * time.Hour)
}

func parseField(field string, min, max int) ([]int, error) {
	if field == "*" {
		return makeRange(min, max), nil
	}

	var result []int
	parts := strings.Split(field, ",")
	for _, part := range parts {
		// Check for step value
		if strings.Contains(part, "/") {
			stepParts := strings.SplitN(part, "/", 2)
			step, err := strconv.Atoi(stepParts[1])
			if err != nil || step <= 0 {
				return nil, fmt.Errorf("invalid step: %s", part)
			}
			start := min
			if stepParts[0] != "*" {
				start, err = strconv.Atoi(stepParts[0])
				if err != nil {
					return nil, fmt.Errorf("invalid step base: %s", part)
				}
			}
			for i := start; i <= max; i += step {
				result = append(result, i)
			}
			continue
		}

		// Check for range
		if strings.Contains(part, "-") {
			rangeParts := strings.SplitN(part, "-", 2)
			rangeStart, err := strconv.Atoi(rangeParts[0])
			if err != nil {
				return nil, fmt.Errorf("invalid range start: %s", part)
			}
			rangeEnd, err := strconv.Atoi(rangeParts[1])
			if err != nil {
				return nil, fmt.Errorf("invalid range end: %s", part)
			}
			for i := rangeStart; i <= rangeEnd; i++ {
				result = append(result, i)
			}
			continue
		}

		// Single value
		val, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid value: %s", part)
		}
		if val < min || val > max {
			return nil, fmt.Errorf("value %d out of range [%d, %d]", val, min, max)
		}
		result = append(result, val)
	}

	return result, nil
}

func makeRange(min, max int) []int {
	r := make([]int, 0, max-min+1)
	for i := min; i <= max; i++ {
		r = append(r, i)
	}
	return r
}

func contains(slice []int, val int) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

// VisualConfigToCron converts a visual schedule config to a cron expression.
func VisualConfigToCron(frequency string, onMinute int, timeStr string, weekdays []string, monthlyDays []int) (string, error) {
	hour, minute := 0, 0
	if timeStr != "" {
		h, m, err := parseTimeString(timeStr)
		if err != nil {
			return "", err
		}
		hour, minute = h, m
	}

	switch frequency {
	case "hourly":
		return fmt.Sprintf("%d * * * *", onMinute), nil
	case "daily":
		return fmt.Sprintf("%d %d * * *", minute, hour), nil
	case "weekly":
		dowStr := convertWeekdays(weekdays)
		return fmt.Sprintf("%d %d * * %s", minute, hour, dowStr), nil
	case "monthly":
		dayStr := convertMonthlyDays(monthlyDays)
		return fmt.Sprintf("%d %d %s * *", minute, hour, dayStr), nil
	default:
		return "", fmt.Errorf("unsupported frequency: %s", frequency)
	}
}

func parseTimeString(timeStr string) (int, int, error) {
	// Parse "2:30 PM" or "12:00 AM" format
	timeStr = strings.TrimSpace(timeStr)
	upper := strings.ToUpper(timeStr)
	isPM := strings.HasSuffix(upper, "PM")
	isAM := strings.HasSuffix(upper, "AM")

	clean := strings.TrimSuffix(strings.TrimSuffix(upper, "PM"), "AM")
	clean = strings.TrimSpace(clean)

	parts := strings.SplitN(clean, ":", 2)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid time format: %s", timeStr)
	}

	hour, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, err
	}
	minute, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, err
	}

	if isPM && hour != 12 {
		hour += 12
	} else if isAM && hour == 12 {
		hour = 0
	}

	return hour, minute, nil
}

func convertWeekdays(weekdays []string) string {
	dayMap := map[string]string{
		"sun": "0", "mon": "1", "tue": "2", "wed": "3",
		"thu": "4", "fri": "5", "sat": "6",
	}
	var nums []string
	for _, d := range weekdays {
		if num, ok := dayMap[strings.ToLower(d)]; ok {
			nums = append(nums, num)
		}
	}
	if len(nums) == 0 {
		return "*"
	}
	return strings.Join(nums, ",")
}

func convertMonthlyDays(days []int) string {
	if len(days) == 0 {
		return "*"
	}
	var strs []string
	for _, d := range days {
		strs = append(strs, strconv.Itoa(d))
	}
	return strings.Join(strs, ",")
}
