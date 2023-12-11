module github.com/cdvelop/unixid

go 1.20

require (
	github.com/cdvelop/model v0.0.78
	github.com/cdvelop/timeserver v0.0.26
)

require (
	github.com/cdvelop/strings v0.0.8 // indirect
	github.com/cdvelop/timetools v0.0.28 // indirect
)

replace github.com/cdvelop/model => ../model

replace github.com/cdvelop/timetools => ../timetools

replace github.com/cdvelop/timeserver => ../timeserver
