#!/bin/sh

output_file=tmp/write_to_stderr.out

bin/tisp test/write_to_stderr.tisp 2> $output_file &&
diff $output_file test/write_to_stderr.out
