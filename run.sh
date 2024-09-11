#!/bin/bash
export RUN_ENV="pro"

go build  

nohup ./nostr_push >> /home/appl/push_nostr/nostr_push/nostr_push.log 2>&1&