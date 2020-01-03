module github.com/bif/bif-sdk-go

go 1.13

require (
	github.com/bif/go-bif v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.0.0-20191209160850-c0dbc17a3553
)

replace (
	github.com/bif/bif-wasm => ../bif-wasm
	github.com/bif/go-bif => ../go-bif
)
