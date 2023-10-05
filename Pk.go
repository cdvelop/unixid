package unixid

import (
	"github.com/cdvelop/input"
	"github.com/cdvelop/model"
)

// Primary key
// parámetro opcional:
// "show": el campo se mostrara el usuario por defecto estará oculto
func InputPK(options ...string) *model.Input {

	in := pk{
		show: false,
		Permitted: input.Permitted{
			Numbers: true,
			Minimum: 1,
			Maximum: 20,
		},
	}

	for _, opt := range options {
		if opt == "show" {
			in.show = true
		}
	}

	return &model.Input{
		InputName: in.Name(),
		Tag:       in,
		Validate:  in,
		TestData:  in,
	}
}

type pk struct {
	show bool
	input.Permitted
}

func (p pk) Name() string {
	return "Pk"
}

func (p pk) HtmlName() string {
	if p.show {
		return "text"
	}
	return "hidden"
}

// representación
func (p pk) HtmlTag(id, field_name string, allow_skip_completed bool) string {
	var required string
	if !allow_skip_completed {
		required = ` required`
	}
	// p.Number.HtmlTag.HtmlTag()
	// return p.BuildHtmlTag(p.HtmlName(), p.Name(), id, field_name, true)
	return `<input type="` + p.HtmlName() + `" id="` + id + `" name="` + field_name + `" data-name="` + p.Name() + `"` + required + `>`
}

func (p pk) GoodTestData() (out []string) {

	temp := []string{
		"56988765432",
		"1234567",
		"0",
		"123456789",
		"100",
		"5000",
		"423456789",
		"31",
		"523756789",
		"10000232326263727",
		"29",
		"923726789",
		"3234567",
		"823456789",
		"29",
	}

	for _, v := range temp {
		if len(v) >= p.Minimum && len(v) <= p.Maximum {
			out = append(out, v)
		}
	}

	return
}

func (p pk) WrongTestData() (out []string) {
	out = []string{"1-1", "-100", "h", "h1", "-1", " ", "", "#", "& ", "% &"}

	return
}
