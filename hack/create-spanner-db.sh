#!/bin/bash

database=database
instance=database
processing_units=100
config=regional-us-central1

gcloud spanner instances create $instance \
  --processing-units=$processing_units \
  --config=$config \
  --description="gke-connectivity-test"

gcloud spanner databases create $database \
  --instance=$instance
