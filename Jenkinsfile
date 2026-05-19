pipeline {
    agent any

    environment {
        DOCKER_HUB_USER = 'nadzalla'
        DOCKER_HUB_ID   = 'logistic-login'
        GIT_REPO_URL    = 'https://github.com/nadzallad/Cloud2.git'

        IMAGE = "${DOCKER_HUB_USER}/payment-service:${env.BUILD_NUMBER}"
    }

    stages {

        // 1. CHECKOUT
        stage('Checkout Repo') {
            steps {
                deleteDir()
                git branch: 'main', url: "${GIT_REPO_URL}"
            }
        }

        // 2. UNIT TEST (HARUS PASS)
        stage('Unit Test') {
            steps {
                dir('PaymentService') {
                    sh '''
                    echo "===== DEBUG ====="
                    pwd
                    ls -la

                    echo "===== GO CHECK ====="
                    which go || echo "GO NOT FOUND"
                    go version || true

                    echo "===== GO MOD ====="
                    go mod tidy
                    go mod download

                    echo "===== RUN TEST ====="
                    go test -v ./...
                    '''
                }
            }
        }

        // 3. LINT / VET (HARUS BERSIH)
        stage('Lint / Vet') {
            steps {
                dir('PaymentService') {
                    sh '''
                    echo "Running Go Vet..."
                    go vet ./...
                    '''
                }
            }
        }

        // 4. BUILD IMAGE (HARUS BERHASIL)
        stage('Build Image') {
            steps {
                sh '''
                echo "Building Docker Image..."
                docker build -t $IMAGE ./PaymentService
                '''
            }
        }

        // 5. FUNCTIONAL TEST (HARUS PASS)
        stage('Functional Test') {
            steps {
                sh '''
                echo "Start container for testing..."
                docker rm -f test-payment 2>/dev/null || true
                docker run -d -p 8082:8082 --name test-payment $IMAGE

                echo "Waiting service..."
                sleep 5

                echo "Running Functional Test..."
                cd PaymentService
                go test -run TestPaymentAPI_Success

                echo "Cleanup..."
                docker stop test-payment
                docker rm test-payment
                '''
            }
        }

        // 6. PUSH IMAGE (HARUS BERHASIL)
        stage('Push Image') {
            steps {
                withCredentials([usernamePassword(
                    credentialsId: "${DOCKER_HUB_ID}",
                    usernameVariable: 'USERNAME',
                    passwordVariable: 'PASSWORD'
                )]) {
                    sh '''
                    echo "$PASSWORD" | docker login -u "$USERNAME" --password-stdin
                    docker push $IMAGE
                    '''
                }
            }
        }

        // 7. DEPLOY (HARUS BERHASIL)
        stage('Deploy') {
            steps {
                sh '''
                echo "Deploy ke Kubernetes..."
                kubectl apply -f k8s/ --validate=false
                '''
            }
        }

        // 8. VERIFY
        stage('Verify') {
            steps {
                sh '''
                echo "Pipeline SUCCESS"
                '''
            }
        }
    }
}