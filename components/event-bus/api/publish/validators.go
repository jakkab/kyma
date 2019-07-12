package publish

import (
	"regexp"
	"time"
)

var (
	isValidEventID = regexp.MustCompile(AllowedEventIDChars).MatchString

	// channel name components
	isValidSourceID         = regexp.MustCompile(AllowedSourceIDChars).MatchString
	isValidEventType        = regexp.MustCompile(AllowedEventTypeChars).MatchString
	isValidEventTypeVersion = regexp.MustCompile(AllowedEventTypeVersionChars).MatchString
)

//ValidatePublish validates a publish POST request
func ValidatePublish(r *Request, opts *EventOptions) *Error {
	if len(r.SourceID) == 0 {
		return ErrorResponseMissingFieldSourceID()
	}
	if len(r.EventType) == 0 {
		return ErrorResponseMissingFieldEventType()
	}
	if len(r.EventTypeVersion) == 0 {
		return ErrorResponseMissingFieldEventTypeVersion()
	}
	if len(r.EventTime) == 0 {
		return ErrorResponseMissingFieldEventTime()
	}
	if r.Data == nil {
		return ErrorResponseMissingFieldData()
	} else if d, ok := (r.Data).(string); ok && d == "" {
		return ErrorResponseMissingFieldData()
	}

	//validate the event components lengths
	if len(r.SourceID) > opts.MaxSourceIDLength {
		return errorInvalidSourceIDLength(opts.MaxSourceIDLength)
	}
	if len(r.EventType) > opts.MaxEventTypeLength {
		return errorInvalidEventTypeLength(opts.MaxEventTypeLength)
	}
	if len(r.EventTypeVersion) > opts.MaxEventTypeVersionLength {
		return errorInvalidEventTypeVersionLength(opts.MaxEventTypeVersionLength)
	}

	// validate the fully-qualified topic name components
	if !isValidSourceID(r.SourceID) {
		return ErrorResponseWrongSourceID(r.SourceIDFromHeader)
	}
	if !isValidEventType(r.EventType) {
		return ErrorResponseWrongEventType()
	}
	if !isValidEventTypeVersion(r.EventTypeVersion) {
		return ErrorResponseWrongEventTypeVersion()
	}

	if _, err := time.Parse(time.RFC3339, r.EventTime); err != nil {
		return ErrorResponseWrongEventTime()
	}
	if len(r.EventID) > 0 && !isValidEventID(r.EventID) {
		return ErrorResponseWrongEventID()
	}
	return nil
}