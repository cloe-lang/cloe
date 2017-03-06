#!/bin/sh

. test/lib/lib.sh

[ $(bin/tisp test/multiple_outputs.tisp | count_lines) -eq 3 ]
