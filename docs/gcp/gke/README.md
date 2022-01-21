# Build MxTransporter in GKE

This section describes the procedure for building MxTransporter in the GKE environment.<br>
All commands to create each GCP resources are wrapped in ```Makefile```.

## Prepare
### Command
Need to use following commands. The version listed is the verified version.

```
bq       v2.0.71
docker   v20.10.8
gcloud   v355.0.0
kubectl  v1.22.1
helm     v3.6.3
make     v3.81
```

### Environment variables files
Before starting the construction, create ```.env``` and ```.secrets.env``` in the current directory by referring to ```.env.template``` and ```secrets.env.template```.

If you want to export change streams to BigQuery or Pub/Sub, write the following description in ```.secrets.env```.

```
EXPORT_DESTINATION=bigquery

or

EXPORT_DESTINATION=pubsub

or

EXPORT_DESTINATION=bigquery,pubsub
```

### BigQuery schema  (optional)
If you want to export change streams to BigQuery, specify the following table schema.

Table schema
```
[
    {
      "mode": "NULLABLE",
      "name": "id",
      "type": "STRING"
    },
    {
      "mode": "NULLABLE",
      "name": "operationType",
      "type": "STRING"
    },
    {
      "mode": "NULLABLE",
      "name": "clusterTime",
      "type": "TIMESTAMP"
    },
    {
      "mode": "NULLABLE",
      "name": "fullDocument",
      "type": "STRING"
    },
    {
      "mode": "NULLABLE",
      "name": "ns",
      "type": "STRING"
    },
    {
      "mode": "NULLABLE",
      "name": "documentKey",
      "type": "STRING"
    },
    {
      "mode": "NULLABLE",
      "name": "updateDescription",
      "type": "STRING"
    }
]
```

## Procedure
**Optional: Setup BigQuery**

You need to create a dataset and a table.

Create a dataset.
It's interactive, so specify the dataset name.

```
$ make create-bigquery-dataset
```

Create a table.
It's interactive, so specify the dataset name and table name.

```
$ make create-bigquery-table
```

Set the expiration date of the partition set in table.
The expiration value is specified in ```.env``` as ```BIGQUERY_TABLE_PARTITIONING_EXPIRATION_TIME```.
It's interactive, so specify the dataset name and table name.

```
$ make set-bigquery-table-expiration-date
```

<br>

**1. Create GKE cluster, node group and IAM resources.**

```
$ make build
```

<br>

**2. Create kubernetes secrets.**

Collect the environment variables written in ```secrets.env``` and create them in a cluster as kubernetes secrets.

```
$ make secrets
```

<br>

**3. Deploy kubernetes resources.**

If you have set an Optional environment variable in .secrets.env, edit the container env in ./templates/stateless.yaml.
Set only the environment variables required for the container running on kubernetes.

Create kubernetes resources with helm.

Following command creates a StatefulSet, HeadlessService, Horizontal Pod Autoscaler, and PVC.

```
$ make deploy
```

The following processing is performed here.<br>
・build docker image.<br>
・push docker image to gcr repository.<br>
・build helm variables.<br>
・deploy with helm templates.<br>

<br>

**4. Upgrade kubernetes resources.**

You can upgrade kubernetes resources with the following command.

Note that if you do not update ```GCR_REPO_TAG```, the new docker image will be referenced and the container will not be created.

```
$ make upgrade
```

<br>

# Architects

![image](https://user-images.githubusercontent.com/37132477/141406547-41edf9eb-5a17-4191-9ee3-3f13ba17ec07.png)

A pod is created for each collection, and a persistent volume is linked to each pod.
Since the StatefulSet is created, even if the pod stops, you can get the change streams by referring to the resume token saved in the persistent volume again.