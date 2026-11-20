package admin

import (
	"context"
	"errors"

	"github.com/holzcloud/holzcloud-cms/internal/dsr"
)

// The privacy area as an assistant reaches it (Phase 33.2, DSR-05).
//
// One operation and no more: the record of processing activities, which
// describes the website and holds no person's data. Searching a person,
// exporting, correcting or erasing their data and reading privacy requests
// stay in the admin, as 33.1 D-16 decided for opt-outs.

// OpProcessingRecord is the website's record of processing activities, as
// the Record of processing screen generates it, in lang or, when lang is
// empty, in the website's main language.
func (h *Handler) OpProcessingRecord(ctx context.Context, websiteID int64, lang string) (dsr.ProcessingRecord, error) {
	if h.dataRights == nil {
		return dsr.ProcessingRecord{}, errors.New("the record of processing is not available on this server")
	}
	ws, err := h.domains.GetWebsite(ctx, websiteID)
	if err != nil {
		return dsr.ProcessingRecord{}, err
	}
	if ws == nil {
		return dsr.ProcessingRecord{}, errors.New("there is no such website")
	}
	if lang == "" {
		lang = websiteLang(ws)
	}
	in, err := h.recordInput(ctx, ws, lang)
	if err != nil {
		return dsr.ProcessingRecord{}, err
	}
	return dsr.BuildRecord(in), nil
}
