#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

gcloud beta functions deploy dailytoggl \
  --project ${PROJECT} \
  --quiet \
  --region ${REGION} \
  --allow-unauthenticated \
  --entry-point DailyToggl \
  --memory 256MB \
  --runtime go113 \
  --timeout 60s \
  --set-env-vars AUTH_TOKEN=${AUTH_TOKEN} \
  --max-instances 1 \
  --trigger-http