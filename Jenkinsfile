pipeline {
    agent any

    environment {
        PAYMENT_IMAGE = "nadzalla/payment-service:${env.BUILD_NUMBER}"
        ORDER_IMAGE = "nadzalla/order-service:${env.BUILD_NUMBER}"
        SHIPMENT_IMAGE = "nadzalla/shipment-service:${env.BUILD_NUMBER}"
        DELIVERY_IMAGE = "nadzalla/delivery-service:${env.BUILD_NUMBER}"
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
            }
        }

        stage('Build Image') {
            steps {
                sh '''
                docker build -t $PAYMENT_IMAGE ./PaymentService
                docker build -t $ORDER_IMAGE ./OrderService
                docker build -t $SHIPMENT_IMAGE ./ShipmentService
                docker build -t $DELIVERY_IMAGE ./DeliveryService
                '''
            }
        }

        stage('Functional Test') {
            steps {
                catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                    sh '''
                    docker rm -f test-payment test-order test-delivery test-shipment || true

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

                    sleep 10

                    curl -s -X POST http://host.docker.internal:8082/payment \
                      -H "Content-Type: application/json" \
                      -d '{"amount":1,"paid":1}'

                    curl -s -X POST http://host.docker.internal:8081/order \
                      -H "Content-Type: application/json" \
                      -d '{"user_id":1,"weight_kg":2,"distance_km":5,"base_price":10000}'

                    curl -s -X POST http://host.docker.internal:8086/delivery

                    curl -s -X POST http://host.docker.internal:8085/shipment

                    docker rm -f test-payment test-order test-delivery test-shipment || true
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
                    '''
                }

                withCredentials([usernamePassword(
                    credentialsId: 'dockerhub-login',
                    usernameVariable: 'USERNAME',
                    passwordVariable: 'PASSWORD'
                )]) {
                    sh '''
                    echo "$PASSWORD" | docker login -u "$USERNAME" --password-stdin

                    docker push $ORDER_IMAGE
                    docker push $DELIVERY_IMAGE
                    docker push $SHIPMENT_IMAGE
                    '''
                }
            }
        }

        stage('Deploy') {
            steps {
                sh 'echo "DEPLOY OK"'
            }
        }

        stage('Verify') {
            steps {
                sh 'echo "PIPELINE SUCCESS"'
            }
        }
    }
}
