pipeline {
    agent any

    environment {
        PAYMENT_IMAGE = "nadzalla/payment-service:${env.BUILD_NUMBER}"
        ORDER_IMAGE = "nadzalla/order-service:${env.BUILD_NUMBER}"
        SHIPMENT_IMAGE = "nadzalla/shipment-service:${env.BUILD_NUMBER}"
        DELIVERY_IMAGE = "nadzalla/delivery-service:${env.BUILD_NUMBER}"
        PICKUP_IMAGE = "nadzalla/pickup-service:${env.BUILD_NUMBER}"
        WAREHOUSE_IMAGE = "nadzalla/warehouse-service:${env.BUILD_NUMBER}"
        TRACKING_IMAGE = "nadzalla/tracking-service:${env.BUILD_NUMBER}"
        NOTIFICATION_IMAGE = "nadzalla/notification-service:${env.BUILD_NUMBER}"

        NETWORK = "microservices-net"
    }

    stages {

        // =====================================================
        // CHECKOUT
        // =====================================================
        stage('Checkout Repo') {
            steps {
                deleteDir()
                git branch: 'main', url: 'https://github.com/nadzallad/Cloud2.git'
            }
        }

        // =====================================================
        // UNIT TEST
        // =====================================================
        stage('Unit Test') {
            steps {

                dir('PaymentService') {
                    catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                        sh 'go test -v -run TestValidatePayment ./...'
                    }
                }

                dir('OrderService') {
                    catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                        sh 'go test -short ./...'
                    }
                }

                dir('DeliveryService') {
                    catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                        sh 'go test ./...'
                    }
                }

                dir('ShipmentService') {
                    catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                        sh 'go test ./...'
                    }
                }

                dir('PickupService') {
                    catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                        sh 'go test ./...'
                    }
                }

                dir('WarehouseService') {
                    catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                        sh 'go test ./...'
                    }
                }

                dir('TrackingService') {
                    catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                        sh 'go test -short ./...'
                    }
                }

                dir('NotificationService') {
                    catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                        sh '''
                            go test -short ./... \
                            -run "TestValidateNotification"
                        '''
                    }
                }
            }
        }

        // =====================================================
        // LINT / VET
        // =====================================================
        stage('Lint / Vet') {
            steps {

                dir('PaymentService') {
                    sh 'go vet ./...'
                }

                dir('OrderService') {
                    sh 'go vet ./...'
                }

                dir('DeliveryService') {
                    sh 'go vet ./...'
                }

                dir('ShipmentService') {
                    sh 'go vet ./...'
                }

                dir('PickupService') {
                    sh 'go vet ./...'
                }

                dir('WarehouseService') {
                    sh 'go vet ./...'
                }

                dir('TrackingService') {
                    sh 'go vet ./...'
                }

                dir('NotificationService') {
                    sh 'go vet ./...'
                }
            }
        }

        // =====================================================
        // BUILD IMAGE
        // =====================================================
        stage('Build Image') {
            steps {
                sh '''
                docker build -t $PAYMENT_IMAGE ./PaymentService
                docker build -t $ORDER_IMAGE ./OrderService
                docker build -t $SHIPMENT_IMAGE ./ShipmentService
                docker build -t $DELIVERY_IMAGE ./DeliveryService
                docker build -t $PICKUP_IMAGE ./PickupService
                docker build -t $WAREHOUSE_IMAGE ./WarehouseService
                docker build -t $TRACKING_IMAGE ./TrackingService
                docker build -t $NOTIFICATION_IMAGE ./NotificationService
                '''
            }
        }

        // =====================================================
        // FUNCTIONAL TEST
        // =====================================================
        stage('Functional Test') {
            steps {
                catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {

                    sh '''
                    echo "START FUNCTIONAL TEST"

                    # ================= CLEANUP =================
                    docker rm -f mongo-test || true

                    docker rm -f \
                    test-payment \
                    test-order \
                    test-delivery \
                    test-shipment \
                    test-pickup \
                    test-warehouse \
                    test-tracking \
                    test-notification || true

                    docker network rm $NETWORK || true

                    # ================= CREATE NETWORK =================
                    docker network create $NETWORK

                    docker network connect $NETWORK jenkins-server || true

                    # ================= MONGODB =================
                    docker run -d \
                      --name mongo-test \
                      --network $NETWORK \
                      -e MONGO_INITDB_ROOT_USERNAME=admin \
                      -e MONGO_INITDB_ROOT_PASSWORD=admin123 \
                      mongo

                    echo "WAIT MONGO"
                    sleep 15

                    # ================= PAYMENT =================
                    docker run -d --name test-payment \
                      --network $NETWORK \
                      -e DB_HOST=host.docker.internal \
                      -e DB_NAME=payment_db \
                      -e DB_PASS=admin123 \
                      -p 8082:8082 \
                      $PAYMENT_IMAGE

                    # ================= ORDER =================
                    docker run -d --name test-order \
                      --network $NETWORK \
                      -p 8081:8081 \
                      $ORDER_IMAGE

                    # ================= DELIVERY =================
                    docker run -d --name test-delivery \
                      --network $NETWORK \
                      -e DB_HOST=host.docker.internal \
                      -e DB_NAME=delivery_db \
                      -e DB_USER=postgres \
                      -e DB_PASSWORD=admin123 \
                      -p 8086:8086 \
                      $DELIVERY_IMAGE

                    # ================= SHIPMENT =================
                    docker run -d --name test-shipment \
                      --network $NETWORK \
                      -e DB_HOST=host.docker.internal \
                      -e DB_NAME=shipment_db \
                      -e DB_USER=postgres \
                      -e DB_PASSWORD=admin123 \
                      -p 8085:8085 \
                      $SHIPMENT_IMAGE

                    # ================= PICKUP =================
                    docker run -d --name test-pickup \
                      --network $NETWORK \
                      -p 8083:8083 \
                      $PICKUP_IMAGE

                    # ================= WAREHOUSE =================
                    docker run -d --name test-warehouse \
                      --network $NETWORK \
                      -p 8084:8084 \
                      $WAREHOUSE_IMAGE

                    # ================= TRACKING =================
                    docker run -d \
                      --name test-tracking \
                      --network $NETWORK \
                      -e MONGO_URI="mongodb://admin:admin123@mongo-test:27017/?authSource=admin" \
                      -p 8087:8087 \
                      $TRACKING_IMAGE

                    # ================= NOTIFICATION =================
                    docker run -d \
                      --name test-notification \
                      --network $NETWORK \
                      -e MONGO_URI="mongodb://admin:admin123@mongo-test:27017/?authSource=admin" \
                      -p 8088:8088 \
                      $NOTIFICATION_IMAGE

                    echo "WAIT APPLICATION"
                    sleep 20

                    echo "CHECK CONTAINER"
                    docker ps -a

                    # ================= API TEST =================

                    curl -s -X POST http://host.docker.internal:8082/payment \
                    -H "Content-Type: application/json" \
                    -d '{"amount":1,"paid":1}'

                    curl -s -X POST http://host.docker.internal:8081/order \
                    -H "Content-Type: application/json" \
                    -d '{"user_id":1,"weight_kg":2,"distance_km":5,"base_price":10000}'

                    curl -s -X POST http://host.docker.internal:8086/delivery \
                    -H "Content-Type: application/json" \
                    -d '{"tracking_number":"LOG-0-1779347830","address":"Bandung"}'

                    curl -s -X POST http://host.docker.internal:8085/shipment \
                    -H "Content-Type: application/json" \
                    -d '{"tracking_number":"LOG-0-1779347830"}'

                    curl -s -X POST http://host.docker.internal:8083/pickup \
                    -H "Content-Type: application/json" \
                    -d '{"order_id":"ORD1","payment_status":"paid","weight":2}'

                    curl -s http://host.docker.internal:8084/warehouse

                    # ================= TRACKING TEST =================
                    cd TrackingService
                    go test -run TestTrackingAPI_Success
                    cd ..

                    # ================= NOTIFICATION TEST =================
                    cd NotificationService
                    go test -run TestNotificationAPI
                    cd ..
                    '''
                }
            }
        }

        // =====================================================
        // PUSH IMAGE
        // =====================================================
        stage('Push Image') {
            steps {
                withCredentials([usernamePassword(
                    credentialsId: 'logistic-login',
                    usernameVariable: 'USERNAME',
                    passwordVariable: 'PASSWORD'
                )]) {

                    sh '''
                    echo "$PASSWORD" | docker login -u "$USERNAME" --password-stdin

                    docker push $PAYMENT_IMAGE
                    docker push $ORDER_IMAGE
                    docker push $DELIVERY_IMAGE
                    docker push $SHIPMENT_IMAGE
                    docker push $PICKUP_IMAGE
                    docker push $WAREHOUSE_IMAGE
                    docker push $TRACKING_IMAGE
                    docker push $NOTIFICATION_IMAGE
                    '''
                }
            }
        }

        // =====================================================
        // DEPLOY
        // =====================================================
        stage('Deploy') {
            steps {
                sh '''
                echo "START FULL DEPLOYMENT"

                # ================= CLEAN OLD =================
                docker rm -f mongo-prod || true

                docker rm -f \
                prod-payment \
                prod-order \
                prod-delivery \
                prod-shipment \
                prod-pickup \
                prod-warehouse \
                prod-tracking \
                prod-notification || true

                docker network rm $NETWORK || true
                docker network create $NETWORK

                # ================= MONGODB =================
                docker run -d \
                  --name mongo-prod \
                  --network $NETWORK \
                  -e MONGO_INITDB_ROOT_USERNAME=admin \
                  -e MONGO_INITDB_ROOT_PASSWORD=admin123 \
                  mongo

                sleep 15

                # ================= PAYMENT =================
                docker run -d --name prod-payment \
                  --network $NETWORK \
                  -e DB_HOST=host.docker.internal \
                  -e DB_NAME=payment_db \
                  -e DB_PASS=admin123 \
                  -p 8082:8082 \
                  $PAYMENT_IMAGE

                # ================= ORDER =================
                docker run -d --name prod-order \
                  --network $NETWORK \
                  -p 8081:8081 \
                  $ORDER_IMAGE

                # ================= DELIVERY =================
                docker run -d --name prod-delivery \
                  --network $NETWORK \
                  -e DB_HOST=host.docker.internal \
                  -e DB_NAME=delivery_db \
                  -e DB_USER=postgres \
                  -e DB_PASSWORD=admin123 \
                  -p 8086:8086 \
                  $DELIVERY_IMAGE

                # ================= SHIPMENT =================
                docker run -d --name prod-shipment \
                  --network $NETWORK \
                  -e DB_HOST=host.docker.internal \
                  -e DB_NAME=shipment_db \
                  -e DB_USER=postgres \
                  -e DB_PASSWORD=admin123 \
                  -p 8085:8085 \
                  $SHIPMENT_IMAGE

                # ================= PICKUP =================
                docker run -d --name prod-pickup \
                  --network $NETWORK \
                  -p 8083:8083 \
                  $PICKUP_IMAGE

                # ================= WAREHOUSE =================
                docker run -d --name prod-warehouse \
                  --network $NETWORK \
                  -p 8084:8084 \
                  $WAREHOUSE_IMAGE

                # ================= TRACKING =================
                docker run -d \
                  --name prod-tracking \
                  --network $NETWORK \
                  -e MONGO_URI="mongodb://admin:admin123@mongo-prod:27017/?authSource=admin" \
                  -p 8087:8087 \
                  $TRACKING_IMAGE

                # ================= NOTIFICATION =================
                docker run -d \
                  --name prod-notification \
                  --network $NETWORK \
                  -e MONGO_URI="mongodb://admin:admin123@mongo-prod:27017/?authSource=admin" \
                  -p 8088:8088 \
                  $NOTIFICATION_IMAGE

                echo "WAIT SERVICES"
                sleep 20

                echo "CHECK CONTAINERS"

                docker ps | grep prod-payment || exit 1
                docker ps | grep prod-order || exit 1
                docker ps | grep prod-delivery || exit 1
                docker ps | grep prod-shipment || exit 1
                docker ps | grep prod-pickup || exit 1
                docker ps | grep prod-warehouse || exit 1
                docker ps | grep prod-tracking || exit 1
                docker ps | grep prod-notification || exit 1

                echo "ALL CONTAINERS RUNNING"
                '''
            }
        }

        // =====================================================
        // VERIFY
        // =====================================================
        stage('Verify') {
            steps {
                sh '''
                echo "VERIFY FULL SYSTEM"

                PAYMENT=$(curl -s -X POST http://host.docker.internal:8082/payment \
                -H "Content-Type: application/json" \
                -d '{"amount":1,"paid":1}')

                echo $PAYMENT | grep "PAID" || exit 1

                ORDER=$(curl -s -X POST http://host.docker.internal:8081/order \
                -H "Content-Type: application/json" \
                -d '{"user_id":1,"weight_kg":2,"distance_km":5,"base_price":10000}')

                TRACKING=$(echo $ORDER | jq -r '.tracking_number')

                if [ -z "$TRACKING" ] || [ "$TRACKING" = "null" ]; then
                    echo "Tracking not found"
                    exit 1
                fi

                curl -s -X POST http://host.docker.internal:8086/delivery \
                -H "Content-Type: application/json" \
                -d "{\"tracking_number\":\"$TRACKING\",\"address\":\"Bandung\"}" || exit 1

                curl -s -X POST http://host.docker.internal:8085/shipment \
                -H "Content-Type: application/json" \
                -d "{\"tracking_number\":\"$TRACKING\"}" || exit 1

                curl -s -X POST http://host.docker.internal:8083/pickup \
                -H "Content-Type: application/json" \
                -d '{"order_id":"ORD1","payment_status":"paid","weight":2}' || exit 1

                curl -s http://host.docker.internal:8084/warehouse || exit 1

                # ================= TRACKING VERIFY =================
                curl -s http://host.docker.internal:8087/tracking || true

                # ================= NOTIFICATION VERIFY =================
                curl -s http://host.docker.internal:8088/notification || true

                echo "ALL SERVICES VERIFIED"
                '''
            }
        }
    }

    // =====================================================
    // POST
    // =====================================================
    post {

        success {
            echo 'PIPELINE SUCCESS'
        }

        failure {
            echo 'PIPELINE FAILED'
        }

        always {

            sh '''
            docker network disconnect $NETWORK jenkins-server || true

            docker rm -f mongo-test mongo-prod || true

            docker rm -f \
            test-payment \
            test-order \
            test-delivery \
            test-shipment \
            test-pickup \
            test-warehouse \
            test-tracking \
            test-notification || true

            docker network rm $NETWORK || true
            '''
        }
    }
}
