#!/usr/bin/env bash

PANE_ID="$1"

zellij action send-keys --pane-id "$PANE_ID" "Ctrl c"
sleep 0.1
zellij action write-chars --pane-id "$PANE_ID" "go run ."
zellij action send-keys --pane-id "$PANE_ID" "Enter"