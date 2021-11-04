#!/bin/bash

source .env

TARGET_COLLECTION_ARRAY=(`set | grep TARGET_MONGODB_COLLECTION`)

cat > ./helm/values.yaml <<_EOF_
# statefulSet values
image:
  repository: $ECR_REPO
  pullPolicy: IfNotPresent
  tag: "$ECR_REPO_TAG"

replicaCount: 1

containers:
  resources:
    requests:
      cpu: 100m
      memory: 128Mi
    limits:
      cpu: 200m
      memory: 256Mi
  volumeMounts:
    name: dev
    mountPath: $PERSISTENT_VOLUME_DIR

volume:
  accessModes: [ "ReadWriteOnce" ]
  storage: 1Gi

secrets:
  name: $KUBERNETES_SECRET_NAME

# autoscaling values
autoscaling:
  enabled: true
  minReplicas: 1
  maxReplicas: 1
  targetCPUUtilizationPercentage: 80

# service values
service:
  clusterIP: None

# Specify the collection for which you want to get change streams.
targetMongoDBCollections:
_EOF_

if [ TARGET_COLLECTION_ARRAY ] ; then
  i=0
  for TARGET_COLLECTION in ${TARGET_COLLECTION_ARRAY[@]}; do
      echo "  - ${TARGET_COLLECTION#*=}" >> ./helm/values.yaml
      let i++
  done
fi