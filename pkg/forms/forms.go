package forms

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// Parsed form key-value pairs and form-related errors
type Form struct {
	url.Values
	Errors errors
}

// Create new form instance
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Check for null required values and populate form errors field as apt
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// Check for numerical value
func (f *Form) Numerical(fields ...string) {
	for _, field := range fields {
		value := strings.TrimSpace(f.Get(field))
		if value == "" {
			continue
		}
		_, err := strconv.Atoi(value)
		if err != nil {
			f.Errors.Add(field, "This field should have a numerical value")
		}
	}
}

// Check for max length violations and populate form errors field as apt
func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (maximum is %d)", d))
	}
}

// Convert field value of the form yyyy-mm-dd to Time reference or report error
func (f *Form) Date(field string) *time.Time {
	date, err := time.Parse("2006-01-02", f.Get(field))
	if err != nil {
		f.Errors.Add(field, "Date not valid")
		return nil
	}

	return &date
}

// Compare start and end dates to make sure they are a valid pair. Also checks that the fields contain valid dates. If any error detected, a pair of nil values will be returned and form errors field will be populated as apt, Otherwise two Time references will be returned
func (f *Form) DatePair(startField, endField string) (*time.Time, *time.Time) {
	startDate := f.Date(startField)
	endDate := f.Date(endField)

	if (startDate == nil) || (endDate == nil) {
		return nil, nil
	}

	if endDate.Before(*startDate) {
		f.Errors.Add(startField, "Start date can't be after end date")
		f.Errors.Add(endField, "End date can't be before start date")
		return nil, nil
	}
	// Make tournament relevant until the end of the day. This method is okay because Date always returns time at 00:00:00
	endDate.Add(time.Hour * 23 + time.Minute * 59 + time.Second * 59)
	return startDate, endDate
}

// Check if any error exists in form errors field
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
