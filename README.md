![image](https://github.com/cam-inc/MxTransporter/blob/main/logo/mxt_logo.png)

MxTransporter is a middleware that accurately carries change streams of MongoDB in real time. For infrastructure, you can easily use this middleware by creating a container image with Dockerfile on any platform and deploying it.

With MxTransporter, real-time data can be reproduced and retained on the data utilization side, and data utilization will become even more active.

:jp: Japanese version of the README is [here](/README_JP.md).

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=cam-inc_MxTransporter&metric=coverage)](https://sonarcloud.io/summary/new_code?id=cam-inc_MxTransporter)
<br>

# Features
- Flexible export destination

It supports data warehouses, streaming services, etc. as destinations for Change Streams after they have been retrieved and formatted.

- Simultaneous multi-export destination

The formatted Change Streams information allows you to select multiple data warehouses, streaming services, etc. as destinations at the same time.

- Container base

MxTransporter creates a docker container image and can deploy it to your favorite container environment.
In addition, templates for easy construction in container orchestration services for AWS and GCP environments are also available in [/docs](/docs/README.md).

- No data loss

By utilizing the token required for data re-acquisition called "resume token" included in change streams, even if MxTransporter goes down, the data at the time of down can be re-acquired when it is restarted.


<br>

# Quick start

## Build in container orchestration services with samples
We have prepared a samples to build on AWS and GCP container orchestration services.
This can be easily constructed by setting environment variables as described and executing commands.

See [/docs](/docs/README.md) for details.

<br>

## Deploy to your container environment
With the Dockerfile provided, you can easily run MxTransporter by simply building the container image and deploying it to your favorite container environment.

### Requirement

- Execute the command ```make build-image``` in ```./Makefile```, Build a Dockerfile, create an image, and create a container based on that image.

- Mount the persistent volume on the container to store the resume token. See the change streams section of this README for more information.

- Allow access from the container to MongoDB

- Add permissions so that the container can access the export destination.

- Have the container read the required environment variables. All you need to do is pass the environment variables in ```.env.template``` to the container.

<br>

## Run locally
### Procedure
1. Create ```.env``` by referring to ```.env.template```.

2. Allow access from the IP of the local machine on the mongoDB.

3. Set permissions to access BigQuery, PubSub or Kinesis Data Streams from local.

For details on how to set it, refer to Registering AWS and GCP Credentials Locally.

4. Run

Run ```go run ./cmd/main.go``` in the root directory.

<br>

# Architects

![image](https://user-images.githubusercontent.com/37132477/141405958-109351c4-fb47-4e3e-8146-4ecf055b0654.png)

1. MxTransporter watches MongoDB Collection.
2. When the Collection is updated, MxTransporter gets the change streams.
3. Format change streams according to the export destination.
4. Put change streams to the export destination.
5. If the put is successful, the resume token included in the change streams is taken out and saved in the persistent volume.

<br>

# Specification

## MongoDB

### Connection to MongoDB
Allow the public IP of the MxTransporter container on the mongoDB side. This allows you to watch the changed streams that occur.

### Change Streams
Change streams output the change events that occurred in the database and are the same as the logs stored in oplog. And it has a unique token called resume token, which can be used to get events after a specific event.

In this system, resume token is saved in Persistent Volume associated with the container, and when a new container is started, the resume token is referenced and change streams acquisition starts from that point.

The resume token of the change streams just before the container stopped is stored in the persistent volume, so you can refer to it and get again the change streams that you missed while the container stopped and the new container started again.

The resume token can be set to the following as the save destination.

#### local file
The resume token is stored in the directory where the persistent volume is mounted.<br>
You can choose to save to a local file by setting ```RESUME_TOKEN_VOLUME_TYPE = file```.

```RESUME_TOKEN_VOLUME_DIR``` is an environment variable given to the container.

```
{$RESUME_TOKEN_VOLUME_DIR}/{$RESUME_TOKEN_FILE_NAME}.dat
```

The resume token is saved in a file called ``` {RESUME_TOKEN_FILE_NAME} .dat```. <br>
```RESUME_TOKEN_FILE_NAME``` is an optional environment variable, so if you don't set it, it will be saved in a file named ```{MONGODB_COLLECTION} .dat```.

```
$ pwd
{$PERSISTENT_VOLUME_DIR}

$ ls
{RESUME_TOKEN_FILE_NAME}.dat

$ cat {RESUME_TOKEN_FILE_NAME}.dat
T7466SLQD7J49BT7FQ4DYERM6BYGEMVD9ZFTGUFLTPFTVWS35FU4BHUUH57J3BR33UQSJJ8TMTK365V5JMG2WYXF93TYSA6BBW9ZERYX6HRHQWYS
```

#### External storage
It is also possible to save to cloud storage.

You can choose to save to S3 or GCS by setting ```RESUME_TOKEN_VOLUME_TYPE = s3 or RESUME_TOKEN_VOLUME_TYPE = gcs```. <br>
Set ``` RESUME_TOKEN_VOLUME_DIR``` as the object key.

Set the following environment variables as you like.

```
RESUME_TOKEN_VOLUME_BUCKET_NAME
RESUME_TOKEN_FILE_NAME
RESUME_TOKEN_BUCKET_REGION
RESUME_TOKEN_SAVE_INTERVAL_SEC
```

When getting change-streams by referring to resume token, it is designed to specify resume token in ```startAfrter``` of ```Collection.Watch()```.

<br>

## Export change streams
MxTransporter export change streams to the following description.

- Google Cloud BigQuery
- Google Cloud Pub/Sub
- Amazon Kinesis Data Streams
- Standard output

Set the environment variables as follows.
```
EXPORT_DESTINATION=bigquery

or

EXPORT_DESTINATION=kinesisStream

or

EXPORT_DESTINATION=pubsub

or

EXPORT_DESTINATION=file
```


### BigQuery
Create a BigQuery Table with a schema like the one below.

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

### Pub/Sub
Set the following environment variables to specify the topic name to which Change Streams will be exported.
```
PUBSUB_TOPIC_NAME
```

Change streams are sent to that subscription in a pipe (|) separated CSV.

### Kinesis Data Streams
No special preparation is required. If you want to separate the data warehouse table for each MongoDB collection for which you want to get change streams, use Kinesis Data Firehose and devise the output destination.

Change streams are sent to that in a pipe (|) separated CSV.

### Standard output
It is tandard output or file output.
This feature assumes the case of relaying data via a sidecar-powered agent (fluentd, fluentbit, etc.).
Due to the specifications of Docker log, if it is output as standard, it will be chunked at 16K Bytes, so if you want to avoid it, please use a file.
```
{"logType": "{FILE_EXPORTER_LOG_TYPE_KEY}","{FILE_EXPORTER_CHANGE_STREAM_KEY}":{// Change Stream Data //}}
```

<br>

## Format
Format before putting change streams to export destination. The format depends on the destination.

### BigQuery
Format to match the table schema and insert a value into each BigQuery Table field for each change streams.

### Pub/Sub
It is formatted into a pipe (|) separated CSV and put.

```
{"_data":"T7466SLQD7J49BT7FQ4DYERM6BYGEMVD9ZFTGUFLTPFTVWS35FU4BHUUH57J3BR33UQSJJ8TMTK365V5JMG2WYXF93TYSA6BBW9ZERYX6HRHQWYS
"}|insert|2021-10-01 23:59:59|{"_id":"6893253plm30db298659298h”,”name”:”xxx”}|{“coll”:”xxx”,”db”:”xxx”}|{“_id":"6893253plm30db298659298h"}|null
```

### Kinesis Data Streams
It is formatted into a pipe (|) separated CSV and put.

```
{"_data":"T7466SLQD7J49BT7FQ4DYERM6BYGEMVD9ZFTGUFLTPFTVWS35FU4BHUUH57J3BR33UQSJJ8TMTK365V5JMG2WYXF93TYSA6BBW9ZERYX6HRHQWYS
"}|insert|2021-10-01 23:59:59|{"_id":"6893253plm30db298659298h”,”name”:”xxx”}|{“coll”:”xxx”,”db”:”xxx”}|{“_id":"6893253plm30db298659298h"}|null
```

### Standard output
It is basic JSON. It is possible to change the key of ChangeStream, add a Time field by specifying the environment variable option.
```
{"logType": "{FILE_EXPORTER_LOG_TYPE_KEY}","{FILE_EXPORTER_CHANGE_STREAM_KEY}":{// Change Stream Data //},"{FILE_EXPORTER_TIME_KEY}":"2022-04-20T01:47:39.228Z"}
```

<br>

## Exclude Field
When exporting Change Streams, you can exclude unnecessary fields in the ```fullDocument```.

Set the following environment variables. If you specify more than one, separate them with ```,```.

Note that you can only exclude fields within the ``` fullDocument``` field.

```
MONGO_WATCH_PIPELINE_EXCLUDE_CS_FIELD

e.g.
MONGO_WATCH_PIPELINE_EXCLUDE_CS_FIELD=fulldocument.<Unnecessary fields>,fulldocument.<Unnecessary fields>
```

Please refer to the following MongoDB documentation.
https://www.mongodb.com/docs/manual/reference/operator/aggregation/unset/#mongodb-pipeline-pipe.-unset

<br>

# Contributors
| [<img src="https://avatars.githubusercontent.com/KenFujimoto12" width="130px;"/><br />Kenshirou](https://github.com/KenFujimoto12) <br />   | [<img src="https://avatars.githubusercontent.com/syama666" width="130px;"/><br />Yoshinori Sugiyama](https://github.com/syama666) <br />   |
| :---: | :---: |
<br>


# Copyright

CAM, Inc. All rights reserved.
