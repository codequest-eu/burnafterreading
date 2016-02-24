#!/bin/sh

trap 'kill -9 $(jobs -p)' EXIT
echo -e $CERT_PEM > cert.pem
echo -e $CERT_KEY > cert.key
$@
