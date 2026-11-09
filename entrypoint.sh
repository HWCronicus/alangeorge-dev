#!/bin/sh
set -e
mkdir -p .ssh

if [ -n "$SSH_HOST_KEY" ]; then
  printf '%s\n' "$SSH_HOST_KEY" > .ssh/id_ed25519
  chmod 600 .ssh/id_ed25519
elif [ ! -f .ssh/id_ed25519 ]; then
  ssh-keygen -t ed25519 -f .ssh/id_ed25519 -N ''
fi

./server &
./tui &
wait