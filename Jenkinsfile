pipeline {
    agent any

    environment {
        PAYMENT_IMAGE = "nadzalla/payment-service:${env.BUILD_NUMBER}"
        ORDER_IMAGE = "nadzalla/order-service:${env.BUILD_NUMBER}"
        SHIPMENT_IMAGE = "nadzalla/shipment-service:${env.BUILD_NUMBER}"
        DELIVERY_IMAGE = "nadzalla/delivery-service:${env.BUILD_NUMBER}"
        PICKUP_IMAGE = "nadzalla/pickup-service:${env.BUILD_NUMBER}"
        WAREHOUSE_IMAGE = "nadzalla/warehouse-service:${env.BUILD_NUMBER}"
    }

    stages {

        stage('Checkout Repo') {
            steps {
                deleteDir()
                git branch: 'main', url: 'https://github.com/nadzallad/Cloud2.git'
            }
        }

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
            }
        }

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
            }
        }

        stage('Build Image') {
            steps {
                sh '''
                docker build -t $PAYMENT_IMAGE ./PaymentService
                docker build -t $ORDER_IMAGE ./OrderService
                docker build -t $SHIPMENT_IMAGE ./ShipmentService
                docker build -t $DELIVERY_IMAGE ./DeliveryService
                docker build -t $PICKUP_IMAGE ./PickupService
                docker build -t $WAREHOUSE_IMAGE ./WarehouseService
                '''
            }
        }

        stage('Functional Test') {
            steps {
                catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                    sh '''
                    docker rm -f test-payment test-order test-delivery test-shipment test-pickup test-warehouse || true

                    docker run -d --name test-payment \
                      -e DB_HOST=host.docker.internal \
                      -e DB_NAME=payment_db \
                      -e DB_PASS=admin123 \
                      -p 8082:8082 \
                      $PAYMENT_IMAGE

                    docker run -d --name test-order \
                      -p 8081:8081 \
                      $ORDER_IMAGE

                    docker run -d --name test-delivery \
                      -e DB_HOST=host.docker.internal \
                      -e DB_NAME=delivery_db \
                      -e DB_USER=postgres \
                      -e DB_PASSWORD=admin123 \
                      -p 8086:8086 \
                      $DELIVERY_IMAGE

                    docker run -d --name test-shipment \
                      -e DB_HOST=host.docker.internal \
                      -e DB_NAME=shipment_db \
                      -e DB_USER=postgres \
                      -e DB_PASSWORD=admin123 \
                      -p 8085:8085 \
                      $SHIPMENT_IMAGE

                    docker run -d --name test-pickup \
                      -p 8083:8083 \
                      $PICKUP_IMAGE

                    docker run -d --name test-warehouse \
                      -p 8084:8084 \
                      $WAREHOUSE_IMAGE

                    sleep 10

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

                    docker rm -f test-payment test-order test-delivery test-shipment test-pickup test-warehouse || true
                    '''
                }
            }
        }

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
                    '''
                }
            }
        }

        stage('Deploy') {
    steps {
        sh '''
        echo "Starting FULL deployment..."

        # ================= CLEAN OLD =================
        docker rm -f prod-payment prod-order prod-delivery prod-shipment prod-pickup prod-warehouse || true

        # ================= PAYMENT =================
        docker run -d --name prod-payment \
          -e DB_HOST=host.docker.internal \
          -e DB_NAME=payment_db \
          -e DB_PASS=admin123 \
          -p 8082:8082 \
          $PAYMENT_IMAGE

        # ================= ORDER =================
        docker run -d --name prod-order \
          -p 8081:8081 \
          $ORDER_IMAGE

        # ================= DELIVERY =================
        docker run -d --name prod-delivery \
          -e DB_HOST=host.docker.internal \
          -e DB_NAME=delivery_db \
          -e DB_USER=postgres \
          -e DB_PASSWORD=admin123 \
          -p 8086:8086 \
          $DELIVERY_IMAGE

        # ================= SHIPMENT =================
        docker run -d --name prod-shipment \
          -e DB_HOST=host.docker.internal \
          -e DB_NAME=shipment_db \
          -e DB_USER=postgres \
          -e DB_PASSWORD=admin123 \
          -p 8085:8085 \
          $SHIPMENT_IMAGE

        # ================= PICKUP =================
        docker run -d --name prod-pickup \
          -p 8083:8083 \
          $PICKUP_IMAGE

        # ================= WAREHOUSE =================
        docker run -d --name prod-warehouse \
          -p 8084:8084 \
          $WAREHOUSE_IMAGE

        echo "Waiting services to be ready..."
        sleep 12

        # ================= HEALTH CHECK =================
        echo "Checking containers..."

        docker ps | grep prod-payment || exit 1
        docker ps | grep prod-order || exit 1
        docker ps | grep prod-delivery || exit 1
        docker ps | grep prod-shipment || exit 1
        docker ps | grep prod-pickup || exit 1
        docker ps | grep prod-warehouse || exit 1

        echo "ALL CONTAINERS RUNNING"
        '''
    }
}

        stage('Verify') {
            steps {
                sh '''
                echo "Verifying full system..."

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

                echo "ALL SERVICES OK"
                '''
            }
        }
    }
}
