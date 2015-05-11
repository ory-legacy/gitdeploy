#!/bin/sh
L=/app/bin/flynn && curl -sL -A \"`uname -sp`\" https://dl.flynn.io/cli | zcat >$L && chmod +x $L
ssh-keygen -t rsa -N "" -f "/app/.ssh/id_rsa"