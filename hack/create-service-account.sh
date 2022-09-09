#!/bin/bash

# based on https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity

gsa_name=gcp-spanner-admin
gsa=gcp-spanner-admin@piotrostr-resources.iam.gserviceaccount.com
ksa=kube-spanner-admin
project_id=piotrostr-resources
namespace=default

# create kubernetes service account (KSA)
kubectl create serviceaccount kube-spanner-admin \
    --namespace default

# create google service account (GSA)
gcloud iam service-accounts create $gsa_name \
  --project $project_id

# add GSA permissions
gcloud projects add-iam-policy-binding $project_id \
  --member "serviceAccount:$gsa" \
  --role "roles/spanner.databaseAdmin"

# add GSA workloadIdentityUser for KSA
gcloud iam service-accounts add-iam-policy-binding $gsa \
    --role roles/iam.workloadIdentityUser \
    --member "serviceAccount:piotrostr-resources.svc.id.goog[$namespace/$ksa]"

# annotate the KSA
kubectl annotate serviceaccount $ksa \
    --namespace $namespace \
    iam.gke.io/gcp-service-account=$gsa
