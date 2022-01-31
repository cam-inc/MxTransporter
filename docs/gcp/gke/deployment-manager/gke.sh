#!/bin/bash

source .env

cat > ./deployment-manager/gke.yaml <<_EOF_
imports:
  - path: templates/gke.py
    name: gke.py
resources:
  - name: mxtransporter-cluster
    type: gke.py
    properties:
      clusterLocationType: Regional
      region: $GKE_REGION
      cluster:
        name: $GKE_CLUSTER_NAME
        description: mxtransporter-cluster
        network: $GKE_CLUSTER_NETWORK
        subnetwork: $GKE_CLUSTER_SUBNETWORK
        initialClusterVersion: $GKE_CLUSTER_VERSION
        nodePools:
          - name: mxtransporter-pool
            initialNodeCount: 1
            locations:
              - $GKE_NODE_LOCATION_1
              - $GKE_NODE_LOCATION_2
              - $GKE_NODE_LOCATION_3
            autoscaling:
              enabled: True
              minNodeCount: 1
              maxNodeCount: 2
            management:
              autoUpgrade: False
              autoRepair: True
            config:
              machineType: $GKE_NODE_MACHINE_TYPE
              localSsdCount: 0
              diskSizeGb: 100
              preemptible: False
              diskType: pd-standard
              oauthScopes:
                - https://www.googleapis.com/auth/cloud-platform
              metadata:
                disable-legacy-endpoints: "true"
              serviceAccount: mxtransporter@$PROJECT_NAME.iam.gserviceaccount.com
        networkPolicy:
          enabled: True
        loggingService: logging.googleapis.com/kubernetes
        monitoringService: monitoring.googleapis.com/kubernetes
        privateCluster: True
        masterIpv4CidrBlock: $GKE_CLUSTER_MASTER_IP_CIDR_BLOCK
        masterAuthorizedNetworksConfig:
          enabled: False
          cidrBlocks: []
        ipAllocationPolicy:
          useIpAliases: True
          clusterIpv4CidrBlock: $GKE_PODS_IP_CIDR_BLOCK
        maintenancePolicy:
          window:
            recurringWindow:
              window:
                # AM9:00-PM15:00 weekday
                startTime: 2020-01-01T00:00:00Z
                endTime: 2020-01-01T06:00:00Z
              recurrence: FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR
_EOF_