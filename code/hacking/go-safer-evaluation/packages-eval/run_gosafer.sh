#!/bin/bash

echo "file_name,line_number,message"

go-safer . 2>&1 | while read -r line; do 
	file_name=$(echo $line | cut -d':' -f1 | sed 's!.*/!!')
	line_number=$(echo $line | cut -d':' -f2)
	message=$(echo $line | cut -d':' -f4- | tr "," ";")

	echo $file_name','$line_number','$message
done
