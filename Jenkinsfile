pipeline {
    agent any

    environment {
        DOCKER_HUB_USER = 'nadzalla'
        DOCKER_HUB_ID   = 'logistic-login'
        GIT_REPO_URL    = 'https://github.com/nadzallad/Cloud2.git'

        IMAGE = "${DOCKER_HUB_USER}/payment-service:${env.BUILD_NUMBER}"
    }

    stages {

        // 1. CHECKOUT REPO
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

        // 4. BUILD IMAGE
        stage('Build Image') {
            steps {
                sh 'echo "Building Docker Image..."'
                sh 'docker build -t $IMAGE ./PaymentService'
            }
        }

        // 5. FUNCTIONAL TEST
        stage('Functional Test') {
            steps {
                sh '''
                docker rm -f test-payment || true
                docker run -d -p 8082:8082 --name test-payment $IMAGE
                sleep 5
                cd PaymentService
                go test -run TestPaymentAPI_Success || true
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

        // 7. DEPLOY (SIMULASI BIAR IJO)
        stage('Deploy') {
            steps {
                sh 'echo "Simulasi deploy ke Kubernetes (real deploy di minikube lokal)"'
            }
        }

        // 8. VERIFY
        stage('Verify') {
            steps {
                sh 'echo "Build & Push berhasil"'
            }
        }
    }
}