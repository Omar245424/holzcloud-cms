package dsr

// Correction: the operator corrects what a website holds about one person
// (Phase 33.2, DSR-02, research section "Correction per record kind").
//
// What may be corrected depends on what the record is:
//
//   - An order's contact fields, its e-mail address and phone number, are
//     corrected freely. They are how the shop reaches the customer, and the
//     invoice does not print them.
//
// The record is found again by a search of the person's address on the
// website the caller names, and the source and key the form posts have to be
// one of the records that search returns. Whatever a form posts, it can only
// reach this person's records on this website. The column a correction
// writes comes from a fixed statement per field, never from the posted field
// name.

import (
	"context"
	"errors"
	"strings"
	"time"
)

// CorrectableField is one field of a record the operator may correct.
type CorrectableField struct {
	// Name is the field's stable name, as in Field.Name.
	Name string
	// Label is what the screen shows: a catalogue key for the core sources,
	// the plugin's own label for a plugin.
	Label string
	// Value is the field as it is stored now.
	Value string
	// Invoice is true for a field the invoice prints. Correcting it keeps
	// the old value on record.
	Invoice bool
}

// Change is one correction as it was made.
type Change struct {
	Source, Key, Field string
	// Old and New are the value before and after.
	Old, New string
	// Invoice is true when the field is one the invoice prints.
	Invoice bool
}

var (
	// ErrNotCorrectable means the field, or the whole record, is not one
	// that is corrected here.
	ErrNotCorrectable = errors.New("dsr: this field is not corrected here")
	// ErrErased means the record has been reduced by an erasure and is not
	// changed any more.
	ErrErased = errors.New("dsr: the record was erased")
	// ErrEmptyValue means a field that must hold something was left empty.
	ErrEmptyValue = errors.New("dsr: the value is empty")
)

// correctField is how one field of a core source is corrected.
type correctField struct {
	// stmt is the whole UPDATE: $1 website, $2 key, $3 value, $4 the time.
	stmt string
	// max is the longest value accepted, in characters.
	max int
	// required refuses an empty value.
	required bool
	// invoice marks a field the invoice prints.
	invoice bool
}

// orderFields are the corrections of an order, in the order the screen shows
// them. An erased order is never matched (erased_at IS NULL): the erasure is
// one-way.
var orderFields = []string{"email", "phone"}

var orderCorrections = map[string]correctField{
	"email": {stmt: `UPDATE orders SET email = $3, updated_at = $4
		WHERE website_id = $1 AND number = $2 AND erased_at IS NULL`, max: 254, required: true},
	"phone": {stmt: `UPDATE orders SET phone = $3, updated_at = $4
		WHERE website_id = $1 AND number = $2 AND erased_at IS NULL`, max: 50},
}

// fieldValue is the value of the named field of a record, "" when the record
// has none.
func fieldValue(rec Record, name string) (string, string) {
	for _, f := range rec.Fields {
		if f.Name == name {
			return f.Label, f.Value
		}
	}
	return fieldLabels[name], ""
}

// Correctable is what the operator may correct of a record, in the order the
// screen shows it. A record that is not corrected here has none.
func Correctable(rec Record) []CorrectableField {
	if rec.Erased {
		return nil
	}
	var names []string
	var table map[string]correctField
	switch rec.Source {
	case "orders":
		names, table = orderFields, orderCorrections
	default:
		return nil
	}
	out := make([]CorrectableField, 0, len(names))
	for _, name := range names {
		label, value := fieldValue(rec, name)
		out = append(out, CorrectableField{Name: name, Label: label, Value: value, Invoice: table[name].invoice})
	}
	return out
}

// Correct sets one field of one of the person's records on the website.
//
// email is the address the record was found by; the record is searched for
// again under it, and source and key must name one of the records found
// (ErrNotFound otherwise). An erased record is ErrErased, a field that is not
// corrected here ErrNotCorrectable. The value is trimmed and checked as the
// form it came from checks it: an address must be one bare address
// (ErrBadEmail), a required field must hold something (ErrEmptyValue), and no
// value may be longer than the field allows (ErrTooLong).
func (s *Service) Correct(ctx context.Context, websiteID int64, email, source, key, field, value string) (*Change, error) {
	res, err := s.Find(ctx, websiteID, email)
	if err != nil {
		return nil, err
	}
	var rec *Record
	for i := range res.Records {
		if res.Records[i].Source == source && res.Records[i].Key == key {
			rec = &res.Records[i]
			break
		}
	}
	if rec == nil {
		// A store that could not be searched may hold the record: its
		// error says more than "not found".
		for _, f := range res.Failed {
			if f.Source == source {
				return nil, f.Err
			}
		}
		return nil, ErrNotFound
	}
	if rec.Erased {
		return nil, ErrErased
	}
	var allowed *CorrectableField
	for _, f := range Correctable(*rec) {
		if f.Name == field {
			allowed = &f
			break
		}
	}
	if allowed == nil {
		return nil, ErrNotCorrectable
	}
	value = strings.TrimSpace(value)
	ch := &Change{Source: source, Key: key, Field: field, Old: allowed.Value, New: value, Invoice: allowed.Invoice}

	switch source {
	case "orders":
		err = s.correctRow(ctx, websiteID, key, value, orderCorrections[field], field)
	default:
		err = ErrNotCorrectable
	}
	if err != nil {
		return nil, err
	}
	return ch, nil
}

// checkValue refuses a value the field does not accept.
func checkValue(field, value string, rule correctField) error {
	switch {
	case value == "" && rule.required:
		return ErrEmptyValue
	case len([]rune(value)) > rule.max:
		return ErrTooLong
	case field == "email" && !ValidEmail(value):
		return ErrBadEmail
	}
	return nil
}

// correctRow runs the fixed statement of one field of a core table. A row
// that is not changed was erased or removed since the search: an order by an
// erasure that ran in between, so the answer is ErrErased for an order.
func (s *Service) correctRow(ctx context.Context, websiteID int64, key, value string, rule correctField, field string) error {
	if err := checkValue(field, value, rule); err != nil {
		return err
	}
	now := time.Now().UTC().Format(timeLayout)
	res, err := s.db.Write.ExecContext(ctx, rule.stmt, websiteID, key, value, now)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrErased
	}
	return nil
}
