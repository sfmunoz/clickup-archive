#!/bin/bash
#
# Date:   Sat Aug 30 09:40:48 AM UTC 2025
# Issue:  https://github.com/sfmunoz/cms/issues/13
#

set -e

cd "$(dirname "$0")"

[ "$HUGO_THEME" = "" ] && HUGO_THEME="headless"
export HUGO_THEME

MINIFY_FLAG="--minify"
[ "$MINIFY" = "0" ] && MINIFY_FLAG=""

set -x

# cleanup: hugo doesn't do it
rm -rf public
hugo build --gc $MINIFY_FLAG --panicOnWarning "$@"

