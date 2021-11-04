#!/bin/bash

source .env

cat > ./nodegroup.yaml <<_EOF_
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig
metadata:
  name: "mxtransporter-cluster"
  version: "$EKS_VERSION"
  region: "$EKS_REGION"
managedNodeGroups:
  - name: "mxtransporter-node-group"
    privateNetworking: true
    labels:
      role: service
    instanceType: "$NODE_INSTANCE_TYPE"
    desiredCapacity: 2
    minSize: 2
    maxSize: 5
    volumeSize: 20
    ssh:
      publicKeyName: "$SSH_KEY_NAME"
    tags:
      # EC2 tags required for cluster-autoscaler auto-discovery
      k8s.io/cluster-autoscaler/enabled: "true"
      k8s.io/cluster-autoscaler/dev-camplat-secure-cluster: "owned"
_EOF_