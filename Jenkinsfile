pipeline {
    agent any

    environment {
        PAYMENT_IMAGE = "payment-service"
        ORDER_IMAGE = "order-service"
        PICKUP_IMAGE = "pickup-service"
        TAG = "latest"
    }

    stages {

        stage('Checkout Repo') {
            steps {
                deleteDir()
                git branch: 'main', url: 'https://github.com/nadzallad/Cloud2.git'
            }
        }

        // ========================
        // UNIT TEST (FIX: skip functional)
        // ========================
        stage('Unit Test') {
            steps {
                bat '''
                go list ./... | findstr /V functional > packages.txt
                for /f %%i in (packages.txt) do go test %%i
                '''
            }
        }

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
                cd PaymentService
                docker build -t payment-service:latest .
                cd ..

                cd OrderService
                docker build -t order-service:latest .
                cd ..

                cd PickupService
                docker build -t pickup-service:latest .
                cd ..
                '''
            }
        }

        // ========================
        // FUNCTIONAL TEST
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

                timeout /t 6

                curl -X POST http://localhost:8081/payment ^
                -H "Content-Type: application/json" ^
                -d "{\\"amount\\":10000,\\"paid\\":10000}"

                curl -X POST http://localhost:8080/order ^
                -H "Content-Type: application/json" ^
                -d "{\\"user_id\\":1,\\"weight_kg\\":2,\\"distance_km\\":5,\\"base_price\\":10000}"

                curl -X POST http://localhost:8082/pickup ^
                -H "Content-Type: application/json" ^
                -d "{\\"order_id\\":\\"ORD1\\",\\"payment_status\\":\\"paid\\",\\"weight\\":2}"
                '''
            }
        }

        // ========================
        // PUSH IMAGE
        // ========================
        stage('Push Images') {
            steps {
                bat '''
                docker tag order-service:latest ghryalvrt/order-service:latest
                docker push ghryalvrt/order-service:latest

                docker tag payment-service:latest ghryalvrt/payment-service:latest
                docker push nadzalla/payment-service:latest

                docker tag pickup-service:latest ghryalvrt/pickup-service:latest
                docker push naurafaizah/pickup-service:latest
                '''
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

        stage('Verify Deployment') {
            steps {
                bat 'kubectl get pods && kubectl get svc'
            }
        }
    }
}
