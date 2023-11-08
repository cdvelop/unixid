module github.com/cdvelop/unixid

go 1.20

require (
	github.com/cdvelop/model v0.0.66
	github.com/cdvelop/timeserver v0.0.1
)

require github.com/cdvelop/timetools v0.0.2 // indirect

replace github.com/cdvelop/model => ../model

replace github.com/cdvelop/timetools => ../timetools

replace github.com/cdvelop/timeserver => ../timeserver
