#!/bin/bash

EXPECTED_START_COMMAND="kibana-me-logs"
RESULTS_PER_PAGE=${RESULTS_PER_PAGE:-10}

if [[ "$(which cf)X" == "X" ]]; then
  echo "Please install cf"
  exit 1
fi
if [[ "$(which jq)X" == "X" ]]; then
  echo "Please install jq - http://stedolan.github.io/jq/download/"
  exit 1
fi

function upgrade_app {
  app=$1
  app_name=$(echo $app | jq -r -c .entity.name)
  space_url=$(echo $app | jq -r -c .entity.space_url)
  space_name=$(cf curl $space_url | jq -r -c .entity.name)
  org_url=$(cf curl $space_url | jq -r -c .entity.organization_url)
  org_name=$(cf curl $org_url | jq -r -c .entity.name)
  echo "cf target -o $org_name -s $space_name; cf push $app_name"
  cf target -o $org_name -s $space_name; cf push $app_name
}

next_url="/v2/apps?results-per-page=${RESULTS_PER_PAGE}"
while [[ "${next_url}" != "null" ]]; do
  cf curl ${next_url}
  app_urls=$(cf curl ${next_url} | jq -r -c ".resources[].metadata.url")
  for app_url in $app_urls; do
    app=$(cf curl $app_url | jq -r -c .)
    detected_start_command=$(echo $app | jq -r -c .entity.detected_start_command)
    if [[ "${detected_start_command}" == "${EXPECTED_START_COMMAND}" ]]; then
      upgrade_app $app
    fi
  done
  next_url=$(cf curl ${next_url} | jq -r -c ".next_url")
done
