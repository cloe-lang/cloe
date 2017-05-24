#!/bin/sh

bin/tisp test/infinite_recursion.tisp > /dev/null &
pid=$!

sleep 1 # Wait for memory usage to be stable.

ok=false
last_mem=0

for _ in $(seq 10)
do
  mem=$(ps ho vsize $pid | tail -1)

  if [ $last_mem -ge $mem ]
  then
    ok=true
    break
  fi

  last_mem=$mem
  sleep 1
done &&

kill $pid &&
$ok
