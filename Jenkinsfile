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
                    echo "Cleanup container lama..."
                    docker rm -f test-payment || true

                    echo "Run container..."
                    docker run -d -p 8082:8082 --name test-payment $IMAGE

                    echo "Tunggu service hidup..."
                    sleep 5

                    echo "Running Functional Test..."
                    cd PaymentService
                    go test -run TestPaymentAPI_Success

                    echo "Stop & remove container..."
                    docker stop test-payment
                    docker rm test-payment
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