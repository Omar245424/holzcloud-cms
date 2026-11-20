package ai

import (
	"context"
	"errors"

	"github.com/holzcloud/holzcloud-cms/internal/dsr"
)

// The privacy area (Phase 33.2, DSR-05).
//
// An assistant may read a website's record of processing activities: it
// describes what the website does with personal data and holds none. It may
// not reach a person's data — no search by address, no export, no request,
// no opt-out, no erasure — which stays in the admin, behind a person who
// checked who is asking (33.1 D-16, carried over).

// privacyOps is what the privacy tools need from the admin handler.
type privacyOps interface {
	OpProcessingRecord(ctx context.Context, websiteID int64, lang string) (dsr.ProcessingRecord, error)
}

var errNoPrivacyOps = errors.New("the record of processing is not available on this server")

func privacyOpsOf(d Deps) (privacyOps, error) {
	ops, ok := d.Ops.(privacyOps)
	if !ok {
		return nil, errNoPrivacyOps
	}
	return ops, nil
}

func privacyTools(d Deps) []Tool {
	return []Tool{getProcessingRecord(d)}
}

func getProcessingRecord(d Deps) Tool {
	return Tool{
		Name: "get_processing_record",
		Description: "The record of processing activities of a website (GDPR Art. 30, Swiss FADP Art. 12), " +
			"as the admin's Settings → Privacy → Record of processing screen generates it from what the " +
			"website uses, the legal frameworks it follows and its retention periods: the controller, one " +
			"activity per processing with purpose, data subjects, categories of data, legal basis, " +
			"recipients, transfers and retention, the security measures, and which exemption each framework " +
			"allows. Cells marked as to be filled in are facts the operator must add (hosting and mail " +
			"provider, server country, data protection officer, TLS at the proxy). It contains no person's " +
			"data; a person's data is not reachable over this connection.",
		InputSchema: Schema{
			Type: "object",
			Properties: map[string]Property{
				"website": websiteArg,
				"lang": {Type: "string", Description: "language of the record's texts, such as de, fr, it or en; " +
					"left out, the website's main language"},
			},
			Required: []string{"website"},
		},
		Run: func(c Call) (any, error) {
			var a struct {
				Website int64  `json:"website"`
				Lang    string `json:"lang"`
			}
			if err := c.Into(&a); err != nil {
				return nil, err
			}
			ws, err := visibleWebsite(c, d, a.Website)
			if err != nil {
				return nil, err
			}
			ops, err := privacyOpsOf(d)
			if err != nil {
				return nil, err
			}
			return ops.OpProcessingRecord(c.Ctx, ws.ID, a.Lang)
		},
	}
}
