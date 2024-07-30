package unixid

import "testing"

func Test_IdpkTABLA(t *testing.T) {
	// table_name := "user"

	TestData := map[string]struct {
		table_name string
		fieldName  string
		ExpectedID bool
		ExpectedPK bool
	}{
		"id corresponde a tabla sin guion":                               {table_name: "usuario", fieldName: "idusuario", ExpectedID: true, ExpectedPK: true},
		"id corresponde a tabla guion bajo _":                            {table_name: "especialidad", fieldName: "id_especialidad", ExpectedID: true, ExpectedPK: true},
		"campo con solo id es pk":                                        {table_name: "user", fieldName: "id", ExpectedID: true, ExpectedPK: true},
		"corresponde a tabla y key contiene parte el nombre de la tabla": {table_name: "especialidad", fieldName: "especialidades", ExpectedID: false, ExpectedPK: false},
		"id fk de otra tabla sin guion":                                  {table_name: "usuario", fieldName: "idfactura", ExpectedID: true, ExpectedPK: false},
		"id fk de otra tabla con guion bajo _":                           {table_name: "usuario", fieldName: "id_factura", ExpectedID: true, ExpectedPK: false},
		"no primary key presente":                                        {table_name: "usuario", fieldName: "factura", ExpectedID: false, ExpectedPK: false},
		"menos de 2 caracteres id no presente":                           {table_name: "usuario", fieldName: "i", ExpectedID: false, ExpectedPK: false},
	}

	for testName, dt := range TestData {

		h, err := NewHandler()
		if err != nil {
			t.Fatal(err)
		}

		t.Run((testName), func(t *testing.T) {
			pk, pk_this_table := h.FieldType(dt.table_name, dt.fieldName)

			if pk != dt.ExpectedID {
				t.Fail()
			}
			if pk_this_table != dt.ExpectedPK {
				t.Fail()
			}
		})

	}
}
