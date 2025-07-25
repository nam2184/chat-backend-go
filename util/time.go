package util

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type CustomTime time.Time

func (ct *CustomTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	parse, err := time.Parse("2006-01-02T15:04:05", v)
	if err != nil {
		return err
  }
	*ct = CustomTime(parse)
	return nil
}

type Duration time.Duration

func (d *Duration) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	var durationStr string
	if err := dec.DecodeElement(&durationStr, &start); err != nil {
		return err
	}
  duration, err := parseISODuration(durationStr) 
	if err != nil {
		return err
	}
	*d = Duration(duration)
	return nil
}

// Function to parse the ISO 8601 duration string
func parseISODuration(duration string) (time.Duration, error) {
    // Regular expression to match the ISO 8601 duration format
    re := regexp.MustCompile(`^PT(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?$`)
    matches := re.FindStringSubmatch(duration)

    if matches == nil {
        return 0, fmt.Errorf("invalid duration format")
    }
    var totalDuration time.Duration

    if len(matches) > 1 && matches[1] != "" {
        hours := matches[1]
        totalDuration += time.Duration(parseInt(hours)) * time.Hour
    }

    if len(matches) > 2 && matches[2] != "" {
        minutes := matches[2]
        totalDuration += time.Duration(parseInt(minutes)) * time.Minute
    }

    if len(matches) > 3 && matches[3] != "" {
        seconds := matches[3]
        totalDuration += time.Duration(parseInt(seconds)) * time.Second
    }
    return totalDuration, nil
}

// Helper function to convert string to int
func parseInt(s string) int {
    val, err := strconv.Atoi(s)
    if err != nil {
        return 0
    }
    return val
}

