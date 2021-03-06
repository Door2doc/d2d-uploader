package rest

import (
	"reflect"
	"testing"
	"time"

	"github.com/door2doc/d2d-uploader/pkg/uploader/db"
)

func tm(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	t = t.In(time.UTC)
	return &t
}

func u(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	tu := t.In(time.UTC)
	return &tu
}

func TestVisitorRecordFromDB(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Amsterdam")
	if err != nil {
		t.Fatal(err)
	}

	for name, test := range map[string]struct {
		Given *db.VisitorRecord
		Want  *VisitorRecord
	}{
		"minimal": {
			Given: &db.VisitorRecord{Bezoeknummer: 12, MutatieID: 100, Locatie: "qqq"},
			Want:  &VisitorRecord{Bezoeknummer: 12, MutatieID: 100, Locatie: "qqq"},
		},
		"Aangemeld": {
			Given: &db.VisitorRecord{Aangemeld: time.Date(2017, time.July, 24, 12, 0, 0, 0, loc)},
			Want:  &VisitorRecord{Aangemeld: tm("2017-07-24T10:00:00Z")},
		},
		"Binnenkomst": {
			Given: &db.VisitorRecord{
				BinnenkomstDatum: "2017-01-01T12:21:25.65Z",
				BinnenkomstTijd:  "00:20:00",
			},
			Want: &VisitorRecord{
				Binnenkomst: tm("2017-01-01T00:20:00+01:00"),
			},
		},
		"Binnenkomst date-time": {
			Given: &db.VisitorRecord{
				BinnenkomstDatum: "2017-01-01",
				BinnenkomstTijd:  "00:20:00",
			},
			Want: &VisitorRecord{
				Binnenkomst: tm("2017-01-01T00:20:00+01:00"),
			},
		},
		"TriageTijd": {
			Given: &db.VisitorRecord{
				BinnenkomstDatum: "2017-01-01",
				BinnenkomstTijd:  "00:20:00",
				TriageTijd:       "00:24",
			},
			Want: &VisitorRecord{
				Binnenkomst: tm("2017-01-01T00:20:00+01:00"),
				Triage:      tm("2017-01-01T00:24:00+01:00"),
			},
		},
		"TriageTijd, geen aankomst": {
			Given: &db.VisitorRecord{
				TriageTijd: "00:24",
			},
			Want: &VisitorRecord{},
		},
		"TriageTijd, volgende dag": {
			Given: &db.VisitorRecord{
				BinnenkomstDatum: "2017-01-01",
				BinnenkomstTijd:  "23:20:00",
				TriageTijd:       "00:24",
			},
			Want: &VisitorRecord{
				Binnenkomst: tm("2017-01-01T23:20:00+01:00"),
				Triage:      tm("2017-01-02T00:24:00+01:00"),
			},
		},
		"NaarKamer": {
			Given: &db.VisitorRecord{
				BinnenkomstDatum: "2017-01-01",
				BinnenkomstTijd:  "00:20:00",
				NaarKamerTijd:    "00:20",
			},
			Want: &VisitorRecord{
				Binnenkomst: tm("2017-01-01T00:20:00+01:00"),
				NaarKamer:   tm("2017-01-01T00:20:00+01:00"),
			},
		},
		"EersteContact": {
			Given: &db.VisitorRecord{
				BinnenkomstDatum: "2017-01-01",
				BinnenkomstTijd:  "00:20:00",
				BijArtsTijd:      "03:06",
			},
			Want: &VisitorRecord{
				Binnenkomst: tm("2017-01-01T00:20:00+01:00"),
				BijArts:     tm("2017-01-01T03:06:00+01:00"),
			},
		},
		"ArtsKlaar": {
			Given: &db.VisitorRecord{
				BinnenkomstDatum: "2017-01-01",
				BinnenkomstTijd:  "00:20:00",
				ArtsKlaarTijd:    "02:23",
			},
			Want: &VisitorRecord{
				Binnenkomst: tm("2017-01-01T00:20:00+01:00"),
				ArtsKlaar:   tm("2017-01-01T02:23:00+01:00"),
			},
		},
		"GereedOpname": {
			Given: &db.VisitorRecord{
				BinnenkomstDatum: "2017-01-01",
				BinnenkomstTijd:  "00:20:00",
				GereedOpnameTijd: "04:52",
			},
			Want: &VisitorRecord{
				Binnenkomst:  tm("2017-01-01T00:20:00+01:00"),
				GereedOpname: tm("2017-01-01T04:52:00+01:00"),
			},
		},
		"Vertrek": {
			Given: &db.VisitorRecord{
				BinnenkomstDatum: "2017-01-01",
				BinnenkomstTijd:  "00:20:00",
				VertrekTijd:      "07:21",
			},
			Want: &VisitorRecord{
				Binnenkomst: tm("2017-01-01T00:20:00+01:00"),
				Vertrek:     tm("2017-01-01T07:21:00+01:00"),
			},
		},
		"Eind": {
			Given: &db.VisitorRecord{
				BinnenkomstDatum: "2017-01-01",
				BinnenkomstTijd:  "00:20:00",
				EindTijd:         "07:21",
			},
			Want: &VisitorRecord{
				Binnenkomst: tm("2017-01-01T00:20:00+01:00"),
				Einde:       tm("2017-01-01T07:21:00+01:00"),
			},
		},
		"Mutatie Eindtijd/Status": {
			Given: &db.VisitorRecord{
				BinnenkomstDatum: "2017-01-01",
				BinnenkomstTijd:  "00:20:00",
				MutatieEindTijd:  "03:15",
				Mutatiestatus:    "NH",
			},
			Want: &VisitorRecord{
				Binnenkomst:   tm("2017-01-01T00:20:00+01:00"),
				MutatieEinde:  tm("2017-01-01T03:15:00+01:00"),
				Mutatiestatus: "NH",
			},
		},
		"Kamer": {
			Given: &db.VisitorRecord{
				Kamer: "D12",
			},
			Want: &VisitorRecord{
				Kamer: "D12",
			},
		},
		"Bed": {
			Given: &db.VisitorRecord{
				Bed: "01",
			},
			Want: &VisitorRecord{
				Bed: "01",
			},
		},
		"Ingangsklacht": {
			Given: &db.VisitorRecord{
				Ingangsklacht: "kortademigheid volwassene",
			},
			Want: &VisitorRecord{
				Ingangsklacht: "kortademigheid volwassene",
			},
		},
		"Specialisme": {
			Given: &db.VisitorRecord{
				Specialisme: "INT",
			},
			Want: &VisitorRecord{
				Specialisme: "INT",
			},
		},
		"Triage": {
			Given: &db.VisitorRecord{
				Urgentie: "04",
			},
			Want: &VisitorRecord{
				Urgentie: "04",
			},
		},
		"Herkomst": {
			Given: &db.VisitorRecord{
				Herkomst: "EIG",
			},
			Want: &VisitorRecord{
				Herkomst: "EIG",
			},
		},
		"Vervoerder": {
			Given: &db.VisitorRecord{
				Vervoerder: "AMB",
			},
			Want: &VisitorRecord{
				Vervoerder: "AMB",
			},
		},
		"Ontslagbestemming": {
			Given: &db.VisitorRecord{
				Ontslagbestemming: "NH",
			},
			Want: &VisitorRecord{
				Ontslagbestemming: "NH",
			},
		},
		"OpnameAfdeling": {
			Given: &db.VisitorRecord{
				OpnameAfdeling: "HAOA",
			},
			Want: &VisitorRecord{
				OpnameAfdeling: "HAOA",
			},
		},
		"OpnameSpecialisme": {
			Given: &db.VisitorRecord{
				OpnameSpecialisme: "INT",
			},
			Want: &VisitorRecord{
				OpnameSpecialisme: "INT",
			},
		},
		"Leeftijd": {
			Given: &db.VisitorRecord{
				BinnenkomstDatum: "2017-01-01",
				BinnenkomstTijd:  "00:20:00",
				Geboortedatum:    time.Date(1977, time.July, 24, 12, 0, 0, 0, time.UTC),
			},
			Want: &VisitorRecord{
				Binnenkomst: tm("2017-01-01T00:20:00+01:00"),
				Leeftijd:    "3",
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			if test.Want == nil {
				t.Skip("not defined")
			}

			got, err := VisitorRecordFromDB(test.Given, loc)
			if err != nil {
				t.Fatal(err)
			}

			got.Binnenkomst = u(got.Binnenkomst)
			got.Triage = u(got.Triage)
			got.NaarKamer = u(got.NaarKamer)
			got.BijArts = u(got.BijArts)
			got.ArtsKlaar = u(got.ArtsKlaar)
			got.Vertrek = u(got.Vertrek)
			got.Aangemeld = u(got.Aangemeld)
			got.GereedOpname = u(got.GereedOpname)
			got.Einde = u(got.Einde)
			got.MutatieEinde = u(got.MutatieEinde)

			if !reflect.DeepEqual(got, test.Want) {
				t.Errorf("VisitorRecordFromDB(\n\t%#v) == \n\t%v, got \n\t%v", test.Given, test.Want, got)
			}
		})
	}
}
