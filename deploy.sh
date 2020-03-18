#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail


<<<$ENV_VARS_YAML cat > env-vars.yaml

gcloud beta functions deploy dailytoggl \
  --project ${PROJECT} \
  --quiet \
  --region ${REGION} \
  --allow-unauthenticated \
  --entry-point DailyToggl \
  --memory 256MB \
  --runtime go113 \
  --timeout 60s \
  --env-vars-file env-vars.yaml \
  --max-instances 1 \
  --trigger-http