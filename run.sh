#!/bin/bash
#
# Date:   Sat Aug 30 09:12:24 UTC 2025
# Issue:  https://github.com/sfmunoz/cms/issues/11
#

set -e

cd "$(dirname "$0")"

[ "$HUGO_THEME" = "" ] && HUGO_THEME="headless"
export HUGO_THEME

set -x

# cleanup: hugo doesn't do it
rm -rf public
hugo server -D --disableFastRender --noHTTPCache
