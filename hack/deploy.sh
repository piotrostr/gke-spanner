#!/bin/bash

image_name=gcr.io/piotrostr-resources/spanner-go-experiment

gcloud auth configure-docker

docker build -t $image_name . && \
  docker push $image_name && \
  kubectl apply -f .
