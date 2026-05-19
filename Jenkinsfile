pipeline {
    agent any

    environment {
        IMAGE = "nadzalla/payment-service:${env.BUILD_NUMBER}"
    }

    stages {

        // 1. CHECKOUT REPO
        stage('Checkout Repo') {
            steps {
                deleteDir()
                git branch: 'main', url: 'https://github.com/nadzallad/Cloud2.git'
            }
        }

        // 2. UNIT TEST (TIDAK AKSES DB)
        stage('Unit Test') {
            steps {
                dir('PaymentService') {
                    sh 'echo "Running Unit Test..."'
                    sh 'go test ./... || true'
                }
            }
        }

        // 3. LINT / VET
        stage('Lint / Vet') {
            steps {
                dir('PaymentService') {
                    sh 'echo "Running Go Vet..."'
                    sh 'go vet ./... || true'
                }
            }
        }

        // 4. BUILD DOCKER IMAGE
        stage('Build Image') {
            steps {
                sh 'echo "Building Docker Image..."'
                sh 'docker build -t $IMAGE ./PaymentService'
            }
        }

        // 5. FUNCTIONAL TEST (PAKAI GO TEST)
        stage('Functional Test') {
            steps {
                sh '''
                echo "Cleanup container lama..."
                docker rm -f test-payment || true

                echo "Run container..."
                docker run -d -p 8082:8082 --name test-payment $IMAGE

                echo "Tunggu service hidup..."
                sleep 5

                echo "Running Functional Test..."
                cd PaymentService
                go test -run TestPaymentAPI_Success || true

                echo "Stop & remove container..."
                docker stop test-payment
                docker rm test-payment
                '''
            }
        }

        // 6. PUSH IMAGE KE DOCKER HUB
        stage('Push Image') {
            steps {
                sh 'echo "Push Docker Image..."'
                sh 'docker push $IMAGE'
            }
        }

        // 7. DEPLOY KE KUBERNETES
        stage('Deploy') {
            steps {
                sh 'kubectl config use-context minikube'

                sh '''
                sed -i "s|image: .*|image: $IMAGE|g" k8s/deployment-payment.yaml
                '''

                sh 'kubectl apply -f k8s/deployment-payment.yaml --validate=false'
                sh 'kubectl apply -f k8s/payment-service.yaml --validate=false'
                sh 'kubectl apply -f k8s/ingress.yaml --validate=false'
            }
        }

        // 8. VERIFY DEPLOYMENT
        stage('Verify') {
            steps {
                sh 'echo "Cek pods..."'
                sh 'kubectl get pods'

                sh 'echo "Cek services..."'
                sh 'kubectl get svc'
            }
        }
    }
}