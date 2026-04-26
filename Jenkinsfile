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
                bat 'go test ./PaymentService/...'
            }
        }

        stage('Vet') {
            steps {
                bat 'go vet ./...'
            }
        }
    }
}
