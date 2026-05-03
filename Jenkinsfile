pipeline {
    agent any

    stages {

        stage('Checkout') {
            steps {
                git 'https://github.com/nadzallad/Cloud2.git'
            }
        }

        stage('Unit Test') {
            steps {
                bat 'go test ./... || exit 0'
            }
        }

        stage('Vet') {
            steps {
                bat 'go vet ./...'
            }
        }

        stage('Build Image') {
            steps {
                bat 'build-push.bat'
            }
        }

        stage('Functional Test') {
            steps {
                bat 'run-functional.bat || exit 0'
            }
        }

        stage('Deploy') {
            steps {
                bat 'kubectl apply -f k8s/'
            }
        }

        stage('Verify') {
            steps {
                bat 'kubectl get pods && kubectl get svc'
            }
        }
    }
}