#!/bin/sh
L=/usr/local/bin/flynn && curl -sL -A \"`uname -sp`\" https://dl.flynn.io/cli | zcat >$L && chmod +x $L
ssh-keygen -t rsa -N "" -R wv5w.flynnhub.com -f "/usr/local/.ssh/id_rsa"