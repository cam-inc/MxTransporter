![image](https://github.com/cam-inc/MxTransporter/blob/main/logo/mxt_logo.png)

MxTransporter は MongoDB の Change Streams を正確に、リアルタイムで送信先に運ぶミドルウェアです。MxTransporter は、用意された Dockerfile からコンテナイメージを作成し、お好きな環境へデプロイすることで簡単に利用することができます。

また、[GitHub Container Registry](https://github.com/cam-inc/MxTransporter/pkgs/container/mxtransporter) 上にコンテナイメージを公開しているため、ビルドすることなく、コンテナイメージを取得して利用することも可能です。

MxTransporterにより、データ活用側でリアルタイムなデータを再現、保持でき、データ活用がより一層活発になるでしょう。

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=cam-inc_MxTransporter&metric=coverage)](https://sonarcloud.io/summary/new_code?id=cam-inc_MxTransporter)
<br>


# 特徴
- 柔軟に選べるエクスポート先

取得し、整形された後の Change Streams のエクスポート先としてデータウェアハウスとストリーミングサービスなどをサポートしています。

- 複数宛先への同時エクスポート

整形された Change Streams の情報は、複数のデータウェアハウス、ストリーミングサービスなどを、宛先として同時に選択することができます。

- コンテナ基盤

MxTransporter は用意された Dockerfile をビルドし、作成されたコンテナイメージをお好きな環境にデプロイするだけで使用できます。さらに、AWS や GCP のコンテナオーケストレーションサービス上で構築するためのサンプルが [/docs](/docs/README_JP.md) に用意されているので、それを元に簡単に構築もできます。

- ロスレス転送

Change Streams に含まれている resume token はデータ再取得のために役立ちます。それを活用することで、MxTransporter が停止した際、復旧後に停止時間分の Change Streams を再取得することができます。

<br>

# クイックスタート

## サンプルを用いて、コンテナオーケストレーションサービス上に構築する
AWS や GCP のコンテナオーケストレーションで構築するためのサンプルを用意しています。環境変数をセットしてコマンドを実行することで簡単に構築できます。

詳しくは[/docs](/docs/README_JP.md) をご覧ください。

<br>

## 好きなコンテナ環境にデプロイする
用意された Dockerfile を元にコンテナイメージを作成し、好きなコンテナ実行環境にデプロイすることで、簡単に MxTransporter を実行することができます。

### 必要事項

- ```./Makefile```にあるコマンド```make build-image```を実行し、Dockerfile をビルドし、コンテナイメージを作成します。それを元に好きな環境でコンテナを作成します。

- resume token を保存するためにコンテナに永続ボリュームをマウントします。こちらに関して詳しくは、この README の **Change Streams** セクションを御覧ください。

- コンテナから MongoDB に対してアクセスを許可します。

- コンテナ環境にエクスポート先のデータウェアハウスやストリーミングサービスへアクセスするための権限を与えます。

- コンテナに必要な環境変数を読み取らせます。```.env.template```に必要な環境変数があるのでそれを渡します。

<br>

## ローカルで実行
### 手順
1. ```.env.template```を参考に```.env```を作成する。

2. MongoDB にローカルマシンのIPからのアクセスを許可します。

3. BigQuery、PubSub、Kinesis Data Streams に対して、ローカルマシンからのアクセスを許可します。

設定方法の詳細について、AWS と GCP のドキュメントを参照してください。

4. 実行

本リポジトリのルートディレクトリで```go run ./cmd/main.go```を実行します。

<br>

# アーキテクチャ

![image](https://user-images.githubusercontent.com/37132477/141405958-109351c4-fb47-4e3e-8146-4ecf055b0654.png)

1. MxTransporter が MongoDB のコレクションを参照しています。
2. コレクションに更新があると、MxTransporter が Change Streamsを取得します。
3. エクスポート先に合うように Change Streams のフォーマットを整形します。
4. エクスポート先に Change Streams を送ります。
5. 送信が成功したら、 Change Streams に含まれている resume token を永続ボリュームに保存します。

<br>

# 仕様

## MongoDB

### MongoDBへの接続
MongoDB 側に MxTransporter のコンテナの Public IP からのアクセスを許可します。これにより、Change Streams を取得することができます。

### Change Streams
Change Streams はデータベースで起きた変更イベントを出力し、ログとして oplog に保存されます。resume token と呼ばれるユニークなトークンを含んでおり、それにより特定のイベントの Change Streams を取得することができます。

本システムではコンテナに紐付いた永続ボリュームに resume token が保存され、新しいコンテナがスタートしたときに、resume token が参照され、その時点から Change Streams の取得が開始されます。

コンテナが停止する直前の Change Streams の resume token は永続ボリュームに保存されるため、コンテナが停止し、新しいコンテナがスタートしたときに resume token を参照して、逃した Change Streams を再取得できます。

resume token は以下を保存先として設定できます。

#### ローカルファイル
resume token は、永続ボリュームがマウントされているディレクトリに保存されます。<br>
```RESUME_TOKEN_VOLUME_TYPE=file```とすることで、ローカルファイルへの保存を選択できます。

```RESUME_TOKEN_VOLUME_DIR``` という環境変数をコンテナに与えます。

```
{$RESUME_TOKEN_VOLUME_DIR}/{$RESUME_TOKEN_FILE_NAME}.dat
```

resume token は```{RESUME_TOKEN_FILE_NAME}.dat```というファイルに保存されます。<br>
```RESUME_TOKEN_FILE_NAME```はオプショナルな環境変数なので、設定しない場合は```{MONGODB_COLLECTION}.dat```という名前のファイルに保存されます

```
$ pwd
{$PERSISTENT_VOLUME_DIR}

$ ls
{RESUME_TOKEN_FILE_NAME}.dat

$ cat {RESUME_TOKEN_FILE_NAME}.dat
T7466SLQD7J49BT7FQ4DYERM6BYGEMVD9ZFTGUFLTPFTVWS35FU4BHUUH57J3BR33UQSJJ8TMTK365V5JMG2WYXF93TYSA6BBW9ZERYX6HRHQWYS
```

#### 外部ストレージ
クラウドストレージに保存することも可能です。

```RESUME_TOKEN_VOLUME_TYPE=s3 or RESUME_TOKEN_VOLUME_TYPE=gcs```とすることで、S3かGCSへの保存を選択できます。<br>
オブジェクトキーとしては```RESUME_TOKEN_VOLUME_DIR```を設定してください。

以下の環境変数を任意で設定してください。
```
RESUME_TOKEN_VOLUME_BUCKET_NAME
RESUME_TOKEN_FILE_NAME
RESUME_TOKEN_BUCKET_REGION
RESUME_TOKEN_SAVE_INTERVAL_SEC
```

resume token を参照して Change Streams を取得する場合、```Collection.Watch()```の```startAfrter```で resume tokenを指定するように設計されています。

<br>

## Change Streams をエクスポートする
MxTransporter は以下の宛先に Change Streams をエクスポートします。

- Google Cloud BigQuery
- Google Cloud Pub/Sub
- Amazon Kinesis Data Streams
- Standard output
- Local file

以下のように環境変数を設定します。
```
EXPORT_DESTINATION=bigquery

or

EXPORT_DESTINATION=kinesisStream

or

EXPORT_DESTINATION=pubsub

or

EXPORT_DESTINATION=stdout

or

EXPORT_DESTINATION=file
```

### BigQuery
次のようなスキーマで BigQuery テーブルを作成します。

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
以下の環境変数を設定し、Change Streamsをエクスポートするトピック名を指定します。
```
PUBSUB_TOPIC_NAME
```

パイプ(|)区切りのCSV形式で Change Streams はサブスクリプションに送られます。

### Kinesis Data Streams
特段、準備は必要有りません。Change Streams を取得する MongoDB ごとにデータウェアハウステーブルを分離する場合は、Kinesis Data Firehose を使用して、出力先を指定します。

パイプ(|)区切りのCSV形式で Change Streams はサブスクリプションに送られます。

### Standard output or File
標準出力またはファイル出力します。
この機能はサイドカーで動くエージェント(fluentd, fluentbit 等)経由でデータをリレーするケースを想定しています。
Dockerログの仕様上、標準出力した場合は16K Bytesでチャンクされるので、それを避ける場合はファイルを使用してください。
```
{"logType": "{FILE_EXPORTER_LOG_TYPE_KEY}","{FILE_EXPORTER_CHANGE_STREAM_KEY}":{// Change Stream Data //}}
```

<br>

## フォーマット
Change Streams をエクスポート先に送る前にフォーマットを整えます。形式はエクスポート先によって異なります。

### BigQuery
テーブルスキーマに合うように、Change Streams の各値がテーブルフィールドに入ります。

### Pub/Sub
パイプ(|)で区切られたCSV形式にフォーマットが整えられます。

```
{"_data":"T7466SLQD7J49BT7FQ4DYERM6BYGEMVD9ZFTGUFLTPFTVWS35FU4BHUUH57J3BR33UQSJJ8TMTK365V5JMG2WYXF93TYSA6BBW9ZERYX6HRHQWYS
"}|insert|2021-10-01 23:59:59|{"_id":"6893253plm30db298659298h”,”name”:”xxx”}|{“coll”:”xxx”,”db”:”xxx”}|{“_id":"6893253plm30db298659298h"}|null
```

### Kinesis Data Streams
パイプ(|)で区切られたCSV形式にフォーマットが整えられます。

```
{"_data":"T7466SLQD7J49BT7FQ4DYERM6BYGEMVD9ZFTGUFLTPFTVWS35FU4BHUUH57J3BR33UQSJJ8TMTK365V5JMG2WYXF93TYSA6BBW9ZERYX6HRHQWYS
"}|insert|2021-10-01 23:59:59|{"_id":"6893253plm30db298659298h”,”name”:”xxx”}|{“coll”:”xxx”,”db”:”xxx”}|{“_id":"6893253plm30db298659298h"}|null
```

### Standard output or File
基本的なJSONです。環境変数オプション指定によりChangeStreamのキーを変更したり、Timeフィールドを追加することが可能です。
```
{"logType": "{FILE_EXPORTER_LOG_TYPE_KEY}","{FILE_EXPORTER_CHANGE_STREAM_KEY}":{// Change Stream Data //},"{FILE_EXPORTER_TIME_KEY}":"2022-04-20T01:47:39.228Z"}
```

<br>

# Contributors
| [<img src="https://avatars.githubusercontent.com/KenFujimoto12" width="130px;"/><br />Kenshirou](https://github.com/KenFujimoto12) <br /> | [<img src="https://avatars.githubusercontent.com/syama666" width="130px;"/><br />Yoshinori Sugiyama](https://github.com/syama666) <br /> |
| :---------------------------------------------------------------------------------------------------------------------------------------: | :--------------------------------------------------------------------------------------------------------------------------------------: |
<br>


# Copyright

CAM, Inc. All rights reserved.

