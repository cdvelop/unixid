package unixid

// Primary key
// parámetro opcional:
// "show": el campo se mostrara el usuario por defecto estará oculto
func InputPK(options ...string) *pk {

	in := &pk{
		show: false,
	}

	for _, opt := range options {
		if opt == "show" {
			in.show = true
		}
	}

	return in
}

type pk struct {
	show bool
}

func (p pk) Name() string {
	return "InputPK"
}

func (p pk) HtmlName() string {
	if p.show {
		return "text"
	}
	return "hidden"
}

// representación
func (p pk) BuildContainerView(id, field_name string, allow_skip_completed bool) string {
	var required string
	if !allow_skip_completed {
		required = ` required`
	}
	// p.Number.BuildContainerView.BuildContainerView()
	// return p.BuildHtmlTag(p.HtmlName(), p.Name(), id, field_name, true)
	return `<input type="` + p.HtmlName() + `" id="` + id + `" name="` + field_name + `" data-name="` + p.Name() + `" value=""` + required + `>`
}

func (pk) ValidateField(data_in string, skip_validation bool, options ...string) error {
	if !skip_validation {

		_, err := validateID(data_in)

		return err
	}
	return nil
}

func (p pk) GoodTestData() (out []string) {

	return []string{
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
}

func (p pk) WrongTestData() (out []string) {
	out = []string{"1-1", "-100", "h", "h1", "-1", " ", "", "#", "& ", "% &"}

	return
}
