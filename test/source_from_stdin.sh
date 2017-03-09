#!/bin/sh

cat test/read_stdin.tisp | bin/tisp-parse > /dev/null &&
cat test/read_stdin.tisp | bin/tisp > /dev/null
