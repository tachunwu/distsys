# Deterministic Snapshot
## API
提供和 Pebble 完全一樣的 API，只是變成分散式的，有以下 API
* Get, Set, Delete (return ErrNotFound if the requested key is not in the store)
* DeleteRange
* RangeKeySet
* RangeKeyUnset

## Single Partition
Publish Transaction 到 table.<partition_id> 然後監聽 service.<txn_id>，序列化處理完成就可以可用 VLL。

## Multi Partition 
當偵測到是 Multi Partition 強制送到 Multi Partition Stream 排序執行。