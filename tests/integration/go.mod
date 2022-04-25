module github.com/iden3/merkletree-proof/tests/integration

go 1.18

require (
	github.com/iden3/go-iden3-crypto v0.0.13
	github.com/iden3/go-merkletree-sql v1.0.1
	github.com/iden3/merkletree-proof v0.0.0
	github.com/stretchr/testify v1.7.1
)

require (
	github.com/davecgh/go-spew v1.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)

replace github.com/iden3/merkletree-proof => ../../
