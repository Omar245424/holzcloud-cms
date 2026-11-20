package dsr

// The machine-readable answer to an access request (D-02): one JSON file
// holding every record the search found, the plugins it could not ask, and
// the sources that failed. Keys are English and stable; values are the data
// as stored.
//
// This is where the export's shape is decided, and with it what never enters
// it. No source reads a password hash, a second-factor secret, a token hash, a
// recovery code or an order's confirmation token, so none of them can reach a
// Record and none can reach this file. They are credentials, not the
// person's data in any useful sense: handing one over would let whoever holds
// the file act as the account or open the order, and the right of access
// stops where it would harm the rights of others (GDPR Art. 15(4)).

import (
	"sort"
	"time"

	"github.com/holzcloud/holzcloud-cms/internal/i18n"
	"github.com/holzcloud/holzcloud-cms/internal/legal"
)

// Export is the file.
//
// Beside the records it carries what GDPR Art. 15(1) and FADP Art. 25(2)
// ask an answer to say: who the controller is, why each kind of data is
// processed, on what basis, for whom, for how long and where it came from,
// and what the person can do about it. Describe fills that in from the
// website's record of processing, so the answer to one person and the
// record kept for the authority are one text.
type Export struct {
	GeneratedAt string                    `json:"generated_at"`
	Website     ExportWebsite             `json:"website"`
	Email       string                    `json:"email"`
	Controller  map[string]string         `json:"controller"`
	Frameworks  []string                  `json:"frameworks"`
	Categories  []ExportCategory          `json:"categories"`
	Rights      []ExportRight             `json:"rights"`
	Records     map[string][]ExportRecord `json:"records"`
	NotCovered  []ExportHolder            `json:"not_covered"`
	Incomplete  []ExportFailure           `json:"incomplete"`
	// Lang is the language of the descriptive texts; the records are as
	// stored.
	Lang string `json:"lang"`
}

// ExportCategory describes one kind of record the export holds.
type ExportCategory struct {
	Kind string `json:"kind"`
	// Activity is the name of the processing activity of the record.
	Activity   string `json:"activity"`
	Purpose    string `json:"purpose"`
	LegalBasis string `json:"legal_basis"`
	Recipients string `json:"recipients"`
	Retention  string `json:"retention"`
	Source     string `json:"source"`
}

// ExportRight is one sentence about the person's rights under one
// framework.
type ExportRight struct {
	Framework string `json:"framework"`
	Text      string `json:"text"`
}

// Where a kind of record comes from (Art. 15(1)(g) GDPR), catalogue keys.
var (
	sourceProvided = i18n.N("Provided by you.")
	sourceAccount  = i18n.N("Provided by you, or by the administrator who created your account.")
	sourceMail     = i18n.N("Written by this website from what you provided.")
)

// sourceOf is the source sentence of a kind.
func sourceOf(kind string) string {
	switch kind {
	case RetentionAccount:
		return sourceAccount
	case "mail":
		return sourceMail
	}
	return sourceProvided
}

// The rights per framework, with where to complain. The wording follows the
// published privacy notice (internal/privacy/texts/<lang>/privacy), so the
// person reads the same rights in the answer as on the website.
var rightsTexts = map[legal.Framework][]string{
	legal.EU: {
		i18n.N("Under the GDPR you have the right of access (Art. 15), to rectification (Art. 16), to erasure (Art. 17), to restriction of processing (Art. 18), to data portability (Art. 20) and to object (Art. 21)."),
		i18n.N("You can lodge a complaint with a supervisory authority, in particular in the Member State of your habitual residence, your place of work or the place of the alleged infringement (Art. 77 GDPR)."),
	},
	legal.UK: {
		i18n.N("Under the UK GDPR you have the right of access, to rectification, to erasure, to restriction of processing, to data portability and to object (Art. 15 to 21 UK GDPR)."),
		i18n.N("You can complain to the controller first (section 164A DPA 2018), and to the Information Commission, formerly the Information Commissioner's Office (ICO) (section 165 DPA 2018)."),
	},
	legal.CH: {
		i18n.N("Under the FADP you have the right of access (Art. 25), to have your data handed over or transferred (Art. 28), to have it corrected (Art. 32(1)), and to have processing prohibited, disclosure prevented or the data deleted (Art. 32(2))."),
		i18n.N("You can report a possible breach of the data protection rules to the Federal Data Protection and Information Commissioner (FDPIC) (Art. 49 FADP)."),
	},
	legal.US: {
		i18n.N("Under the US state privacy laws such as the CCPA you have the right to know, to delete and to correct your personal information, to opt out of its sale and sharing, and not to be discriminated against for using these rights."),
		i18n.N("If a request is refused, the privacy notice of this website says how to appeal."),
	},
}

// Describe adds to the export what the record of processing says about it:
// the controller, the frameworks, one category per kind of record the
// export holds, and the person's rights, in the record's language.
func (e *Export) Describe(rec ProcessingRecord) {
	t := func(s string) string { return i18n.T(rec.Lang, s) }
	e.Lang = rec.Lang
	e.Controller = map[string]string{}
	for k, v := range rec.Controller {
		e.Controller[k] = v
	}
	e.Frameworks = append([]string{}, rec.Frameworks...)

	present := map[string]bool{}
	for _, recs := range e.Records {
		for _, r := range recs {
			present[r.Kind] = true
		}
	}
	e.Categories = []ExportCategory{}
	seen := map[string]bool{}
	// In the order of the record's activities, so the answer reads like the
	// record; a kind is described by the first activity that produces it.
	for _, a := range rec.Activities {
		for _, kind := range a.Kinds {
			if !present[kind] || seen[kind] {
				continue
			}
			seen[kind] = true
			e.Categories = append(e.Categories, ExportCategory{
				Kind: kind, Activity: a.Name, Purpose: a.Purpose, LegalBasis: a.LegalBasis,
				Recipients: a.Recipients, Retention: a.Retention, Source: t(sourceOf(kind)),
			})
		}
	}
	// A kind no activity describes, such as one a plugin invented: it is
	// named, and what only the operator knows is left to fill in.
	fill := t(fillIn)
	sources := make([]string, 0, len(e.Records))
	for source := range e.Records {
		sources = append(sources, source)
	}
	sort.Strings(sources)
	for _, source := range sources {
		for _, r := range e.Records[source] {
			if seen[r.Kind] {
				continue
			}
			seen[r.Kind] = true
			e.Categories = append(e.Categories, ExportCategory{
				Kind: r.Kind, Activity: fill, Purpose: fill, LegalBasis: fill,
				Recipients: fill, Retention: fill, Source: t(sourceOf(r.Kind)),
			})
		}
	}

	e.Rights = []ExportRight{}
	for _, f := range rec.Frameworks {
		for _, s := range rightsTexts[legal.Framework(f)] {
			e.Rights = append(e.Rights, ExportRight{Framework: f, Text: t(s)})
		}
	}
}

// ExportWebsite names the website the export is about.
type ExportWebsite struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// ExportRecord is one record.
type ExportRecord struct {
	Key    string            `json:"key"`
	Kind   string            `json:"kind"`
	Time   string            `json:"time"`
	Fields map[string]string `json:"fields"`
}

// ExportHolder is a plugin that keeps data for the website and was not
// asked.
type ExportHolder struct {
	Plugin  string `json:"plugin"`
	Name    string `json:"name"`
	Entries int    `json:"entries"`
	Reason  string `json:"reason"`
}

// ExportFailure is a source that could not be searched. The error text is
// the program's own, never a value from a record.
type ExportFailure struct {
	Source string `json:"source"`
	Error  string `json:"error"`
}

// BuildExport turns a search result into the file. Every list is empty
// rather than null, so the file always has the same keys.
func BuildExport(res *Result, websiteID int64, websiteName string, now time.Time) Export {
	out := Export{
		GeneratedAt: now.UTC().Format(time.RFC3339),
		Website:     ExportWebsite{ID: websiteID, Name: websiteName},
		Email:       res.Email,
		Controller:  map[string]string{},
		Frameworks:  []string{},
		Categories:  []ExportCategory{},
		Rights:      []ExportRight{},
		Records:     map[string][]ExportRecord{},
		NotCovered:  []ExportHolder{},
		Incomplete:  []ExportFailure{},
	}
	for _, rec := range res.Records {
		fields := make(map[string]string, len(rec.Fields))
		for _, f := range rec.Fields {
			fields[f.Name] = f.Value
		}
		var when string
		if !rec.Time.IsZero() {
			when = rec.Time.UTC().Format(time.RFC3339)
		}
		out.Records[rec.Source] = append(out.Records[rec.Source], ExportRecord{
			Key: rec.Key, Kind: rec.Kind, Time: when, Fields: fields,
		})
	}
	for _, h := range res.NotCovered() {
		out.NotCovered = append(out.NotCovered, ExportHolder{
			Plugin: h.PluginID, Name: h.Name, Entries: h.Entries, Reason: h.Reason,
		})
	}
	for _, f := range res.Failed {
		out.Incomplete = append(out.Incomplete, ExportFailure{Source: f.Source, Error: f.Err.Error()})
	}
	return out
}
