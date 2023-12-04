module github.com/cdvelop/unixid

go 1.20

require (
	github.com/cdvelop/model v0.0.76
	github.com/cdvelop/timeserver v0.0.24
)

require (
	github.com/cdvelop/strings v0.0.7 // indirect
	github.com/cdvelop/timetools v0.0.25 // indirect
)

replace github.com/cdvelop/model => ../model

replace github.com/cdvelop/timetools => ../timetools

replace github.com/cdvelop/timeserver => ../timeserver
