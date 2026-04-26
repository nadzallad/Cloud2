pipeline {
    agent any

    environment {
        IMAGE_NAME = "payment-service"
        IMAGE_TAG = "latest"
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

        stage('Build Docker Image') {
            steps {
                bat '''
                cd PaymentService
                docker build -t payment-service:latest .
                cd ..
                '''
            }
        }

        stage('Functional Test') {
            steps {
                bat '''
                start /b go run PaymentService/main.go
                timeout /t 5

                curl -X POST http://localhost:8081/payment ^
                -H "Content-Type: application/json" ^
                -d "{\"amount\":10000,\"paid\":10000}"
                '''
            }
        }

        stage('Push Image') {
            steps {
                bat '''
                docker tag payment-service:latest nadzalla/payment-service:latest
                docker push nadzalla/payment-service:latest
                '''
            }
        }

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
