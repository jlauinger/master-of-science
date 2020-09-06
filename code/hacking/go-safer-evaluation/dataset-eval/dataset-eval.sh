#!/bin/bash

EVAL_DIR=/root/code/hacking/go-safer-evaluation
FILES=$EVAL_DIR/dataset-eval-files.txt

>&2 echo "..examining go-safer"
for file in $(cat $FILES); do
	>&2 echo "   $file"
	cd $(dirname $file)
	echo "--- $file ---"
	go-safer . 2>&1
	echo
done > $EVAL_DIR/dataset-eval-results-gosafer.txt

>&2 echo "..examining go vet"
for file in $(cat $FILES); do
	>&2 echo "   $file"
	cd $(dirname $file)
	echo "--- $file ---"
	go vet 2>&1
	echo
done > $EVAL_DIR/dataset-eval-results-govet.txt

>&2 echo "..examining gosec"
for file in $(cat $FILES); do
	>&2 echo "   $file"
	cd $(dirname $file)
	echo "--- $file ---"
	gosec . 2>&1
	echo
done > $EVAL_DIR/dataset-eval-results-gosec.txt

