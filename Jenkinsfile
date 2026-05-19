pipeline {

    environment {
        DOCKER_HUB_USER = 'nadzalla'
        DOCKER_HUB_ID   = 'logistic-login'
        GIT_REPO_URL    = 'https://github.com/nadzallad/Cloud2.git'
        IMAGE = "${DOCKER_HUB_USER}/payment-service:${env.BUILD_NUMBER}"
    }

    triggers {
        // AUTO BUILD tiap ada perubahan (backup kalau webhook ga jalan)
        pollSCM('* * * * *')
    }

    stages {

        // 1. CHECKOUT
        stage('Checkout Repo') {
            steps {
                deleteDir()
                git branch: 'main', url: "${GIT_REPO_URL}"
            }
        }

        // 2. UNIT TEST
        stage('Unit Test') {
            steps {
                dir('PaymentService') {
                    sh '''
                    echo "Running Unit Test..."

                    go mod tidy
                    go mod download

                    go test -v ./...
                    '''
                }
            }
        }

        // 3. LINT
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

        // 4. BUILD IMAGE
        stage('Build Image') {
            steps {
                sh '''
                echo "Building Docker Image..."
                docker build -t $IMAGE ./PaymentService
                '''
            }
        }

        // 5. FUNCTIONAL TEST
        stage('Functional Test') {
            steps {
                sh '''
                echo "Start container..."
                docker rm -f test-payment 2>/dev/null || true
                docker run -d -p 8082:8082 --name test-payment $IMAGE

                echo "Waiting API..."
                until curl -s http://localhost:8082; do
                  sleep 2
                done

                echo "Run Functional Test..."
                cd PaymentService
                go test -run TestPaymentAPI_Success

                docker stop test-payment
                docker rm test-payment
                '''
            }
        }

        // 6. PUSH IMAGE
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

        // 7. DEPLOY
        stage('Deploy') {
            steps {
                sh '''
                echo "Install kubectl..."
                curl -LO https://dl.k8s.io/release/v1.29.0/bin/linux/amd64/kubectl
                chmod +x kubectl
                mv kubectl /usr/local/bin/

                echo "Deploy to Kubernetes..."
                kubectl apply -f k8s/ --validate=false
                '''
            }
        }

        // 8. VERIFY
        stage('Verify') {
            steps {
                sh 'echo "PIPELINE SUCCESS"'
            }
        }
    }

    post {
        success {
            echo '✅ SUCCESS: Pipeline completed!'
        }
        failure {
            echo '❌ FAILED: Check logs!'
        }
    }
}