package rest

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/door2doc/d2d-uploader/pkg/uploader/db"
)

// VisitorRecord defines a single mutation on a visit.
type VisitorRecord struct {
	// Identity
	Locatie      string `json:"code_locatie"`
	Afdeling     string `json:"code_afdeling"`
	Bezoeknummer int    `json:"bezoeknummer"`
	MutatieID    int    `json:"mutatie_id"`
	Kamer        string `json:"kamer,omitempty"`
	Bed          string `json:"bed,omitempty"`
	Leeftijd     string `json:"leeftijd,omitempty"`

	// Process
	Aangemeld     *time.Time `json:"dt_aangemeld,omitempty"`
	Binnenkomst   *time.Time `json:"dt_binnenkomst,omitempty"`
	Triage        *time.Time `json:"dt_triage,omitempty"`
	NaarKamer     *time.Time `json:"dt_naar_kamer,omitempty"`
	BijArts       *time.Time `json:"dt_bij_arts,omitempty"`
	ArtsKlaar     *time.Time `json:"dt_arts_klaar,omitempty"`
	GereedOpname  *time.Time `json:"dt_gereed_opname,omitempty"`
	Vertrek       *time.Time `json:"dt_vertrek,omitempty"`
	Einde         *time.Time `json:"dt_einde,omitempty"`
	MutatieEinde  *time.Time `json:"dt_mutatie_einde,omitempty"`
	Mutatiestatus string     `json:"code_mutatiestatus,omitempty"`
	Vervallen     bool       `json:"is_vervallen"`

	// Misc
	Ingangsklacht     string `json:"code_ingangsklacht,omitempty"`
	Urgentie          string `json:"code_urgentie,omitempty"`
	Specialisme       string `json:"code_specialisme,omitempty"`
	Herkomst          string `json:"code_herkomst,omitempty"`
	Vervoerder        string `json:"code_vervoerder,omitempty"`
	Ontslagbestemming string `json:"code_ontslagbestemming,omitempty"`
	OpnameAfdeling    string `json:"code_opnameafdeling,omitempty"`
	OpnameSpecialisme string `json:"code_opnamespecialisme,omitempty"`
}

func (v *VisitorRecord) fromDB(r *db.VisitorRecord, loc *time.Location) error {
	var err error
	v.Bezoeknummer = r.Bezoeknummer
	v.MutatieID = r.MutatieID
	v.Locatie = r.Locatie
	v.Afdeling = r.Afdeling
	v.Kamer = r.Kamer
	v.Bed = r.Bed
	v.Ingangsklacht = r.Ingangsklacht
	v.Urgentie = r.Urgentie
	v.Specialisme = r.Specialisme
	v.Herkomst = r.Herkomst
	v.Vervoerder = r.Vervoerder
	v.Ontslagbestemming = r.Ontslagbestemming
	v.OpnameAfdeling = r.OpnameAfdeling
	v.OpnameSpecialisme = r.OpnameSpecialisme

	if !r.Aangemeld.IsZero() {
		aangemeld := r.Aangemeld.Format("2006-01-02 15:04")
		t, err := time.ParseInLocation("2006-01-02 15:04", aangemeld, loc)
		if err == nil {
			v.Aangemeld = &t
		}
	}
	v.Binnenkomst, err = datumTijd(r.BinnenkomstDatum, r.BinnenkomstTijd, loc)
	if err != nil {
		return err
	}
	v.Triage, err = datumTijdRef(v.Binnenkomst, r.TriageTijd, loc)
	if err != nil {
		return err
	}
	v.NaarKamer, err = datumTijdRef(v.Binnenkomst, r.NaarKamerTijd, loc)
	if err != nil {
		return err
	}
	v.BijArts, err = datumTijdRef(v.Binnenkomst, r.BijArtsTijd, loc)
	if err != nil {
		return err
	}
	v.ArtsKlaar, err = datumTijdRef(v.Binnenkomst, r.ArtsKlaarTijd, loc)
	if err != nil {
		return err
	}
	v.GereedOpname, err = datumTijdRef(v.Binnenkomst, r.GereedOpnameTijd, loc)
	if err != nil {
		return err
	}
	v.Vertrek, err = datumTijdRef(v.Binnenkomst, r.VertrekTijd, loc)
	if err != nil {
		return err
	}
	v.Einde, err = datumTijdRef(v.Binnenkomst, r.EindTijd, loc)
	if err != nil {
		return err
	}
	v.MutatieEinde, err = datumTijdRef(v.Binnenkomst, r.MutatieEindTijd, loc)
	if err != nil {
		return err
	}
	v.Mutatiestatus = r.Mutatiestatus
	v.Vervallen = r.Vervallen

	if v.Binnenkomst != nil && !r.Geboortedatum.IsZero() {
		leeftijd := age(v.Binnenkomst, r.Geboortedatum)
		leeftijd = leeftijd / 10
		v.Leeftijd = strconv.Itoa(leeftijd)
	}

	return nil
}

// VisitorRecordFromDB converts a database record to a visitor record.
func VisitorRecordFromDB(r *db.VisitorRecord, loc *time.Location) (*VisitorRecord, error) {
	res := &VisitorRecord{}
	if err := res.fromDB(r, loc); err != nil {
		return nil, err
	}
	return res, nil
}

// VisitorRecordsFromDB converts multiple database records into visitor records.
func VisitorRecordsFromDB(rs []db.VisitorRecord, loc *time.Location) ([]VisitorRecord, error) {
	res := make([]VisitorRecord, len(rs))
	for i := range rs {
		err := res[i].fromDB(&rs[i], loc)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func age(now *time.Time, birthday time.Time) int {
	years := now.Year() - birthday.Year()
	if now.YearDay() < birthday.YearDay() {
		years--
	}
	return years
}

var (
	reParseTijd = regexp.MustCompile(`^(\d?\d:\d\d)(:\d\d)?$`)
)

func normalizeDatum(d string) string {
	t := strings.Index(d, "T")
	if t >= 0 {
		return d[:t]
	}

	return d
}

func normalizeTijd(t string) (string, error) {
	parsedTijd := reParseTijd.FindAllStringSubmatch(t, 1)
	if len(parsedTijd) != 1 || len(parsedTijd[0]) < 2 {
		return "", fmt.Errorf("unrecognized time format: %q", t)
	}
	return parsedTijd[0][1], nil
}

func datumTijd(datum, tijd string, location *time.Location) (*time.Time, error) {
	if datum == "" || tijd == "" {
		return nil, nil
	}

	datum = normalizeDatum(datum)

	nt, err := normalizeTijd(tijd)
	if err != nil {
		return nil, err
	}

	t, err := time.ParseInLocation("2006-01-02 15:04", datum+" "+nt, location)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func datumTijdRef(ref *time.Time, tijd string, location *time.Location) (*time.Time, error) {
	if ref == nil {
		return nil, nil
	}
	if tijd == "" {
		return nil, nil
	}

	nt, err := normalizeTijd(tijd)
	if err != nil {
		return nil, err
	}

	dt := fmt.Sprintf("%d-%02d-%02d %s", ref.Year(), ref.Month(), ref.Day(), nt)
	t, err := time.ParseInLocation("2006-01-02 15:04", dt, location)
	if err != nil {
		return nil, err
	}

	if ref.After(t) {
		t = t.AddDate(0, 0, 1)
	}

	return &t, err
}
