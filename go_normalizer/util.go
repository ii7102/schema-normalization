package gonormalizer

import (
	"time"

	"github.com/ii7102/diploma-schema-normalization/rules"
)

func baseLayoutFormat(bt rules.BaseType) string {
	switch bt {
	case rules.Timestamp:
		return time.TimeOnly
	case rules.Date:
		return time.DateOnly
	case rules.DateTime:
		return time.DateTime
	default:
		return ""
	}
}

func layoutFormats(bt rules.BaseType) []string {
	switch bt {
	case rules.Timestamp:
		return timestampFormats()
	case rules.Date:
		return dateFormats()
	case rules.DateTime:
		return dateTimeFormats()
	default:
		return nil
	}
}

func timestampFormats() []string {
	return []string{
		time.TimeOnly, time.Kitchen, time.Stamp, time.StampMilli, time.StampMicro, time.StampNano,
	}
}

func dateFormats() []string {
	return []string{
		time.DateOnly,
	}
}

func dateTimeFormats() []string {
	return []string{
		time.RFC3339, time.DateTime, time.RFC3339Nano, time.RFC1123, time.RFC1123Z,
		time.RFC850, time.RFC822, time.RFC822Z, time.ANSIC, time.UnixDate, time.RubyDate,
	}
}
