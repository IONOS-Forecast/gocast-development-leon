#!/bin/bash

PROJECT_ROOT=$(realpath $(dirname "${BASH_SOURCE[0]}")/..)

#cities=("berlin" "hamburg")
cities=()
while IFS= read -r line || [[ -n "$line" ]]; do
    cities+=("$line")
done < $(realpath $(dirname "${BASH_SOURCE[0]}")/cities.txt)

echo "Cities array contains: ${cities[@]}"

echo "removing existing weather records"
for c in "${cities[@]}"; do
  rm -f "resources/pg/data/$c".csv
done;

mkdir -p "../resources/pg/data"

echo "generating CSV files to COPY"
for c in "${cities[@]}"; do
  for i in {0..7}; do
    jq --arg city "$c" -r '.weather | del(.[].fallback_source_ids) | .[].city+=$city | map([.[]]) | .[:-1] | .[] | @csv' "resources/weather_records/$c"_"$i"-orig.json >> "resources/pg/data/$c".csv
  done;
  printf "%s" "$(< resources/pg/data/$c.csv)" > "resources/pg/data/$c".csv
done;
