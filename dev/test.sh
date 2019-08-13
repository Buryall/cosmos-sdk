#!/usr/bin/env bash

./killbyname.sh gaiad

rm -rf cache

gaiad testnet --v 4 --output-dir cache --chain-id testchain --base-port 10056 --starting-ip-address 127.0.0.1<<EOF
12345678
EOF
sleep 1

gaiad start --home cache/node0/gaiad --p2p.laddr tcp://127.0.0.1:10056 --p2p.seed_mode=true --log_level "main:info,state:info,*:error" --rpc.laddr tcp://127.0.0.1:10057 > cache/gaiad.0.log 2>&1 &

seed=$(gaiad tendermint show-node-id --home cache/node0/gaiad)

gaiad start --home cache/node1/gaiad --p2p.laddr tcp://127.0.0.1:10156 --p2p.seeds ${seed}@127.0.0.1:10056 --log_level "main:info,state:info,*:error" --rpc.laddr tcp://127.0.0.1:10157 > cache/gaiad.1.log 2>&1 &

gaiad start --home cache/node2/gaiad --p2p.laddr tcp://127.0.0.1:10256 --p2p.seeds ${seed}@127.0.0.1:10056 --log_level "main:info,state:info,*:error" --rpc.laddr tcp://127.0.0.1:10257 > cache/gaiad.2.log 2>&1 &

gaiad start --home cache/node3/gaiad --p2p.laddr tcp://127.0.0.1:10356 --p2p.seeds ${seed}@127.0.0.1:10056 --log_level "main:info,state:info,*:error" --rpc.laddr tcp://127.0.0.1:10357 > cache/gaiad.3.log 2>&1 &

sleep 1

for ((;;)) do
    tail -n 5 ./cache/gaiad.0.log
    tail -n 5 ./cache/gaiad.1.log
    tail -n 5 ./cache/gaiad.2.log
    tail -n 5 ./cache/gaiad.3.log
    sleep 5
done
