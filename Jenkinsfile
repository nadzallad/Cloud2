pipeline {
    agent any

    stages {

        stage('Checkout') {
            steps {
                deleteDir()
                git branch: 'main', url: 'https://github.com/nadzallad/Cloud2.git'
            }
        }

        stage('Unit Test') {
            steps {
                sh 'go test ./...'
            }
        }

        stage('Vet') {
            steps {
                sh 'go vet ./...'
            }
        }

        stage('Build Image') {
            steps {
                sh 'docker build -t payment-service:latest .'
            }
        }

        stage('Push Image') {
            steps {
                sh 'docker tag payment-service:latest yourdockerhub/payment-service:latest'
                sh 'docker push yourdockerhub/payment-service:latest'
            }
        }

        stage('Functional Test') {
            steps {
                sh './run-functional.sh'
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
                sh 'kubectl get svc'
            }
        }
    }
}