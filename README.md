![image](https://user-images.githubusercontent.com/37132477/141405557-b7c9138e-0cf2-43ac-9f9e-b4cf78de60dc.png)

MxTransporter is a middleware that accurately carries change streams of MongoDB in real time. For infrastructure, you can easily use this middleware by creating a container image with Dockerfile on any platform and deploying it.

<br>

# Guide

## Build with samples
We have prepared a samples to build on AWS and GCP container orchestration services.
This can be easily constructed by setting environment variables as described and executing commands.

See ```docs/``` for details.

<br>

## Deploy to your container environment
With the Dockerfile provided, you can easily run MxTransporter by simply building the container image and deploying it to your favorite container environment.

### Requirement

- Build a Dockerfile, create an image, and create a container based on that image.

- Mount the persistent volume on the container to store the resume token. See the Change streams section of this README for more information.

- Allow access from the container to MongoDB

- Add permissions so that the container can access the export destination.

- Have the container read the required environment variables. All you need is a "Run locally" section in this README to add to your ```.env```.

<br>

## Run locally
### Requirement
- Set the following environment variables in ```.env```.

```
# Require 
## e.g. HOST=mongodb+srv://<user name>:<passward>@xxx.yyy.mongodb.net/test?retryWrites=true&w=majority
MONGODB_HOST=
MONGODB_DATABASE=
MONGODB_COLLECTION=

# Optional
## You have to specify this environment variable if you want to export BigQuery, Pub/Sub.
PROJECT_NAME_TO_EXPORT_CHANGE_STREAMS=

# Optional
## You have to specify this environment variable if you want to export BigQuery.
BIGQUERY_DATASET=
BIGQUERY_TABLE=

# Optional
## You have to specify this environment variable if you want to export Kinesis Data Stream.
KINESIS_STREAM_NAME=
KINESIS_STREAM_REGION=

# Require 
## One resume token is saved in the location specified here.
PERSISTENT_VOLUME_DIR=

# Require 
## Specify the location you want to export. 
## e.g. EXPORT_DESTINATION=bigquery 
EXPORT_DESTINATION=

# Require 
## Specify the time zone you run this middleware by referring to the following. (e.g. TIME_ZONE=Asia/Tokyo)
## https://cs.opensource.google/go/go/+/master:src/time/zoneinfo_abbrs_windows.go;drc=72ab424bc899735ec3c1e2bd3301897fc11872ba;l=15
TIME_ZONE=

# Optional
## MxT use zap library.
## Specify the MxT log setting you run this middleware by referring to the following.
### Specify log level, "0" is Info level, "1" is Error level. default is "0".
### e.g. LOG_LEVEL=0
LOG_LEVEL=
### Specify log format, json or console. default is json.
### e.g. LOG_FORMAT=console
LOG_FORMAT=
### Specify log output Directory.
### e.g. LOG_OUTPUT_PATH=/var/log/
### e.g. LOG_OUTPUT_PATH=../../var/log/
LOG_OUTPUT_DIRECTORY=
### Specify log output File.
### e.g. LOG_OUTPUT_FILE=mxt.log
LOG_OUTPUT_FILE=
```

- Allow access from the IP of the local machine on the mongoDB.


- Run ```go run ./cmd/main.go``` in the root directory.

<br>

# Architects

![image](https://user-images.githubusercontent.com/37132477/141405958-109351c4-fb47-4e3e-8146-4ecf055b0654.png)

1. MxTransporter watches MongoDB Collection.
2. When the Collection is updated, MxTransporter gets the change streams.
3. Format change streams according to the export destination.
4. Put change streams to the export destination.
5. If the put is successful, the resume Token included in the change streams is taken out and saved in the persistent volume.

<br>

# Specification

## MongoDB

### Connection to MongoDB
Allow the public IP of the MxTransporter container on the mongoDB side. This allows you to watch the changed streams that occur.

### Change streams
change streams output the change events that occurred in the database and are the same as the logs stored in oplog. And it has a unique token called resume token, which can be used to get events after a specific event.

In this system, resume token is saved in Persistent Volume associated with the container, and when a new container is started, the resume token is referenced and change streams acquisition starts from that point.

The resume token of the change streams just before the container stopped is stored in the persistent volume, so you can refer to it and get again the change streams that you missed while the container stopped and the new container started again.

The resume token is stored in the directory where the PVC is mounted.

```PERSISTENT_VOLUME_DIR``` is an environment variable given to the container.

```
{$PERSISTENT_VOLUME_DIR}/{year}/{month}/{day}
```

The resume token is saved in ```{year}-{month}-{day}.dat```.

```
$ pwd
{$PERSISTENT_VOLUME_DIR}/{year}/{month}/{day}

$ ls
{year}-{month}-{day}.dat

$ cat {year}-{month}-{day}.dat
T7466SLQD7J49BT7FQ4DYERM6BYGEMVD9ZFTGUFLTPFTVWS35FU4BHUUH57J3BR33UQSJJ8TMTK365V5JMG2WYXF93TYSA6BBW9ZERYX6HRHQWYS
```

When getting change-streams by referring to resumu token, it is designed to specify resume token in ```startAfrter``` of ```Collection.Watch()```.

<br>

## Export change streams
MxTransporter export change streams to the following description.

- Google Cloud BigQuery
- Google Cloud Pub/Sub
- Amazon Kinesis Data Streams

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
No special preparation is required. Create a Topic with the MongoDB Database name, and a Subscription with the MongoDB Collection name from which the change streams originated.

Change streams are sent to that subscription in a pipe (|) separated CSV.

### Kinesis Data Streams
No special preparation is required. If you want to separate the data warehouse table for each MongoDB collection for which you want to get change streams, use Kinesis Data Firehose and devise the output destination.

Change streams are sent to that in a pipe (|) separated CSV.

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

<br>

# Contributors
| [<img src="https://avatars.githubusercontent.com/KenFujimoto12" width="130px;"/><br />Kenshirou](https://github.com/KenFujimoto12) <br />   |
| :---: |
<br>


# Copyright

CAM, Inc. All rights reserved.