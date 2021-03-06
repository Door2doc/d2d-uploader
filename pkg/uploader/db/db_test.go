package db

import (
	"context"
	"database/sql"
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/denisenkom/go-mssqldb"
	"github.com/lib/pq"
)

const (
	PqlDSN = "postgres://pguser:pwd@localhost:5436/pgdb?sslmode=disable"
	MssDSN = "sqlserver://sa:MyPassw0rd@localhost:1433/?database=upload"

	TestDSN = PqlDSN
)

var (
	_ pq.Driver
	_ mssql.Driver
)

func u(t time.Time) time.Time {
	if t.IsZero() {
		return t
	}
	return t.In(time.UTC)
}

func TestExecuteVisitorQuery(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for name, test := range map[string]struct {
		Query string
		Want  VisitorRecords
		Err   error
	}{
		"correct": {
			Query: `select * from correct where id = 1`,
			Want: []VisitorRecord{{
				Bezoeknummer:      328996,
				MutatieID:         1091568,
				Locatie:           "A",
				Afdeling:          "seh",
				Aangemeld:         time.Date(2017, time.July, 13, 13, 0, 0, 0, time.UTC),
				BinnenkomstDatum:  "2017-07-13",
				BinnenkomstTijd:   "23:18",
				TriageTijd:        "",
				NaarKamerTijd:     "23:18",
				BijArtsTijd:       "02:40",
				ArtsKlaarTijd:     "",
				GereedOpnameTijd:  "02:06",
				VertrekTijd:       "04:34",
				EindTijd:          "04:34",
				Kamer:             "",
				Bed:               "",
				Ingangsklacht:     "Pneumonie",
				Specialisme:       "04",
				Urgentie:          "",
				Vervoerder:        "2",
				Geboortedatum:     time.Date(1977, time.July, 24, 12, 0, 0, 0, time.UTC),
				OpnameAfdeling:    "",
				OpnameSpecialisme: "",
				Herkomst:          "",
				Ontslagbestemming: "",
				Vervallen:         false,
			}},
		},
		"check all columns": {
			Query: `select * from correct where id = 2`,
			Want: []VisitorRecord{{
				Bezoeknummer:      1,
				MutatieID:         2,
				Locatie:           "locatie",
				Afdeling:          "seh",
				Aangemeld:         time.Date(2018, time.July, 4, 12, 4, 0, 0, time.UTC),
				BinnenkomstDatum:  "binnenkomstdatum",
				BinnenkomstTijd:   "binnenkomsttijd",
				TriageTijd:        "aanvangtriagetijd",
				NaarKamerTijd:     "naarkamertijd",
				BijArtsTijd:       "eerstecontacttijd",
				ArtsKlaarTijd:     "artsklaartijd",
				GereedOpnameTijd:  "gereedopnametijd",
				VertrekTijd:       "vertrektijd",
				EindTijd:          "eindtijd",
				Kamer:             "kamer",
				Bed:               "bed",
				Ingangsklacht:     "ingangsklacht",
				Specialisme:       "specialisme",
				Urgentie:          "triage",
				Vervoerder:        "vervoerder",
				Geboortedatum:     time.Date(1977, time.July, 24, 12, 0, 0, 0, time.UTC),
				OpnameAfdeling:    "opnameafdeling",
				OpnameSpecialisme: "opnamespecialisme",
				Herkomst:          "herkomst",
				Ontslagbestemming: "ontslagbestemming",
				Vervallen:         false,
			}},
		},
		"missing columns": {
			Query: `select null as hello`,
			Err: &SelectionError{
				Missing: []string{
					ColBezoeknummer.Name,
					ColMutatieID.Name,
					ColLocatie.Name,
					ColAfdeling.Name,
					ColAangemeld.Name,
					ColBinnenkomstDatum.Name,
					ColBinnenkomstTijd.Name,
					ColTriageTijd.Name,
					ColNaarKamerTijd.Name,
					ColBijArtsTijd.Name,
					ColArtsKlaarTijd.Name,
					ColGereedOpnameTijd.Name,
					ColVertrekTijd.Name,
					ColEindTijd.Name,
					ColMutatieEindTijd.Name,
					ColMutatieStatus.Name,
					ColKamer.Name,
					ColBed.Name,
					ColIngangsklacht.Name,
					ColSpecialisme.Name,
					ColUrgentie.Name,
					ColVervoerder.Name,
					ColGeboortedatum.Name,
					ColOpnameAfdeling.Name,
					ColOpnameSpecialisme.Name,
					ColHerkomst.Name,
					ColOntslagbestemming.Name,
					ColVervallen.Name,
				},
				Got: []string{"hello"},
			},
		},
		"duplicate columns": {
			Query: `select null as hello, null as hello`,
			Err:   ErrDuplicateColumnNames,
		},
	} {
		t.Run(name, func(t *testing.T) {
			tx, cancel := setup(ctx, t)
			defer cancel()

			got, err := ExecuteVisitorQuery(ctx, tx, test.Query, time.Second)

			for i := range got {
				got[i].Aangemeld = u(got[i].Aangemeld)
				got[i].Geboortedatum = u(got[i].Geboortedatum)
			}

			if !reflect.DeepEqual(got, test.Want) {
				t.Errorf("ExecuteVisitorQuery() == \n\t%v, got \n\t%v", test.Want, got)
			}
			if !reflect.DeepEqual(err, test.Err) {
				t.Errorf("ExecuteVisitorQuery() == \n\t_, %#v; got \n\t_, %#v", test.Err, err)
			}
		})
	}
}

func TestExecuteVisitorQueryPermutations(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	tx, cancel := setup(ctx, t)
	defer cancel()

	parts := []string{
		"1 AS sehid",
		"2 AS sehmutid",
		"'locatie' AS locatie",
		"'afdeling' AS afdeling",
		"NULL AS aangemaakt",
		"'binnenkomstdatum' AS binnenkomstdatum",
		"'binnenkomsttijd' AS binnenkomsttijd",
		"'triagetijd' AS triagetijd",
		"'naarkamertijd' AS naarkamertijd",
		"'eerstecontacttijd' AS eerstecontacttijd",
		"'artsklaartijd' AS artsklaartijd",
		"'afdelinggebeldtijd' AS afdelinggebeldtijd",
		"'gereedopnametijd' AS gereedopnametijd",
		"'vertrektijd' AS vertrektijd",
		"'eindtijd' AS eindtijd",
		"'mutatieeindtijd' AS mutatieeindtijd",
		"'mutatiestatus' AS mutatiestatus",
		"'kamer' AS kamer",
		"'bed' AS bed",
		"'ingangsklacht' AS ingangsklacht",
		"'specialisme' AS specialisme",
		"'triage' AS triage",
		"'vervoerder' AS vervoerder",
		"'bestemming' AS bestemming",
		"NULL AS geboortedatum",
		"'opnameafdeling' AS opnameafdeling",
		"'opnamespecialisme' AS opnamespecialisme",
		"'herkomst' AS herkomst",
		"'ontslagbestemming' AS ontslagbestemming",
		"0 AS vervallen",
		"'x' AS ignoreme_1",
		"'x' AS ignoreme_2",
		"'x' AS ignoreme_3",
		"2 AS ignoreme_4",
		"3 AS ignoreme_5",
		"4 AS ignoreme_6",
		"'x' AS ignoreme_7",
		"'x' AS ignoreme_8",
		"'x' AS ignoreme_9",
	}

	query := generateQuery(parts)
	got, err := ExecuteVisitorQuery(ctx, tx, query, time.Second)
	if err != nil {
		t.Fatal(err)
	}

	want := VisitorRecords{{
		Bezoeknummer:      1,
		MutatieID:         2,
		Locatie:           "locatie",
		Afdeling:          "afdeling",
		Aangemeld:         time.Time{},
		BinnenkomstDatum:  "binnenkomstdatum",
		BinnenkomstTijd:   "binnenkomsttijd",
		TriageTijd:        "triagetijd",
		NaarKamerTijd:     "naarkamertijd",
		BijArtsTijd:       "eerstecontacttijd",
		ArtsKlaarTijd:     "artsklaartijd",
		GereedOpnameTijd:  "gereedopnametijd",
		VertrekTijd:       "vertrektijd",
		EindTijd:          "eindtijd",
		MutatieEindTijd:   "mutatieeindtijd",
		Mutatiestatus:     "mutatiestatus",
		Kamer:             "kamer",
		Bed:               "bed",
		Ingangsklacht:     "ingangsklacht",
		Specialisme:       "specialisme",
		Urgentie:          "triage",
		Vervoerder:        "vervoerder",
		Geboortedatum:     time.Time{},
		OpnameAfdeling:    "opnameafdeling",
		OpnameSpecialisme: "opnamespecialisme",
		Herkomst:          "herkomst",
		Ontslagbestemming: "ontslagbestemming",
		Vervallen:         false,
	}}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ExecuteVisitorQuery() == \n\t%#v, got \n\t%#v", want, got)
	}
}

func generateQuery(parts []string) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	rnd.Shuffle(len(parts), func(i, j int) {
		parts[i], parts[j] = parts[j], parts[i]
	})
	query := "select " + strings.Join(parts, ",")
	return query
}

func setup(ctx context.Context, t *testing.T) (*sql.Tx, context.CancelFunc) {
	if testing.Short() {
		t.Skip("uses database")
	}

	driver := TestDSN[:strings.Index(TestDSN, "://")]
	db, err := sql.Open(driver, TestDSN)
	if err != nil {
		t.Fatal(err)
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	return tx, func() {
		if err := tx.Rollback(); err != nil {
			t.Error(err)
		}
		if err := db.Close(); err != nil {
			t.Error(err)
		}
	}
}

func TestVisitorRecordsAsTable(t *testing.T) {
	v := VisitorRecords{
		{}, {},
	}
	res := v.AsTable()
	if res == "" {
		t.Fail()
	}
}

func timeP(s string) *time.Time {
	res, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}

	return &res
}

func TestExecuteRadiologieQueryPermutations(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	tx, cancel := setup(ctx, t)
	defer cancel()

	parts := []string{
		"1 AS sehid",
		"2 AS ordernr",
		"'status' as status",
		"'1977-07-24 12:00'::timestamp as startdatumtijd",
		"null as einddatumtijd",
		"'module' as module",
		"'specialisme' as specialisme",
		"'x' AS ignoreme_1",
	}

	query := generateQuery(parts)
	got, err := ExecuteRadiologieQuery(ctx, tx, query, time.Second)
	if err != nil {
		t.Fatal(err)
	}

	want := RadiologieOrders{{
		Bezoeknummer: 1,
		Ordernummer:  2,
		Status:       "status",
		Start:        timeP("1977-07-24T12:00:00+00:00"),
		Eind:         nil,
		Module:       "module",
	}}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ExecuteRadiologieQuery() == \n\t%v, want \n\t%v", got, want)
	}
}

func TestExecuteLabQueryPermutations(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	tx, cancel := setup(ctx, t)
	defer cancel()

	parts := []string{
		"1 AS sehid",
		"2 AS ordernr",
		"'status' as status",
		"'1977-07-24 12:00'::timestamp as startdatumtijd",
		"null as einddatumtijd",
		"'module' as module",
		"'specialisme' as specialisme",
		"'x' AS ignoreme_1",
	}

	query := generateQuery(parts)
	got, err := ExecuteLabQuery(ctx, tx, query, time.Second)
	if err != nil {
		t.Fatal(err)
	}

	want := LabOrders{{
		Bezoeknummer: 1,
		Ordernummer:  2,
		Status:       "status",
		Start:        timeP("1977-07-24T12:00:00+00:00"),
		Eind:         nil,
	}}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ExecuteRadiologieQuery() == \n\t%v, want \n\t%v", got, want)
	}
}

func TestExecuteConsultQueryPermutations(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	tx, cancel := setup(ctx, t)
	defer cancel()

	parts := []string{
		"1 AS sehid",
		"2 AS ordernr",
		"'status' as status",
		"'1977-07-24 12:00'::timestamp as startdatumtijd",
		"null as einddatumtijd",
		"'module' as module",
		"'specialisme' as specialisme",
		"'x' AS ignoreme_1",
	}

	query := generateQuery(parts)
	got, err := ExecuteConsultQuery(ctx, tx, query, time.Second)
	if err != nil {
		t.Fatal(err)
	}

	want := ConsultOrders{{
		Bezoeknummer: 1,
		Ordernummer:  2,
		Status:       "status",
		Start:        timeP("1977-07-24T12:00:00+00:00"),
		Eind:         nil,
		Specialisme:  "specialisme",
	}}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ExecuteRadiologieQuery() == \n\t%v, want \n\t%v", got, want)
	}
}
