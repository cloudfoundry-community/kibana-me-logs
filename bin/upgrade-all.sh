#!/bin/bash

EXPECTED_START_COMMAND="kibana-me-logs"

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
  echo $org_name $space_name $app_name
}

app_urls=$(cf curl /v2/apps | jq -r -c ".resources[].metadata.url")
for app_url in $app_urls; do
  app=$(cf curl $app_url | jq -r -c .)
  detected_start_command=$(echo $app | jq -r -c .entity.detected_start_command)
  if [[ "${detected_start_command}" == "${EXPECTED_START_COMMAND}" ]]; then
    upgrade_app $app
  fi
done
