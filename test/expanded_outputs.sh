#!/bin/sh

. test/lib/lib.sh

[ $(bin/tisp test/expanded_outputs.tisp | count_lines) -eq 5 ]
