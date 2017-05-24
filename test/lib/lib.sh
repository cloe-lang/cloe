#!/bin/sh

count_lines() {
  wc -l | grep -o '[0-9]\+' | head -1
}
