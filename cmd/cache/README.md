# Cache
## Motivation
希望用 Stream 的技術用 Go 實作和 mcrouter, memcached 類似的 Semantic，目標如下。
* key/value operation: 提供單一的 Key/Value 介面
* replicated pool: 寫入時能夠寫入多個 Instance 來達到 HA
* fault tolerant: 當 cache miss 不論何種原因能夠執行相對應的處理
* hash partition: 針對 Key 本身做 Partition 
* subject partition: 針對特定 Key Prefix 做 partition

## Architecture
主要為 Two-Level 架構，第一層為 Proxy，負責統一 routing request，第二層則是 Cache Instance，只作為單體的 cache node。

### Layer 1: Proxy
Client 發起同步呼叫將要處理的 Operation(Get/Set/Delete)發佈到 NATS Proxy，js.Publish() 本身是同步的非常的慢(這個設定等同於Async，但是 MaxPending = 1)。

直接提供 key-level topic **(例如：table.partition.key)**，這樣 API 本身提供為以下
* Single-key transaction (GET/SET/DELETE)
* Multi-key transaction
* Multi-parition transaction

### Layer 2: Cache Instance

### Implementation Details
#### Get
向 Proxy 發起 Request-Reply，Proxy 向正在監聽這個 Topic 的 Cache Instance 拿資料。

#### Set
原則上可以用原生的 Pub/Sub，但是為了後續更好支援 Distributed Transaction 用 JetStream Push-based Consumer 實作。

採取 facebook 的 aside policy 先 invalid 再寫入。

#### Del


# Cache Instance
自己去和 JetStream 申請 Consumer。Write 的時候用廣播的，Read 的時候有一個回覆即可。
* 註冊自己的 shard consumer
* 註冊 stream
* 訂閱自己的 CACHE.<shard_id>


# Experiments
SyncSubscribe  10,0000 took 5.232976505s

## Experiments: Single shard Multi consumer
30000 Txn, 10000 GET, 10000 SET, 10000 DEL
1 client:
2 client:
4 client:
8 client:


## Reference
* (Introducing mcrouter: A memcached protocol router for scaling memcached deployments)[https://engineering.fb.com/2014/09/15/web/introducing-mcrouter-a-memcached-protocol-router-for-scaling-memcached-deployments/]
* (Learn NATS by Example)[https://natsbyexample.com/]