set dotenv-load

[private]
default:
  @just --list

new module:
  #!/bin/sh
  mkdir -p {{module}}
  pushd {{module}} > /dev/null && go mod init {{module}} && popd
  echo "package main" >> {{module}}/main.go
  go work use {{module}}

test module:
  go test -v {{module}}

tdd module:
  gow -c -w {{module}} test -v {{module}}

cover module:
  #!/bin/sh
  MOD=$(basename {{module}})
  go test -v $MOD/... -coverprofile .coverage/$MOD.out
  go tool cover -html=.coverage/$MOD.out -o .coverage/$MOD.html

add module package:
  @pushd {{module}} && go get {{package}}

run module:
  go run -C {{module}} main.go

watch module:
  gow run -C {{module}} main.go

build module:
  go build -C {{module}} -o "$(pwd)/.bin/$(basename {{module}})" main.go

clean module:
  pushd {{module}} && go mod tidy

compose_up:
  docker compose up --build --detach

compose_down:
  docker compose down -v

db_create name:
  migrate create -ext sql -dir with-db/migrations {{name}}

db_up module:
  migrate -source file://./{{module}}/migrations -database $DATABASE_URL up

db_down module:
  migrate -source file://./{{module}}/migrations -database $DATABASE_URL down --all

init module:
  terraform -chdir={{module}} init

plan module:
  terraform -chdir={{module}} plan

apply module:
  terraform -chdir={{module}} apply

destroy module:
  terraform -chdir={{module}} destroy
