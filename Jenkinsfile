pipeline {
    agent any

    environment {
        IMAGE = "nadzallad/payment-service:${env.BUILD_NUMBER}"
    }

    stages {

        stage('Checkout') {
            steps {
                git branch: 'main', url: 'https://github.com/nadzallad/Cloud2.git'
            }
        }

        stage('Test Payment') {
            steps {
                dir('PaymentService') {
                    sh 'go test ./...'
                }
            }
        }

        stage('Build Image') {
            steps {
                sh 'docker build -t $IMAGE ./PaymentService'
            }
        }

        stage('Push Image') {
            steps {
                sh 'docker push $IMAGE'
            }
        }

        stage('Deploy') {
            steps {
                sh 'kubectl apply -f k8s/'
            }
        }

        stage('Verify') {
            steps {
                sh 'kubectl get pods'
            }
        }
    }
}