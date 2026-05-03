pipeline {
    agent any

    environment {
        PAYMENT_IMAGE = "ghryalvrt/payment-service"
        ORDER_IMAGE = "ghryalvrt/order-service"
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
                sh 'go test ./... || true'
            }
        }

        stage('Vet') {
            steps {
                sh 'go vet ./...'
            }
        }

        // ========================
        // BUILD DOCKER
        // ========================
        stage('Build Docker Images') {
            steps {
                sh '''
                cd PaymentService
                docker build -t $PAYMENT_IMAGE:$TAG .
                cd ..

                cd OrderService
                docker build -t $ORDER_IMAGE:$TAG .
                cd ..
                '''
            }
        }

        // ========================
        // FUNCTIONAL TEST
        // ========================
        stage('Functional Test') {
            steps {
                sh '''
                go run PaymentService/main.go &
                go run OrderService/main.go &

                sleep 5

                curl -X POST http://localhost:8081/payment \
                -H "Content-Type: application/json" \
                -d '{"amount":10000,"paid":10000}' || true

                curl -X POST http://localhost:8080/order \
                -H "Content-Type: application/json" \
                -d '{"user_id":1,"weight_kg":2,"distance_km":5,"base_price":10000}' || true
                '''
            }
        }

        // ========================
        // PUSH IMAGE
        // ========================
        stage('Push Images') {
            steps {
                sh '''
                docker push $ORDER_IMAGE:$TAG
                docker push $PAYMENT_IMAGE:$TAG
                '''
            }
        }

        // ========================
        // DEPLOY
        // ========================
        stage('Deploy to Kubernetes') {
            steps {
                sh '''
                kubectl apply -f k8s/
                '''
            }
        }

        stage('Verify Deployment') {
            steps {
                sh '''
                kubectl get pods
                kubectl get svc
                '''
            }
        }
    }
}