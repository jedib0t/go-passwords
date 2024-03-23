#!/usr/bin/env bash

TEST_NAME=$1
if [[ -z "${TEST_NAME}" ]]; then
  echo "ERROR: please provide a valid test-name to profile."
  exit 1
fi
PROFILE_TYPE=$2
if [[ -z "${PROFILE_TYPE}" ]]; then
  PROFILE_TYPE="cpu"
fi
if [[ "${PROFILE_TYPE}" != "cpu" ]] && [[ "${PROFILE_TYPE}" != "mem" ]]; then
  echo "ERROR: profile type can only be 'cpu' or 'mem': not '${PROFILE_TYPE}'"
  exit 1
fi

go test -bench=${TEST_NAME} -benchmem -cpuprofile=profile.cpu -memprofile=profile.mem ./...
go tool pprof -http=":8000" ./profile.${PROFILE_TYPE}
