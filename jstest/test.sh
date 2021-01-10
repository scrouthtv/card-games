#!/bin/bash

cd $(dirname $0)
pushd .. > /dev/null
go test -run ExportProps
popd > /dev/null

for f in $(ls test*.mjs); do
	out=$(node "$f")
	if [ $? -ne 0 ]; then
		printf "Test %s failed:\n" "$f"
		echo "$out"
		exit 1
	else
		printf "OK %s\n" "$f"
	fi
done
