package admin

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/holzcloud/holzcloud-cms/internal/activity"
	"github.com/holzcloud/holzcloud-cms/internal/dsr"
	"github.com/holzcloud/holzcloud-cms/internal/i18n"
	"github.com/holzcloud/holzcloud-cms/internal/plugin"
	"github.com/holzcloud/holzcloud-cms/internal/web"
)

// Correcting a person's records from the Privacy screen (Phase 33.2, DSR-02).
//
// The screen shows, under each record the search found, one small form per
// field that may be corrected; dsr.Correctable decides which those are and
// dsr.Service.Correct makes the change. The form posts the address the
// record was found by, the record's source and key, the field and the new
// value. The service searches that address again on the website of the path
// and corrects only a record it finds there, so nothing else can be reached
// through the form.
//
// Correcting changes stored data, so it needs the publishing right
// (refuseReadOnly, asked first), and somebody without it is not shown the
// forms. It needs no password asked again: a correction can itself be
// corrected.
//
// Every correction writes one privacy.correct entry naming the source, the
// key and the field. A contact field's entry carries no value: the old
// address is the very data the person asked to have replaced, and the log
// outlives the correction. An invoice field's entry carries the old and the
// new value, because the original of an invoice has to stay ascertainable.
//
// A refusal goes back to the screen with a sentence, as a refused erasure
// does, so the answer is the same with and without JavaScript.

// privacyCorrectAddress is where the forms post.
func privacyCorrectAddress(websiteID int64) string { return privacyAddress(websiteID) + "/berichtigen" }

// correctRefusals are the sentences of a refused correction.
var correctRefusals = []struct {
	err      error
	sentence string
}{
	{dsr.ErrBadEmail, i18n.N("This is not one e-mail address.")},
	{dsr.ErrEmptyValue, i18n.N("This field cannot be empty.")},
	{dsr.ErrTooLong, i18n.N("This value is too long.")},
	{dsr.ErrErased, i18n.N("This order was erased and cannot be corrected.")},
	{dsr.ErrNotCorrectable, i18n.N("This field cannot be corrected here.")},
	{dsr.ErrNotFound, i18n.N("This record was not found under this address on this website.")},
	{dsr.ErrBadValue, i18n.N("This value does not fit the field: a country is a two-letter code such as CH, and a Swiss postal code has four digits.")},
	{plugin.ErrPrivacyFailed, i18n.N("The plugin did not accept this value.")},
}

// correctRefusal is the sentence of a refused correction, "" for an error
// that is not a refusal.
func correctRefusal(err error) string {
	for _, r := range correctRefusals {
		if errors.Is(err, r.err) {
			return r.sentence
		}
	}
	return ""
}

// HandlePrivacyCorrect corrects one field of one of the person's records.
func (h *Handler) HandlePrivacyCorrect(w http.ResponseWriter, r *http.Request) error {
	ws, ok, err := h.privacyWebsite(w, r)
	if err != nil || !ok {
		return err
	}
	if err := r.ParseForm(); err != nil {
		return err
	}
	ctx := r.Context()
	email := strings.TrimSpace(r.PostForm.Get("email"))
	requestID, _ := strconv.ParseInt(r.PostForm.Get("anfrage"), 10, 64)
	requestID = max(requestID, 0)
	back := privacyAddress(ws.ID) + "?" + eraseQuery(email, requestID)
	if h.refuseReadOnly(w, r, back) {
		return nil
	}
	source := r.PostForm.Get("source")
	key := r.PostForm.Get("key")
	field := r.PostForm.Get("field")

	ch, err := h.dataRights.Correct(ctx, ws.ID, email, source, key, field, r.PostForm.Get("value"))
	if err != nil {
		sentence := correctRefusal(err)
		if sentence == "" {
			return err
		}
		web.SetFlashError(h.sm, ctx, sentence)
		return h.redirect(w, r, back)
	}

	meta := map[string]any{"source": ch.Source, "key": ch.Key, "field": ch.Field, "invoice": ch.Invoice}
	if ch.Invoice {
		// AO § 146(4) forbids changing a record in a way that leaves its
		// original content no longer ascertainable. An order is the seller's
		// copy of an invoice, so an invoice field keeps its old value here,
		// in the activity log, beside the new one. The log's manual purge
		// stays the operator's to use. An erasure of the person that
		// reduces the order removes both values from these entries again
		// (dsr's stripCorrections), and so does the retention run once the
		// order's period has ended (dsr's stripExpiredCorrections). The e-mail address is never an invoice
		// field, so no address reaches the log this way.
		meta["old"], meta["new"] = ch.Old, ch.New
	}
	websiteID := ws.ID
	h.LogActivity(r, activity.Entry{
		Action:     activity.ActionPrivacyCorrect,
		EntityType: "website",
		EntityID:   ws.ID,
		WebsiteID:  &websiteID,
		Metadata:   meta,
	})
	if ch.Invoice {
		web.SetFlashSuccess(h.sm, ctx, "Corrected. An invoice already issued with the old details may need a corrected invoice.")
	} else {
		web.SetFlashSuccess(h.sm, ctx, "Corrected.")
	}
	// A corrected address finds the record only under the new one, so the
	// screen comes back with that.
	if ch.Field == "email" {
		back = privacyAddress(ws.ID) + "?" + eraseQuery(ch.New, requestID)
	}
	return h.redirect(w, r, back)
}
