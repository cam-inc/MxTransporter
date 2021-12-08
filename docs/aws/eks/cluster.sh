#!/bin/bash

source .env

cat > ./cluster.yaml <<_EOF_
apiVersion: eksctl.io/v1alpha5
cloudWatch:
  clusterLogging:
    enableTypes: ["all"]
kind: ClusterConfig
metadata:
  name: "mxtransporter-cluster"
  version: "$EKS_VERSION"
  region: "$EKS_REGION"
privateCluster:
  enabled: false
vpc:
  subnets:
    private:
      $EKS_AVAILABILITY_ZONE_A:
        id: "$EKS_SUBNET_A"
      $EKS_AVAILABILITY_ZONE_B:
        id: "$EKS_SUBNET_B"
  clusterEndpoints:
    publicAccess: true
    privateAccess: true
#  you don't have to write the following description, so default setting ["0.0.0.0/0"]
#  publicAccessCIDRs: ["0.0.0.0/0"]
  manageSharedNodeSecurityGroupRules: true
secretsEncryption:
  keyARN: "$EKS_KEY_ARN"
_EOF_