#!/bin/bash

cluster_name=gke-connectivity-test
region=us-east1

gcloud container clusters create-auto $cluster_name \
   --region=$region

gcloud container clusters get-credentials $cluster_name \
   --region=$region
