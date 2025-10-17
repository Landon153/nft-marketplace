#!/usr/bin/env bash

APP_NAME=nft-marketplace-test-web

docker load -i ${APP_NAME}.tar
docker rm -f ${APP_NAME}

docker run -d --name ${APP_NAME} \
  --network my-net \
  -e VIRTUAL_HOST="test.nft.marketplace.200lab.io" \
  -e VIRTUAL_PORT=3000 \
  -e PORT=3000 \
  -e LETSENCRYPT_HOST="test.nft.marketplace.200lab.io" \
  -e LETSENCRYPT_EMAIL="blockchain.test@test.nft.marketplace.200lab.io" \
  -e ENABLE_IPV6=true \
  ${APP_NAME}

exit
