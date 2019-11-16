#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

RED='\033[0;31m'
GREEN='\033[0;32m'
INVERTED='\033[7m'
NC='\033[0m' # No Color

echo "? go mod verify"
verifyResult=$(go mod verify)
if [[ $? != 0 ]]; then
	echo -e "${RED}✗ go mod verify\n$verifyResult${NC}"
	exit 1
else echo -e "${GREEN}√ go mod verify${NC}"
fi

echo "? go test"
go test ./...
if [[ $? != 0 ]]; then
	echo -e "${RED}✗ go test\n${NC}"
	exit 1
else echo -e "${GREEN}√ go test${NC}"
fi

goFilesToCheck=$(find . -type f -name "*.go" | egrep -v "\/vendor\/|_*/automock/|_*/testdata/|_*export_test.go")

goFmtResult=$(echo "${goFilesToCheck}" | xargs -L1 go fmt)
if [[ $(echo ${#goFmtResult}) != 0 ]]; then
  echo -e "${RED}✗ go fmt${NC}\n$goFmtResult${NC}"
  exit 1;
else echo -e "${GREEN}√ go fmt${NC}"
fi

echo "? ./test/gqlgen.sh"
${DIR}/test/gqlgen.sh
if [[ -n "$(git status -s test/schema)" ]]; then
		echo -e "${RED}✗ gqlgen.sh modified some files, schema and code are out-of-sync${NC}"
		git status -s test/schema
		exit 1
else echo -e "${GREEN}√ ./test/gqlgen.sh${NC}"
fi
