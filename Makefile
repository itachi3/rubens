# Prepend the project directory to the system GOPATH
# so that import path resolution will prioritize
# our third party snapshots.

GOPATH := ${PWD}/_vendor:${GOPATH}
export GOPATH

default: build

build: vet
    go build -v -o ./bin/main_app ./src/main_app