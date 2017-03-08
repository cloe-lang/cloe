#!/bin/sh

[ -n "$(bin/tisp test/print_call_trace.tisp 2>&1 > /dev/null)" ]
