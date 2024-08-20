#!/bin/bash
export RUN_ENV="pro"

CGO_LDFLAGS="-fuse-ld=lld" go build  
##go build

nohup ./nostr_push >> /Users/tangjinxing/Desktop/newProject/nostr_push/nostr_push.log 2>&1&