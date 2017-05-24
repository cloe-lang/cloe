#!/bin/sh

bin/tisp-parse < test/source_from_stdin.tisp > /dev/null &&
bin/tisp < test/source_from_stdin.tisp > /dev/null
