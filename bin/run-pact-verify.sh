#! /usr/bin/env bash
export PACT_FILE=$1
export PROVIDER_URL=$2
export SETUP_SERVER_URL=$3

bundle exec rake pact:verify:javascript
