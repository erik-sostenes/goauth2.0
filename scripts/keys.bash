#!/bin/bash
openssl genpkey -algorithm RSA -out ./scripts/private_key.pem
openssl rsa -pubout -in ./scripts/private_key.pem -out ./scripts/public_key.pem