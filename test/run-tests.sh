export GOPATH=`pwd`
go get github.com/stretchr/testify/assert
go get github.com/codegangsta/cli
go test pact-verify
go test pact-verify/run
