#!/bin/sh
echo "Installing Flynn CLI"
L=/app/bin/flynn && curl -sL -A \"`uname -sp`\" https://dl.flynn.io/cli | zcat >$L && chmod +x $L
echo "Generating SSH Key..."
ssh-keygen -t rsa -N "" -f "/app/.ssh/id_rsa"
echo "Removing StrictHostKeyChecking for Flynn Host"
ssh -o StrictHostKeyChecking=no rkex.flynnhub.com -p 2222