#!/bin/bash

EXPECTED_START_COMMAND="kibana-me-logs"
RESULTS_PER_PAGE=${RESULTS_PER_PAGE:-100}

if [[ "$(which cf)X" == "X" ]]; then
  echo "Please install cf"
  exit 1
fi
if [[ "$(which jq)X" == "X" ]]; then
  echo "Please install jq - http://stedolan.github.io/jq/download/"
  exit 1
fi

debug() {
  if [[ -n ${DEBUG} && ${DEBUG} != '0' ]];
    then echo >&2 '>> ' "$*"
  fi
}

function cf_curl() {
  set -e
  url=$1
  md5=$(echo "${url}" | md5sum | cut -f1 -d " ")
        path="${tmpdir}/${md5}"
  if [[ ! -f $path ]]; then
    debug "No cached data found - cf curl ${url}"
    cf curl "${url}" > ${path}
  fi
  echo ${path}
}

tmpdir=$(mktemp -d)
mkdir -p $tmpdir
trap 'rm -rf ${tmpdir:?nothing to remove}; exit' INT TERM QUIT EXIT
debug "set up workspace directory ${tmpdir}"

echo "Searching for kibana-me-logs apps... This may take a bit"
next_url="/v2/apps?results-per-page=${RESULTS_PER_PAGE}"
while [[ "${next_url}" != "null" ]]; do
  debug "Finding Apps from ${next_url}"
  app_urls=$(cat $(cf_curl ${next_url}) | jq -r -c ".resources[].metadata.url")
  for app_url in $app_urls; do
    debug "Getting app data from $app_url"
    app_name=$(cat $(cf_curl $app_url) | jq -r -c .entity.name)
    space_url=$(cat $(cf_curl $app_url) | jq -r -c .entity.space_url)
    space_name=$(cat $(cf_curl $space_url) | jq -r -c .entity.name)
    org_url=$(cat $(cf_curl $space_url) | jq -r -c .entity.organization_url)
    org_name=$(cat $(cf_curl $org_url) | jq -r -c .entity.name)
    detected_start_command=$(cat $(cf_curl $app_url) | jq -r -c .entity.detected_start_command)
    echo -n "Checking app '${app_name} in '${org_name}/${space_name}'..."
    if [[ "${detected_start_command}" == "${EXPECTED_START_COMMAND}" ]]; then
      if [[ -z ${DRY_RUN} || ${DRY_RUN} == '0' ]]; then
        echo -e "\033[0;32mUPGRADING\033[0m"
        cf target -o $org_name -s $space_name; cf push $app_name
      else
        echo -e "\033[0;33mWOULD UPGRADE\033[0m (dry-run mode enabled)"
      fi
    else
      echo
    fi
  done
  next_url=$(cat $(cf_curl ${next_url}) | jq -r -c ".next_url")
done
