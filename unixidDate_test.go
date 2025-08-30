package unixid_test

import (
	"testing"
	"time"

	. "github.com/cdvelop/tinystring"
	"github.com/cdvelop/unixid"
)

func TestIdToDate(t *testing.T) {
	now := time.Now()
	now_nano := now.UnixNano()
	now_expected := now.Format("2006-01-02 15:04")

	uid, err := unixid.NewUnixID()
	if err != nil {
		t.Fatal(err)
	}

	testData := map[string]struct {
		inputData string
		expected  string
	}{
		"id sin punto ok":             {Fmt("%d", now_nano), now_expected},
		"id 2022 ":                    {"1643317560659315800", "2022-01-27 21:06"},
		"id 2023 ":                    {"1672949057314961800", "2023-01-05 20:04"},
		"id 2022-01-27":               {"1643319071971938900", "2022-01-27 21:31"},
		"id 2021-03-24":               {"1648131093739755800", "2022-03-24 14:11"},
		"id 2022-01-19":               {"1643318806368317300", "2022-01-27 21:26"},
		"id 2022-01-19 con .":         {"1643318806368317300.5", "2022-01-27 21:26"},
		"id 2023-03-30 con .":         {"1680184020131482400.0", "2023-03-30 13:47"},
		"error id 2023-03-30 con E":   {"16801E4020131482400.0", "Character Invalid Not Supported"},
		"error id 2023-03-30 con .. ": {"16801840201.31482400.0", "Invalid Format Found More Point"},
		"error sin data de entrada":   {"", "Character Invalid Not Supported"},
	}

	for prueba, data := range testData {
		t.Run((prueba + ": " + data.inputData), func(t *testing.T) {

			resp, err := uid.UnixNanoToStringDate(data.inputData)

			if err != nil {
				resp = err.Error()
			}

			if resp != data.expected {
				t.Fatalf("ERROR prueba:%v\n- [%v] resultado:\n[%v]\n-expectativa:\n[%v]\n\n", prueba, data.inputData, resp, data.expected)
			}
		})
	}

}
