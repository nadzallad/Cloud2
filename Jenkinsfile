pipeline {
    agent any

    environment {
        IMAGE_NAME = "payment-service"
        IMAGE_TAG = "latest"
    }

    stages {

        stage('Checkout Repo') {
            steps {
                deleteDir()
                git branch: 'main', url: 'https://github.com/nadzallad/Cloud2.git'
            }
        }

        stage('Unit Test') {
            steps {
                bat '''
                for /d %%i in (*) do (
                    if exist %%i\\go.mod (
                        echo ===== UNIT TEST %%i =====
                        cd %%i
                        go test ./...
                        if errorlevel 1 exit /b 1
                        cd ..
                    )
                )
                '''
            }
        }

        stage('Vet') {
            steps {
                bat '''
                for /d %%i in (*) do (
                    if exist %%i\\go.mod (
                        echo ===== VET %%i =====
                        cd %%i
                        go vet ./...
                        if errorlevel 1 exit /b 1
                        cd ..
                    )
                )
                '''
            }
        }

        stage('Build Docker Image') {
            steps {
                bat '''
                cd PaymentService
                docker build -t payment-service:latest .
                cd ..
                '''
            }
        }

        stage('Functional Test') {
            steps {
                bat '''
                echo Starting app...
                cd PaymentService
                start /b go run main.go

                timeout /t 5

                curl -X POST http://localhost:8081/payment ^
                -H "Content-Type: application/json" ^
                -d "{\"amount\":10000,\"paid\":10000}"
                '''
            }
        }

        stage('Push Image') {
            steps {
                bat '''
                docker tag payment-service:latest yourdockerhub/payment-service:latest
                docker push yourdockerhub/payment-service:latest
                '''
            }
        }

        stage('Deploy to Kubernetes') {
            steps {
                bat '''
                kubectl apply -f k8s/
                '''
            }
        }

        stage('Verify Deployment') {
            steps {
                bat '''
                kubectl get pods
                kubectl get svc
                '''
            }
        }
    }
}