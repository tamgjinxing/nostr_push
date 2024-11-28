#!/bin/bash
export RUN_ENV="dev"

CGO_LDFLAGS="-fuse-ld=lld" go build  

nohup ./nostr_push >> /Users/allen/Projects/go_projecrs/nostr_push/nostr_push.log 2>&1&