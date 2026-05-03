@echo off

set DOCKER_USER=nadzalla
set TAG=latest

echo ======================
echo BUILD + PUSH ALL IMAGE
echo ======================

docker build -t %DOCKER_USER%/payment-service:%TAG% PaymentService
docker push %DOCKER_USER%/payment-service:%TAG%

docker build -t %DOCKER_USER%/order-service:%TAG% OrderService
docker push %DOCKER_USER%/order-service:%TAG%

docker build -t %DOCKER_USER%/pickup-service:%TAG% PickupService
docker push %DOCKER_USER%/pickup-service:%TAG%

docker build -t %DOCKER_USER%/warehouse-service:%TAG% WarehouseService
docker push %DOCKER_USER%/warehouse-service:%TAG%

docker build -t %DOCKER_USER%/shipment-service:%TAG% ShipmentService
docker push %DOCKER_USER%/shipment-service:%TAG%

docker build -t %DOCKER_USER%/delivery-service:%TAG% DeliveryService
docker push %DOCKER_USER%/delivery-service:%TAG%

docker build -t %DOCKER_USER%/notification-service:%TAG% NotificationService
docker push %DOCKER_USER%/notification-service:%TAG%

docker build -t %DOCKER_USER%/tracking-service:%TAG% TrackingService
docker push %DOCKER_USER%/tracking-service:%TAG%

echo DONE