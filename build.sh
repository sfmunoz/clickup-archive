#!/bin/bash

set -e -o pipefail

cd "$(dirname "$0")"

[ "$HUGO_THEME" = "" ] && HUGO_THEME="picocss"
export HUGO_THEME

MINIFY_FLAG="--minify"
[ "$MINIFY" = "0" ] && MINIFY_FLAG=""

set -x

# cleanup: hugo doesn't do it
rm -rf public
hugo build --gc $MINIFY_FLAG --panicOnWarning "$@"
