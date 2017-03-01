#!/bin/sh

[ $(bin/tisp test/multiple_outputs.tisp | wc -l | grep -o '[0-9]*' | head -1) \
  -eq 3 ]
