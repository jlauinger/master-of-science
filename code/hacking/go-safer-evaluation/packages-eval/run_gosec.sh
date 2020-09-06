#!/bin/bash

export GREP_OPTIONS=

echo "file_name,line_number,message"

gosec -quiet . 2>/dev/null | grep '^\[' | sed -r "s/\x1B\[(([0-9]{1,2})?(;)?([0-9]{1,2})?)?[m,K,H,f,J]//g" | sed -r "s/\[|\]|\(|\)//g" | while read -r line; do 
	file_name=$(echo $line | cut -d':' -f1 | sed 's!.*/!!')
	line_number=$(echo $line | cut -d':' -f2 | cut -d' ' -f1)
	message=$(echo $line | cut -d' ' -f3- | tr "," ";")

	echo $file_name','$line_number','$message
done
