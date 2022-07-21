# MxTransporter を GKE 環境で構築する

ここでは MxTransporter を GKE 環境で構築するための手順を紹介します。
各 GCP リソースを作成するコマンドはすべて ```Makefike``` でラップされています。

## 準備
### コマンド
以下のコマンドが使用できる必要があります。記載のバージョンは検証済みです。

```
bq       v2.0.71
docker   v20.10.8
gcloud   v355.0.0
kubectl  v1.22.1
helm     v3.6.3
make     v3.81
```

### 環境変数
構築を始める前に、カレントディレクトリにある```.env.template``` と ```secrets.env.template``` を参考にして、```.env``` と ```.secrets.env```を作成してください。

もし、BigQuery と PubSub に Change Streams をエクスポートしたい場合は、```.secrets.env```に以下のような記載をしましょう。

```
EXPORT_DESTINATION=bigquery

or

EXPORT_DESTINATION=pubsub

or

EXPORT_DESTINATION=bigquery,pubsub
```

### Pubsub Ordering (オプション)
メッセージの順序指定を利用したい場合、環境変数```PUBSUB_ORDERING_BY```を設定する必要があります。
https://cloud.google.com/pubsub/docs/ordering

**注意**
メッセージの順序指定はパフォーマンスに悪影響をもたらす可能性があります。
参照: https://medium.com/google-cloud/google-cloud-pub-sub-ordered-delivery-1e4181f60bc8

### BigQuery スキーマ (オプション)
Change Streams を BigQuery にエクスポートしたい場合、以下のようなテーブルスキーマを指定する必要があります。

テーブルスキーマ
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

## 手順
**オプション: BigQuery セットアップ**

データセットとテーブルを作成する必要があります。

データセットを作成します。対話式なので、データセット名を指定します。

```
$ make create-bigquery-dataset
```

テーブルを作成します。対話式なので、テーブル名を指定します。

```
$ make create-bigquery-table
```

テーブルに設定されているパーティションの有効期限を設定します。有効期限の値は、```.env```で```BIGQUERY_TABLE_PARTITIONING_EXPIRATION_TIME```として指定されます。対話式なので、データセット名とテーブル名を指定します。

```
$ make set-bigquery-table-expiration-date
```

<br>

**1. GKE クラスター、ノードグループ、IAM リソース作成**

```
$ make build
```

<br>

**2. Kubernetes Secrets 作成**

```secrets.env```にある環境変数を収集し、それらを Kubernetes Secrets としてクラスターに作成します。

```
$ make secrets
```

<br>

**3. Kubernetes リソースをデプロイする**

もしオプションの環境変数を```.secrets.env```にセットしたのであれば、```./templates/stateless.yaml```内の env パラメータを編集します。Kubernetes で実行されているコンテナに必要な環境変数のみを設定します。

Kubernetes リソースは helm によって作成されます。

以下のコマンドは StatefulSet、HeadlessService、Horizontal Pod Autoscaler、PVC を作成します。

```
$ make deploy
```

このコマンドでは、以下のような処理が行われます。<br>
・GCR リポジトリへ Docker イメージを送信<br>
・helm の variables を作成します<br>
・helm テンプレートをデプロイする<br>

<br>

**4. Kubernetes リソースをアップグレードする**

以下のコマンドで、Kubernetes リソースをアップグレードできます。

環境変数の```GCR_REPO_TAG```を更新しないと、最新の Docker イメージが参照されず、新しいコンテナが作成されない点に注意してください。

```
$ make upgrade
```

<br>

# アーキテクチャ

![image](https://user-images.githubusercontent.com/37132477/141406547-41edf9eb-5a17-4191-9ee3-3f13ba17ec07.png)

MongoDB コレクションごとに Pod が作成され、それぞれの Pod に 永続ボリュームが紐付いています。StatefulSet が作成されるため、ポッドが停止した場合でも、永続ボリュームに保存されている再開トークンを再度参照することで、変更ストリームを取得できます。
