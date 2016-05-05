#!/bin/bash
find ./src -type d -name '__unit__' |
xargs -I {} find {} -name '*.js' |
xargs ./node_modules/mocha/bin/_mocha \
    --opts ./test/mocha.opts
