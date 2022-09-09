#!/bin/bash

instance=database
cluster_name=gke-connectivity-test

gcloud spanner instances delete $instance \
  --quiet

gcloud container clusters delete $cluster_name \
  --region=us-east1 \
  --quiet
