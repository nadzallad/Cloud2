@echo off

set DOCKER_USER=nadzalla
set TAG=latest

echo ======================
echo BUILD + PUSH ALL IMAGE
echo ======================

docker build -t %DOCKER_USER%/payment-service:%TAG% PaymentService
docker push %DOCKER_USER%/payment-service:%TAG%

echo DONE