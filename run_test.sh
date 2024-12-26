#!/bin/bash
export RUN_ENV="dev"

go build  

nohup ./nostr_push >> /Users/allen/Projects/go_projecrs/nostr_push/nostr_push.log 2>&1&