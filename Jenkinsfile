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

        TEST_NETWORK = "test-net"
        PROD_NETWORK = "prod-net"
    }

    stages {
        // CHECKOUT
        stage('Checkout Repo') {
            steps {
                deleteDir()
                git branch: 'main', url: 'https://github.com/nadzallad/Cloud2.git'
            }
        }
        
        // UNIT TEST
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
        
        // LINT / VET
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

        // BUILD IMAGE
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

        // FUNCTIONAL TEST
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
                    
                    # ================= NETWORK =================
                    docker network inspect $TEST_NETWORK >/dev/null 2>&1 || docker network create $TEST_NETWORK
                    
                    # ================= MONGODB =================
                    docker run -d \
                      --name mongo-test \
                      --network $TEST_NETWORK \
                      -e MONGO_INITDB_ROOT_USERNAME=admin \
                      -e MONGO_INITDB_ROOT_PASSWORD=admin123 \
                      mongo
                    
                    echo "WAIT MONGO"
                    sleep 15
                    
                    # ================= SERVICES =================
                    docker run -d --name test-payment \
                      --network $TEST_NETWORK \
                      -e DB_HOST=host.docker.internal \
                      -e DB_NAME=payment_db \
                      -e DB_PASS=admin123 \
                      -p 18082:8082 \
                      $PAYMENT_IMAGE
                    
                    docker run -d --name test-order \
                      --network $TEST_NETWORK \
                      -p 18081:8081 \
                      $ORDER_IMAGE
                    
                    docker run -d --name test-delivery \
                      --network $TEST_NETWORK \
                      -e DB_HOST=host.docker.internal \
                      -e DB_NAME=delivery_db \
                      -e DB_USER=postgres \
                      -e DB_PASSWORD=admin123 \
                      -p 18086:8086 \
                      $DELIVERY_IMAGE
                    
                    docker run -d --name test-shipment \
                      --network $TEST_NETWORK \
                      -e DB_HOST=host.docker.internal \
                      -e DB_NAME=shipment_db \
                      -e DB_USER=postgres \
                      -e DB_PASSWORD=admin123 \
                      -p 18085:8085 \
                      $SHIPMENT_IMAGE
                    
                    docker run -d --name test-pickup \
                      --network $TEST_NETWORK \
                      -p 18083:8083 \
                      $PICKUP_IMAGE
                    
                    docker run -d --name test-warehouse \
                      --network $TEST_NETWORK \
                      -p 18084:8084 \
                      $WAREHOUSE_IMAGE
                    
                    docker run -d \
                      --name test-tracking \
                      --network $TEST_NETWORK \
                      -e MONGO_URI="mongodb://admin:admin123@mongo-test:27017/?authSource=admin" \
                      -p 18087:8087 \
                      $TRACKING_IMAGE
                    
                    docker run -d \
                      --name test-notification \
                      --network $TEST_NETWORK \
                      -e MONGO_URI="mongodb://admin:admin123@mongo-test:27017/?authSource=admin" \
                      -p 18088:8088 \
                      $NOTIFICATION_IMAGE
                    
                    echo "WAIT APPLICATION"
                    sleep 20
                    
                    # ================= API TEST =================
                    curl -s -X POST http://host.docker.internal:18082/payment \
                    -H "Content-Type: application/json" \
                    -d '{"amount":1,"paid":1}'
                    
                    curl -s -X POST http://host.docker.internal:18081/order \
                    -H "Content-Type: application/json" \
                    -d '{"user_id":1,"weight_kg":2,"distance_km":5,"base_price":10000}'
                    
                    curl -s http://host.docker.internal:18084/warehouse
                    
                    echo "FUNCTIONAL TEST DONE"
                    '''
                }
            }
        }

        // PUSH IMAGE
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

        // DEPLOY
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
                
                # ================= NETWORK =================
                docker network inspect $PROD_NETWORK >/dev/null 2>&1 || docker network create $PROD_NETWORK
                
                # ================= MONGODB =================
                docker run -d \
                  --name mongo-prod \
                  --network $PROD_NETWORK \
                  -e MONGO_INITDB_ROOT_USERNAME=admin \
                  -e MONGO_INITDB_ROOT_PASSWORD=admin123 \
                  mongo
                
                sleep 15
                
                # ================= SERVICES =================
                docker run -d --name prod-payment \
                  --network $PROD_NETWORK \
                  -e DB_HOST=host.docker.internal \
                  -e DB_NAME=payment_db \
                  -e DB_PASS=admin123 \
                  -p 8082:8082 \
                  $PAYMENT_IMAGE
                
                docker run -d --name prod-order \
                  --network $PROD_NETWORK \
                  -p 8081:8081 \
                  $ORDER_IMAGE
                
                docker run -d --name prod-delivery \
                  --network $PROD_NETWORK \
                  -e DB_HOST=host.docker.internal \
                  -e DB_NAME=delivery_db \
                  -e DB_USER=postgres \
                  -e DB_PASSWORD=admin123 \
                  -p 8086:8086 \
                  $DELIVERY_IMAGE
                
                docker run -d --name prod-shipment \
                  --network $PROD_NETWORK \
                  -e DB_HOST=host.docker.internal \
                  -e DB_NAME=shipment_db \
                  -e DB_USER=postgres \
                  -e DB_PASSWORD=admin123 \
                  -p 8085:8085 \
                  $SHIPMENT_IMAGE
                
                docker run -d --name prod-pickup \
                  --network $PROD_NETWORK \
                  -p 8083:8083 \
                  $PICKUP_IMAGE
                
                docker run -d --name prod-warehouse \
                  --network $PROD_NETWORK \
                  -p 8084:8084 \
                  $WAREHOUSE_IMAGE
                
                docker run -d \
                  --name prod-tracking \
                  --network $PROD_NETWORK \
                  -e MONGO_URI="mongodb://admin:admin123@mongo-prod:27017/?authSource=admin" \
                  -p 8087:8087 \
                  $TRACKING_IMAGE
                
                docker run -d \
                  --name prod-notification \
                  --network $PROD_NETWORK \
                  -e MONGO_URI="mongodb://admin:admin123@mongo-prod:27017/?authSource=admin" \
                  -p 8088:8088 \
                  $NOTIFICATION_IMAGE
                
                echo "WAIT SERVICES"
                sleep 20
                
                docker ps | grep prod-payment || exit 1
                docker ps | grep prod-order || exit 1
                
                echo "DEPLOY SUCCESS"
                '''
            }
        }

        // VERIFY
        stage('Verify') {
            steps {
                sh '''
                echo "VERIFY FULL SYSTEM"
        
                # ================= WAIT SERVICE READY =================
                for i in {1..10}; do
                  if curl -s http://host.docker.internal:8082/payment > /dev/null; then
                    echo "Payment service ready"
                    break
                  fi
                  echo "Waiting payment service..."
                  sleep 3
                done
        
                # ================= START VERIFY =================
                PAYMENT=$(curl -s -f -X POST http://host.docker.internal:8082/payment \
                -H "Content-Type: application/json" \
                -d '{"amount":1,"paid":1}')
        
                echo $PAYMENT | grep -q "PAID" || exit 1
        
                ORDER=$(curl -s -f -X POST http://host.docker.internal:8081/order \
                -H "Content-Type: application/json" \
                -d '{"user_id":1,"weight_kg":2,"distance_km":5,"base_price":10000}')
        
                TRACKING=$(echo $ORDER | sed -n 's/.*"tracking_number":"\\([^"]*\\)".*/\\1/p')
        
                if [ -z "$TRACKING" ]; then
                    echo "Tracking not found"
                    exit 1
                fi
        
                curl -s -f -X POST http://host.docker.internal:8086/delivery \
                -H "Content-Type: application/json" \
                -d "{\"tracking_number\":\"$TRACKING\",\"address\":\"Bandung\"}"
        
                curl -s -f -X POST http://host.docker.internal:8085/shipment \
                -H "Content-Type: application/json" \
                -d "{\"tracking_number\":\"$TRACKING\"}"
        
                curl -s -f -X POST http://host.docker.internal:8083/pickup \
                -H "Content-Type: application/json" \
                -d '{"order_id":"ORD1","payment_status":"paid","weight":2}'
        
                curl -s -f http://host.docker.internal:8084/warehouse
        
                echo "ALL SERVICES VERIFIED"
                '''
            }
        }
    }

    // POST
    post {

        success {
            echo 'PIPELINE SUCCESS'
        }

        failure {
            echo 'PIPELINE FAILED'
        }

        always {

            sh '''
            docker rm -f mongo-test mongo-prod || true
            
            docker rm -f \
            test-payment test-order test-delivery test-shipment \
            test-pickup test-warehouse test-tracking test-notification || true
            
            docker network rm $TEST_NETWORK || true
            '''
        }
    }
}
