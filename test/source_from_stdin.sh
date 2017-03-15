#!/bin/sh

cat test/source_from_stdin.tisp | bin/tisp-parse > /dev/null &&
cat test/source_from_stdin.tisp | bin/tisp > /dev/null
