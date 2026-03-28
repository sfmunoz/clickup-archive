#!/bin/bash

set -e -o pipefail

cd "$(dirname "$0")"

[ "$HUGO_THEME" = "" ] && HUGO_THEME="picocss"
export HUGO_THEME

set -x

# cleanup: hugo doesn't do it
rm -rf public
hugo server -D --disableFastRender --noHTTPCache
