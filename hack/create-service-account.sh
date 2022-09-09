#!/bin/bash

gcloud iam service-accounts create spanner-rw \
  --display-name "Spanner Read Write" \
  --project piotrostr-resources

gcloud projects add-iam-policy-binding piotrostr-resources \
  --member "serviceAccount:spanner-rw@piotrostr-resources.iam.gserviceaccount.com" \
  --role "roles/spanner.databaseAdmin"
