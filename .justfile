[private]
default:
  @just --list

new module:
  #!/bin/sh
  mkdir -p {{module}}
  pushd {{module}} > /dev/null && go mod init {{module}} && popd
  touch {{module}}/main.go
  go work use {{module}}

test module:
  go test -v {{module}}

cover module:
  #!/bin/sh
  MOD=$(basename {{module}})
  go test -v $MOD -coverprofile .coverage/$MOD.out
  go tool cover -html=.coverage/$MOD.out -o .coverage/$MOD.html

add module package:
  @pushd {{module}} && go get {{package}}
