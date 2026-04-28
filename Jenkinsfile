pipeline {
    agent any

    environment {
        PAYMENT_IMAGE = "payment-service"
        ORDER_IMAGE = "order-service"
        TAG = "latest"
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
                bat 'go test ./...'
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
                '''
            }
        }

        // ========================
        // FUNCTIONAL TEST
        // ========================
        stage('Functional Test') {
            steps {
                bat '''
                start /b go run PaymentService/main.go
                start /b go run OrderService/main.go

                timeout /t 5

                curl -X POST http://localhost:8081/payment ^
                -H "Content-Type: application/json" ^
                -d "{\\"amount\\":10000,\\"paid\\":10000}"

                curl -X POST http://localhost:8080/order ^
                -H "Content-Type: application/json" ^
                -d "{\\"user_id\\":1,\\"weight_kg\\":2,\\"distance_km\\":5,\\"base_price\\":10000}"
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
                docker push ghryalvrt/payment-service:latest
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