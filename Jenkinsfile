pipeline {
    agent any

    environment {
        DOCKER_USER = "ghryalvrt"
        TAG = "latest"
    }

    stages {

        // ========================
        // CHECKOUT
        // ========================
        stage('Checkout Repo') {
            steps {
                deleteDir()
                git branch: 'main', url: 'https://github.com/nadzallad/Cloud2.git'
            }
        }

        // ========================
        // UNIT TEST (BOLEH FAIL)
        // ========================
        stage('Unit Test') {
            steps {
                bat 'go test ./... || exit 0'
            }
        }

        // ========================
        // LINT / VET
        // ========================
        stage('Vet') {
            steps {
                bat 'go vet ./...'
            }
        }

        // ========================
        // BUILD DOCKER
        // ========================
        stage('Build Docker Images') {
            steps {
                bat '''
                docker build -t %DOCKER_USER%/payment-service:%TAG% PaymentService
                docker build -t %DOCKER_USER%/order-service:%TAG% OrderService
                docker build -t %DOCKER_USER%/pickup-service:%TAG% PickupService
                docker build -t %DOCKER_USER%/warehouse-service:%TAG% WarehouseService
                docker build -t %DOCKER_USER%/shipment-service:%TAG% ShipmentService
                docker build -t %DOCKER_USER%/delivery-service:%TAG% DeliveryService
                docker build -t %DOCKER_USER%/notification-service:%TAG% NotificationService
                docker build -t %DOCKER_USER%/tracking-service:%TAG% TrackingService
                '''
            }
        }

        // ========================
        // FUNCTIONAL TEST (BOLEH FAIL)
        // ========================
        stage('Functional Test') {
            steps {
                bat '''
                cd PaymentService
                start /b go run .
                cd ..

                cd OrderService
                start /b go run .
                cd ..

                cd PickupService
                start /b go run .
                cd ..

                cd WarehouseService
                start /b go run .
                cd ..

                cd ShipmentService
                start /b go run .
                cd ..

                cd DeliveryService
                start /b go run .
                cd ..

                cd NotificationService
                start /b go run .
                cd ..

                cd TrackingService
                start /b go run .
                cd ..

                timeout /t 8

                curl -X POST http://localhost:8082/payment ^
                -H "Content-Type: application/json" ^
                -d "{\\"amount\\":10000,\\"paid\\":10000}" || exit 0

                curl -X POST http://localhost:8081/order ^
                -H "Content-Type: application/json" ^
                -d "{\\"user_id\\":1,\\"weight_kg\\":2,\\"distance_km\\":5,\\"base_price\\":10000}" || exit 0

                curl -X POST http://localhost:8083/pickup ^
                -H "Content-Type: application/json" ^
                -d "{\\"order_id\\":\\"ORD1\\",\\"payment_status\\":\\"paid\\",\\"weight\\":2}" || exit 0

                curl -X POST http://localhost:8084/warehouse ^
                -H "Content-Type: application/json" ^
                -d "{\\"stock\\":10}" || exit 0

                curl -X POST http://localhost:8085/shipment ^
                -H "Content-Type: application/json" ^
                -d "{\\"order_id\\":\\"ORD1\\",\\"status\\":\\"shipped\\"}" || exit 0

                curl -X POST http://localhost:8086/delivery ^
                -H "Content-Type: application/json" ^
                -d "{\\"order_id\\":\\"ORD1\\",\\"status\\":\\"delivered\\"}" || exit 0

                curl -X POST http://localhost:8087/track ^
                -H "Content-Type: application/json" ^
                -d "{\\"order_id\\":\\"ORD1\\",\\"status\\":\\"on the way\\"}" || exit 0

                curl -X POST http://localhost:8088/notify ^
                -H "Content-Type: application/json" ^
                -d "{\\"message\\":\\"order created\\"}" || exit 0
                '''
            }
        }

        // ========================
        // PUSH IMAGE
        // ========================
        stage('Push Images') {
            steps {
                withCredentials([usernamePassword(
                    credentialsId: 'docker-creds',
                    usernameVariable: 'DOCKER_USER_LOGIN',
                    passwordVariable: 'DOCKER_PASS'
                )]) {
                    bat '''
                    echo %DOCKER_PASS% | docker login -u %DOCKER_USER_LOGIN% --password-stdin

                    docker push %DOCKER_USER%/payment-service:%TAG%
                    docker push %DOCKER_USER%/order-service:%TAG%
                    docker push %DOCKER_USER%/pickup-service:%TAG%
                    docker push %DOCKER_USER%/warehouse-service:%TAG%
                    docker push %DOCKER_USER%/shipment-service:%TAG%
                    docker push %DOCKER_USER%/delivery-service:%TAG%
                    docker push %DOCKER_USER%/notification-service:%TAG%
                    docker push %DOCKER_USER%/tracking-service:%TAG%
                    '''
                }
            }
        }

        // ========================
        // DEPLOY
        // ========================
        stage('Deploy to Kubernetes') {
            steps {
                bat 'kubectl apply -f k8s/'
            }
        }

        // ========================
        // VERIFY
        // ========================
        stage('Verify Deployment') {
            steps {
                bat '''
                kubectl get pods
                kubectl get svc
                '''
            }
        }
    }
}