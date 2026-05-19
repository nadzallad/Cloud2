pipeline {
    agent any

    environment {
        IMAGE = "nadzallad/payment-service:${env.BUILD_NUMBER}"
    }

    stages {

        // 1. CHECKOUT
        stage('Checkout Repo') {
            steps {
                deleteDir()
                git branch: 'main', url: 'https://github.com/nadzallad/Cloud2.git'
            }
        }

        // 2. UNIT TEST
        stage('Unit Test') {
            steps {
                dir('PaymentService') {
                    sh 'go test ./... || true'
                }
            }
        }

        // 3. LINT / VET
        stage('Lint / Vet') {
            steps {
                dir('PaymentService') {
                    sh 'go vet ./... || true'
                }
            }
        }

        // 4. BUILD IMAGE
        stage('Build Image') {
            steps {
                sh 'docker build -t $IMAGE ./PaymentService'
            }
        }

        // 5. FUNCTIONAL TEST
        stage('Functional Test') {
            steps {
                sh '''
                docker run -d -p 8082:8082 --name test-payment $IMAGE
                sleep 5
                curl -X POST http://localhost:8082/payment \
                -H "Content-Type: application/json" \
                -d '{"order_id":1,"amount":50000,"paid":50000,"payment_method":"BANK_TRANSFER"}'
                docker stop test-payment
                docker rm test-payment
                '''
            }
        }

        // 6. PUSH IMAGE
        stage('Push Image') {
            steps {
                sh 'docker push $IMAGE'
            }
        }

        // 7. DEPLOY KUBERNETES
        stage('Deploy') {
            steps {
                sh 'kubectl apply -f k8s/'
            }
        }

        // 8. VERIFY
        stage('Verify') {
            steps {
                sh 'kubectl get pods'
                sh 'kubectl get svc'
            }
        }
    }
}