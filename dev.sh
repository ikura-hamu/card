#!/usr/bin/env bash

TUI_PANE=$(zellij action new-pane --name bubbletea)

zellij action write-chars --pane-id "$TUI_PANE" "go run ."
zellij action send-keys --pane-id "$TUI_PANE" "Enter"

watchexec -e go -- ./reload.sh "$TUI_PANE"
