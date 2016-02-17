#!/bin/bash
echo -e $CERT_PEM > cert.pem
echo -e $CERT_KEY > cert.key
$@
