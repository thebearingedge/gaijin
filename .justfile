[private]
default:
  @just --list

new module:
  mkdir -p {{module}}
  @pushd {{module}} > /dev/null && go mod init {{module}} && popd
  touch {{module}}/main.go
  go work use {{module}}

test module:
  go test -v {{module}}
