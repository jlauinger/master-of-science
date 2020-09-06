#!/bin/bash

export GREP_OPTIONS=

echo "file_name,line_number,text,label"

for file in *.go; do
	if [[ $file == *"_test.go" ]]; then
		continue
	fi

	egrep -n 'unsafe\.|reflect\.|uintptr' $file | while read -r finding; do
		line_number=$(echo $finding | cut -d':' -f1)
		text=$(echo $finding | cut -d':' -f2- | tr "," ";")
		echo $file","$line_number","$text",UNKNOWN"
	done
done
