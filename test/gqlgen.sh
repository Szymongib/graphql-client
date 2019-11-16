#!/usr/bin/env bash

echo "Generating code from GraphQL schema..."

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

cd ${DIR}/schema
go run ${DIR}/hack/gqlgen.go -v --config ./config.yaml
