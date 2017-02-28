#!/bin/sh

bin/tisp test/infinite_recursion.tisp &
pid=$!

sleep 1 # Wait for memory usage to be stable.

ok=false
last_mem=0

for i in $(seq 10)
do
  mem=$(ps ho vsize $pid)

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
