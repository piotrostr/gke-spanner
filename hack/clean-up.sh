#!/bin/bash

instance=database
cluster_name=gke-connectivity-test

gcloud spanner instances delete $instance \
  --config=regional-us-central1

gcloud container clusters delete $cluster_name \
  --regoin=us-east1
