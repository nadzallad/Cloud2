pipeline {
    agent {
        docker {
            image 'golang:1.22'
            args '-v /var/run/docker.sock:/var/run/docker.sock'
        }
    }

    environment {
        DOCKER_HUB_USER = 'nadzalla'
        DOCKER_HUB_ID   = 'logistic-login'
        GIT_REPO_URL    = 'https://github.com/nadzallad/Cloud2.git'

        IMAGE = "${DOCKER_HUB_USER}/payment-service:${env.BUILD_NUMBER}"
    }

    stages {

        stage('Checkout Repo') {
            steps {
                deleteDir()
                git branch: 'main', url: "${GIT_REPO_URL}"
            }
        }

        stage('Unit Test') {
            steps {
                dir('PaymentService') {
                    sh '''
                    echo "Running Unit Test..."
                    go test ./...
                    '''
                }
            }
        }

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

        stage('Build Image') {
            steps {
                sh '''
                echo "Building Docker Image..."
                docker build -t $IMAGE ./PaymentService
                '''
            }
        }

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

        stage('Deploy') {
            steps {
                sh '''
                echo "Deploy ke Kubernetes..."
                kubectl apply -f k8s/ --validate=false
                '''
            }
        }

        stage('Verify') {
            steps {
                sh 'echo "Pipeline SUCCESS"'
            }
        }
    }
}