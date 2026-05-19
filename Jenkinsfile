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

        // 2. UNIT TEST (MERAH TAPI LANJUT)
        stage('Unit Test') {
            steps {
                dir('PaymentService') {
                    echo "Running Unit Test..."
                    catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                        sh 'go test ./...'
                    }
                }
            }
        }

        // 3. LINT / VET (WAJIB HIJAU)
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

        // 4. BUILD DOCKER IMAGE (WAJIB BERHASIL)
        stage('Build Image') {
            steps {
                sh '''
                echo "Building Docker Image..."
                docker build -t $IMAGE ./PaymentService
                '''
            }
        }

        // 5. FUNCTIONAL TEST (MERAH TAPI LANJUT)
        stage('Functional Test') {
            steps {
                catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                    sh '''
                    echo "Run DB container..."
                    docker rm -f postgres-test || true
                    docker run -d \
                    --name postgres-test \
                    -e POSTGRES_PASSWORD=123 \
                    -e POSTGRES_DB=testdb \
                    --health-cmd="pg_isready" \
                    --health-interval=2s \
                    postgres

                    echo "Waiting DB ready..."
                    until [ "$(docker inspect -f {{.State.Health.Status}} postgres-test)" = "healthy" ]; do
                    sleep 1
                    done

                    echo "Run app container..."
                    docker rm -f test-payment || true
                    docker run -d -p 8082:8082 --name test-payment $IMAGE

                    echo "Waiting API..."
                    until curl -s http://localhost:8082; do
                    sleep 1
                    done

                    echo "Run Functional Test..."
                    cd PaymentService
                    go test -run TestPaymentAPI_Success

                    echo "Cleanup..."
                    docker rm -f test-payment
                    docker rm -f postgres-test
                    '''
                }
            }
        }

        // 6. PUSH IMAGE (WAJIB BERHASIL)
        stage('Push Image') {
            steps {
                withCredentials([usernamePassword(
                    credentialsId: 'logistic-login',
                    usernameVariable: 'USERNAME',
                    passwordVariable: 'PASSWORD'
                )]) {
                    sh '''
                    echo "Login Docker Hub..."
                    echo "$PASSWORD" | docker login -u "$USERNAME" --password-stdin

                    echo "Push Docker Image..."
                    docker push $IMAGE
                    '''
                }
            }
        }

        // 7. DEPLOY (SIMULASI BIAR AMAN)
        stage('Deploy') {
            steps {
                sh '''
                echo "Simulasi deploy ke Kubernetes"
                echo "kubectl apply -f k8s/"
                '''
            }
        }

        // 8. VERIFY
        stage('Verify') {
            steps {
                sh '''
                echo "Build, Push, dan Pipeline selesai"
                '''
            }
        }
    }

    post {
        success {
            echo '✅ PIPELINE SUCCESS (meskipun test merah)'
        }
        failure {
            echo '❌ PIPELINE FAILED (cek vet/build/push)'
        }
    }
}