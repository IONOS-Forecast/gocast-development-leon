#!/bin/bash

PROJECT_ROOT=$(realpath $(dirname "${BASH_SOURCE[0]}")/..)

cities=("berlin" "hamburg" "munich")

echo "removing existing weather records"
for c in "${cities[@]}"; do
  rm -f "resources/pg/data/$c".csv
done;

echo "generating CSV files to COPY"
for c in "${cities[@]}"; do
  for i in {1..7}; do
    jq --arg city "$c" -r '.weather | del(.[].fallback_source_ids) | .[].city+=$city | map([.[]]) | .[:-1] | .[] | @csv' "resources/weather_records/$c"_0"$i"-09-orig.json >> "resources/pg/data/$c".csv
  done;
  printf "%s" "$(< resources/pg/data/$c.csv)" > "resources/pg/data/$c".csv
done;

