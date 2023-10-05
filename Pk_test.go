package unixid_test

import (
	"log"
	"testing"

	"github.com/cdvelop/unixid"
)

var (
	modelPrimaryKey = unixid.InputPK()

	dataPrimaryKey = map[string]struct {
		inputData       string
		skip_validation bool
		expected        string
	}{
		"id 1 correcto?":                 {"1624397134562544800", false, ""},
		"id 2 ok?":                       {"1624397172303448900", false, ""},
		"id 3 ok?":                       {"1634394443466878800", false, ""},
		"numero 5 correcto?":             {"5", false, ""},
		"numero 45 correcto?":            {"45", false, ""},
		"id con letra valido?":           {"E624397172303448900", false, "E no es un numero"},
		"primary key se permite vació ?": {"", false, "tamaño mínimo 1 caracteres"},
		"id cero?":                       {"0", false, ""},
	}
)

func Test_InputPrimaryKey(t *testing.T) {

	for prueba, data := range dataPrimaryKey {
		t.Run((prueba + ": " + data.inputData), func(t *testing.T) {
			err := modelPrimaryKey.Validate.ValidateField(data.inputData, data.skip_validation)
			var resp string
			if err != nil {
				resp = err.Error()
			}

			if resp != data.expected {
				log.Println(prueba)
				log.Fatalf("resultado: [%v] expectativa: [%v]\n%v", resp, data.expected, data.inputData)
			}
		})
	}
}

func Test_TagPrimaryKey(t *testing.T) {
	tag := modelPrimaryKey.Tag.HtmlTag("1", "name", true)
	if tag == "" {
		log.Fatalln("ERROR NO TAG RENDERING ")
	}
}

func Test_GoodInputPrimaryKey(t *testing.T) {
	for _, data := range modelPrimaryKey.TestData.GoodTestData() {
		t.Run((data), func(t *testing.T) {
			if ok := modelPrimaryKey.Validate.ValidateField(data, false); ok != nil {
				log.Fatalf("resultado [%v] [%v]", ok, data)
			}
		})
	}
}

func Test_WrongInputPrimaryKey(t *testing.T) {
	for _, data := range modelPrimaryKey.TestData.WrongTestData() {
		t.Run((data), func(t *testing.T) {
			if ok := modelPrimaryKey.Validate.ValidateField(data, false); ok == nil {
				log.Fatalf("resultado [%v] [%v]", ok, data)
			}
		})
	}
}
