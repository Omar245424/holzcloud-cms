package dsr

// Correction: the operator corrects what a website holds about one person
// (Phase 33.2, DSR-02, research section "Correction per record kind").
//
// What may be corrected depends on what the record is:
//
//   - An order's contact fields, its e-mail address and phone number, are
//     corrected freely. They are how the shop reaches the customer, and the
//     invoice does not print them.
//   - An order's invoice fields — name, company, VAT number, street, postal
//     code, city, country — are corrected too, but as a change that keeps
//     the original on record (Change.Invoice). The order is the seller's
//     copy of an invoice, and AO § 146(4) forbids changing a record in a way
//     that leaves its original content no longer ascertainable. The caller
//     writes old and new value into the activity log; this file only says
//     which fields they are.
//   - An erased order is not corrected at all. The erasure reduced it to the
//     invoice record the law asks to keep (D-03, D-07), and that is one-way.
//   - A legacy form message corrects its name and e-mail address; what the
//     person wrote stays what they wrote.
//   - A plugin's record is corrected by the plugin, through the privacy hook,
//     and only in a field the plugin listed as correctable. A plugin order's
//     fields other than e-mail and phone count as invoice fields, as the
//     erasure keeps them like an order of the shop.
//   - A staff account is corrected on its own screen, by an administrator; a
//     privacy request is the record of a request and a mail a copy of what
//     was sent, and neither is corrected (ErrNotCorrectable).
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
	"fmt"
	"strings"
	"time"

	"github.com/holzcloud/holzcloud-cms/internal/plugin"
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
	// ErrBadValue means the value has a shape the field does not take: a
	// country that is not a two-letter code, a Swiss postal code that is not
	// four digits.
	ErrBadValue = errors.New("dsr: the value does not fit the field")
)

// correctField is how one field is corrected.
type correctField struct {
	// stmt is the whole UPDATE of a core table: $1 website, $2 key, $3
	// value, $4 the time. Empty for a plugin's field.
	stmt string
	// max is the longest value accepted, in characters.
	max int
	// required refuses an empty value.
	required bool
	// invoice marks a field the invoice prints.
	invoice bool
	// unstamped marks a statement of a table without updated_at: it takes
	// no $4.
	unstamped bool
}

// orderFields are the corrections of an order, in the order the screen shows
// them: contact first, then the invoice fields.
var orderFields = []string{"email", "phone", "name", "company", "vat_number", "street", "postal_code", "city", "country"}

// orderCorrections never match an erased order (erased_at IS NULL): the
// erasure is one-way. The checkout requires name, street, postal code, city
// and country and sets no length; the limits here only keep a pasted text
// out of an address field.
var orderCorrections = map[string]correctField{
	"email": {stmt: `UPDATE orders SET email = $3, updated_at = $4
		WHERE website_id = $1 AND number = $2 AND erased_at IS NULL`, max: 254, required: true},
	"phone": {stmt: `UPDATE orders SET phone = $3, updated_at = $4
		WHERE website_id = $1 AND number = $2 AND erased_at IS NULL`, max: 50},
	"name": {stmt: `UPDATE orders SET name = $3, updated_at = $4
		WHERE website_id = $1 AND number = $2 AND erased_at IS NULL`, max: 200, required: true, invoice: true},
	"company": {stmt: `UPDATE orders SET company = $3, updated_at = $4
		WHERE website_id = $1 AND number = $2 AND erased_at IS NULL`, max: 200, invoice: true},
	"vat_number": {stmt: `UPDATE orders SET vat_number = $3, updated_at = $4
		WHERE website_id = $1 AND number = $2 AND erased_at IS NULL`, max: 50, invoice: true},
	"street": {stmt: `UPDATE orders SET street = $3, updated_at = $4
		WHERE website_id = $1 AND number = $2 AND erased_at IS NULL`, max: 200, required: true, invoice: true},
	"postal_code": {stmt: `UPDATE orders SET postal_code = $3, updated_at = $4
		WHERE website_id = $1 AND number = $2 AND erased_at IS NULL`, max: 20, required: true, invoice: true},
	"city": {stmt: `UPDATE orders SET city = $3, updated_at = $4
		WHERE website_id = $1 AND number = $2 AND erased_at IS NULL`, max: 200, required: true, invoice: true},
	"country": {stmt: `UPDATE orders SET country = $3, updated_at = $4
		WHERE website_id = $1 AND number = $2 AND erased_at IS NULL`, max: 2, required: true, invoice: true},
}

// formMessageFields are the corrections of a legacy form message. The table
// has no updated_at, so its statements take three values.
var formMessageFields = []string{"name", "email"}

var formMessageCorrections = map[string]correctField{
	"name": {stmt: `UPDATE form_messages SET name = $3
		WHERE website_id = $1 AND id = $2`, max: 200, required: true, unstamped: true},
	"email": {stmt: `UPDATE form_messages SET email = $3
		WHERE website_id = $1 AND id = $2`, max: 254, required: true, unstamped: true},
}

// pluginRule is the host's own check of a plugin field before the plugin is
// asked. The plugin checks again by its own rules.
func pluginRule(rec Record, name string) correctField {
	rule := correctField{max: 400, required: name == "name" || name == "email"}
	if name == "email" {
		rule.max = 254
	}
	rule.invoice = rec.Kind == plugin.PrivacyKindOrder && name != "email" && name != "phone"
	return rule
}

// fieldValue is the label and value of the named field of a record. A field
// the record does not show (an empty one) gets the core label of its name.
func fieldValue(rec Record, name string) (label, value string) {
	for _, f := range rec.Fields {
		if f.Name == name {
			return f.Label, f.Value
		}
	}
	if l, ok := fieldLabels[name]; ok {
		return l, ""
	}
	return name, ""
}

// Correctable is what the operator may correct of a record, in the order the
// screen shows it. A record that is not corrected here has none.
func Correctable(rec Record) []CorrectableField {
	if rec.Erased {
		return nil
	}
	var names []string
	rule := func(string) correctField { return correctField{} }
	switch {
	case rec.Source == "orders":
		names = orderFields
		rule = func(n string) correctField { return orderCorrections[n] }
	case rec.Source == "form_messages":
		names = formMessageFields
		rule = func(n string) correctField { return formMessageCorrections[n] }
	case strings.HasPrefix(rec.Source, "plugin:"):
		names = rec.Correctable
		rule = func(n string) correctField { return pluginRule(rec, n) }
	default:
		return nil
	}
	out := make([]CorrectableField, 0, len(names))
	for _, name := range names {
		label, value := fieldValue(rec, name)
		out = append(out, CorrectableField{Name: name, Label: label, Value: value, Invoice: rule(name).invoice})
	}
	return out
}

// Correct sets one field of one of the person's records on the website.
//
// email is the address the record was found by; the record is searched for
// again under it, and source and key must name one of the records found
// (ErrNotFound otherwise). An erased record is ErrErased, a field that is not
// corrected here ErrNotCorrectable — for a plugin before the plugin is asked.
// The value is trimmed and checked as the form it came from checks it: an
// address must be one bare address (ErrBadEmail), a required field must hold
// something (ErrEmptyValue), a country is a two-letter code and a Swiss
// postal code four digits (ErrBadValue), and no value may be longer than the
// field allows (ErrTooLong). A plugin's refusal is plugin.ErrPrivacyFailed.
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
	if field == "country" {
		value = strings.ToUpper(value)
	}

	switch {
	case source == "orders":
		rule := orderCorrections[field]
		if err := checkOrderValue(*rec, field, value); err != nil {
			return nil, err
		}
		err = s.correctRow(ctx, websiteID, key, field, value, rule)
	case source == "form_messages":
		err = s.correctRow(ctx, websiteID, key, field, value, formMessageCorrections[field])
	case strings.HasPrefix(source, "plugin:"):
		err = s.correctPlugin(ctx, websiteID, res.Email, *rec, field, value)
	default:
		err = ErrNotCorrectable
	}
	if err != nil {
		return nil, err
	}
	return &Change{Source: source, Key: key, Field: field, Old: allowed.Value, New: value, Invoice: allowed.Invoice}, nil
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

// checkOrderValue is what the checkout checks beyond presence: the country
// is a code (the checkout offers a list of codes), and a Swiss postal code
// has four digits.
func checkOrderValue(rec Record, field, value string) error {
	switch field {
	case "country":
		if len(value) != 2 || !isAsciiUpper(value) {
			return ErrBadValue
		}
	case "postal_code":
		if _, country := fieldValue(rec, "country"); country == "CH" && value != "" && !isDigits(value, 4) {
			return ErrBadValue
		}
	}
	return nil
}

func isAsciiUpper(s string) bool {
	for _, r := range s {
		if r < 'A' || r > 'Z' {
			return false
		}
	}
	return true
}

func isDigits(s string, n int) bool {
	if len(s) != n {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// correctRow runs the fixed statement of one field of a core table. A row
// that is not changed was erased or removed since the search: an order by an
// erasure that ran in between, which is ErrErased.
func (s *Service) correctRow(ctx context.Context, websiteID int64, key, field, value string, rule correctField) error {
	if rule.stmt == "" {
		return ErrNotCorrectable
	}
	if err := checkValue(field, value, rule); err != nil {
		return err
	}
	args := []any{websiteID, key, value}
	if !rule.unstamped {
		args = append(args, time.Now().UTC().Format(timeLayout))
	}
	res, err := s.db.Write.ExecContext(ctx, rule.stmt, args...)
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

// correctPlugin asks the plugin to correct its record. The searched address
// goes with it, so a plugin that keeps the address in more than one place
// (an assembled form's answers) replaces it everywhere.
func (s *Service) correctPlugin(ctx context.Context, websiteID int64, email string, rec Record, field, value string) error {
	if err := checkValue(field, value, pluginRule(rec, field)); err != nil {
		return err
	}
	id := strings.TrimPrefix(rec.Source, "plugin:")
	_, err := s.plugins.Privacy(ctx, websiteID, id, plugin.PrivacyIn{
		Op: plugin.PrivacyCorrect, Email: email, Key: rec.Key, Field: field, Value: value,
	})
	if err != nil {
		return fmt.Errorf("dsr: %s refused the correction: %w", id, err)
	}
	return nil
}
