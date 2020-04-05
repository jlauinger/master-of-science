#!/bin/bash

DOWNLOAD_DIR=/home/johannes/studium/s14/masterarbeit/code/download
PROJECT_RESULTS_FILE=unsafe_usages.csv
TOTAL_RESULTS_FILE=$DOWNLOAD_DIR/total_unsafe_usages.csv

echo "Module;Unsafe Pointer Usages" > $TOTAL_RESULTS_FILE

for PROJECT in $DOWNLOAD_DIR/*; do
  echo -e "\nAnalyzing $PROJECT..."

  if [ ! -f $PROJECT/go.mod ]; then
    echo "skipping because go.mod does not exist."
    continue
  fi

  cd $PROJECT


  go mod vendor 2>/dev/null 1>/dev/null
  if [ "$?" -ne "0" ]; then
    echo "skipping because go mod vendor failed."
    continue
  fi

  echo "Module;Unsafe Pointer Usages" > $PROJECT/$PROJECT_RESULTS_FILE

  for MODULE in $(go mod vendor -v 2>&1 | grep -v "#" | sort | uniq); do
    LINES=$(ag unsafe.Pointer vendor/$MODULE | wc -l)
    echo "$LINES $MODULE"
    echo "$MODULE;$LINES" >> $PROJECT/$PROJECT_RESULTS_FILE
    echo "$MODULE;$LINES" >> $TOTAL_RESULTS_FILE
  done
done
