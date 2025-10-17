#!/usr/bin/env bash

APP_NAME=nft-marketplace-test-web
DEPLOY_CONNECT=$1

if [[ -z "$DEPLOY_CONNECT" ]]; then
	echo 'deploy host cannot be empty.'
	exit 1
fi

echo "Docker building..."
docker build -t ${APP_NAME} .
echo "Docker saving..."
docker save -o ${APP_NAME}.tar ${APP_NAME}

echo "Deploying..."
scp -o StrictHostKeyChecking=no ./${APP_NAME}.tar  ${DEPLOY_CONNECT}:~
ssh -o StrictHostKeyChecking=no ${DEPLOY_CONNECT} 'bash -s' < ./deploy/dev.sh

echo "Cleaning..."
rm -f ./${APP_NAME}.tar
echo "Done"
