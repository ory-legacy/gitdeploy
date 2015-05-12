#!/bin/sh
L=/app/bin/flynn && curl -sL -A \"`uname -sp`\" https://dl.flynn.io/cli | zcat >$L && chmod +x $L
ssh-keygen -t rsa -N "" -f "/app/.ssh/id_rsa"
ssh -o StrictHostKeyChecking=no p3mj.flynnhub.com -p 2222
flynn key add