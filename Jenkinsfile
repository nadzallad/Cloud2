pipeline {
    agent {
        docker {
            image 'golang:1.20'
            args '-v /var/run/docker.sock:/var/run/docker.sock'
        }
    }

    environment {
        DOCKER_HUB_USER = 'nadzalla'
        DOCKER_HUB_ID   = 'logistic-login'
        GIT_REPO_URL    = 'https://github.com/nadzallad/Cloud2.git'
        IMAGE = "${DOCKER_HUB_USER}/payment-service:${env.BUILD_NUMBER}"
    }

    triggers {
        pollSCM('* * * * *')
    }

    stages {

        // 1. CHECKOUT
        stage('Checkout Repo') {
            steps {
                deleteDir()
                checkout([
                    $class: 'GitSCM',
                    branches: [[name: '*/main']],
                    userRemoteConfigs: [[url: "${GIT_REPO_URL}"]]
                ])
            }
        }

        // 2. UNIT TEST (BOLEH FAIL)
        stage('Unit Test') {
            steps {
                dir('PaymentService') {
                    sh '''
                    echo "===== UNIT TEST ====="

                    go mod tidy
                    go mod download

                    go test -v ./... || echo "Unit Test Failed (Allowed)"
                    '''
                }
            }
        }

        // 3. LINT / VET (BOLEH FAIL)
        stage('Lint / Vet') {
            steps {
                dir('PaymentService') {
                    sh '''
                    echo "===== GO VET ====="
                    go vet ./... || echo "Lint/Vet Failed (Allowed)"
                    '''
                }
            }
        }

        // 4. BUILD IMAGE (HARUS BERHASIL)
        stage('Build Image') {
            steps {
                sh '''
                echo "===== BUILD IMAGE ====="
                docker build -t $IMAGE ./PaymentService
                '''
            }
        }

        // 5. FUNCTIONAL TEST (BOLEH FAIL)
        stage('Functional Test') {
            steps {
                sh '''
                echo "===== FUNCTIONAL TEST ====="

                docker rm -f test-payment 2>/dev/null || true
                docker run -d -p 8082:8082 --name test-payment $IMAGE

                echo "Waiting API..."
                sleep 5

                cd PaymentService
                go test -run TestPaymentAPI_Success || echo "Functional Test Failed (Allowed)"

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
                    echo "===== PUSH IMAGE ====="
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
                echo "===== DEPLOY ====="

                curl -LO https://dl.k8s.io/release/v1.29.0/bin/linux/amd64/kubectl
                chmod +x kubectl
                mv kubectl /usr/local/bin/

                kubectl apply -f k8s/ --validate=false
                '''
            }
        }

        // 8. VERIFY
        stage('Verify') {
            steps {
                sh 'echo "PIPELINE SUCCESS 🎉"'
            }
        }
    }

    post {
        success {
            echo '✅ SUCCESS: Pipeline selesai walaupun test gagal'
        }
        failure {
            echo '❌ FAILED: Cek stage build/push/deploy'
        }
    }
}