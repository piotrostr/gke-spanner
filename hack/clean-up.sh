#!/bin/bash

instance=database
cluster_name=gke-connectivity-test
gsa=gcp-spanner-admin@piotrostr-resources.iam.gserviceaccount.com
namespace=default
ksa=kube-spanner-admin

# remove the service account binding
gcloud iam service-accounts remove-iam-policy-binding $gsa \
    --role roles/iam.workloadIdentityUser \
    --member "serviceAccount:piotrostr-resources.svc.id.goog[$namespace/$ksa]"

# delete the spanner instance
gcloud spanner instances delete $instance \
  --quiet

# delete the cluster
gcloud container clusters delete $cluster_name \
  --region=us-east1 \
  --quiet
