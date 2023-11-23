package unixid_test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/cdvelop/unixid"
)

func TestIdToDate(t *testing.T) {
	now := time.Now()
	now_nano := now.UnixNano()
	now_expected := now.Format("2006-01-02 15:04")

	testData := map[string]struct {
		inputData string
		expected  string
	}{
		"id sin punto ok":             {fmt.Sprint(now_nano), now_expected},
		"id 2022 ":                    {"1643317560659315800", "2022-01-27 18:06"},
		"id 2023 ":                    {"1672949057314961800", "2023-01-05 17:04"},
		"id 2022-01-27":               {"1643319071971938900", "2022-01-27 18:31"},
		"id 2021-03-24":               {"1648131093739755800", "2022-03-24 11:11"},
		"id 2022-01-19":               {"1643318806368317300", "2022-01-27 18:26"},
		"id 2022-01-19 con .":         {"1643318806368317300.5", "2022-01-27 18:26"},
		"id 2023-03-30 con .":         {"1680184020131482400.0", "2023-03-30 10:47"},
		"error id 2023-03-30 con E":   {"16801E4020131482400.0", "validateID error id contiene caracteres no válidos"},
		"error id 2023-03-30 con .. ": {"16801840201.31482400.0", "validateID error id contiene más de un punto"},
		"error sin data de entrada":   {"", "validateID error id contiene caracteres no válidos"},
	}

	for prueba, data := range testData {
		t.Run((prueba + ": " + data.inputData), func(t *testing.T) {

			resp, err := unixid.UnixNanoToStringDate(data.inputData)

			if err != "" {
				resp = err
			}

			if resp != data.expected {
				fmt.Println("ERROR prueba:", prueba)
				fmt.Printf("- [%v] resultado:\n[%v]\n-expectativa:\n[%v]\n\n", data.inputData, resp, data.expected)
				log.Fatal()
			}
		})
	}

}
