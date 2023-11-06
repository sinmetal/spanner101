# Spanner101

Spannerハンズオン資料

Spannerのテーブル設計をする上で重要で固有の機能である [インターリーブ](https://cloud.google.com/spanner/docs/schema-and-data-model?hl=en#parent-child) を中心に触ってみるハンズオン。

* pattern1 インターリーブなし 
* pattern2 インターリーブあり
* pattern3 おまけ

実行計画を見るために [spanner-cli](https://github.com/cloudspannerecosystem/spanner-cli) を利用している。

## Singers

公式ドキュメントでも出てくるSingers, Albumsを利用して、シンプルなQueryでインターリーブを体験する。

## DataBoost

Spanner Instanceに負荷をかけずにQueryを実行する [DataBoost](https://cloud.google.com/spanner/docs/databoost/databoost-overview) のハンズオン。

DataBoostで実行できるQueryは [PartitionQuery](https://cloud.google.com/spanner/docs/reads?hl=en#read_data_in_parallel) に限られます。
インターリーブの有無でPartitionQueryにできるQueryの変化を体験します。
