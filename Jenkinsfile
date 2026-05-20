pipeline {
    agent any

    environment {
        IMAGE = "nadzalla/payment-service:${env.BUILD_NUMBER}"
    }

    stages {

        // 1. CHECKOUT
        stage('Checkout Repo') {
            steps {
                deleteDir()
                git branch: 'main', url: 'https://github.com/nadzallad/Cloud2.git'
            }
        }

        // 2. UNIT TEST (REAL LOGIC)
        stage('Unit Test') {
            steps {
                dir('PaymentService') {
                    catchError(stageResult: 'FAILURE') {
                        sh 'go test -v -run TestValidatePayment ./...'
                    }
                }
            }
        }

        // 3. LINT
        stage('Lint / Vet') {
            steps {
                dir('PaymentService') {
                    sh 'go vet ./...'
                }
            }
        }

        // 4. BUILD
        stage('Build Image') {
            steps {
                sh 'docker build -t $IMAGE ./PaymentService'
            }
        }

        // 5. FUNCTIONAL TEST (REAL API + DB)
        stage('Functional Test') {
            steps {
                catchError(stageResult: 'FAILURE') {
                    sh '''
                    echo "RUN APP"

                    docker rm -f test-payment || true

                    docker run -d \
                      --name test-payment \
                      -e DB_HOST=host.docker.internal \
                      -e DB_NAME=payment_db \
                      -e DB_PASS=admin123 \
                      -p 8082:8082 \
                      $IMAGE

                    echo "WAIT APP"
                    sleep 5

                    echo "RUN FUNCTIONAL TEST"
                    cd PaymentService
                    go test -v -run TestPaymentAPI ./...
                    '''
                }
            }
        }

        // 6. PUSH
        stage('Push Image') {
            steps {
                withCredentials([usernamePassword(
                    credentialsId: 'logistic-login',
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

        // 7. DEPLOY
        stage('Deploy') {
            steps {
                catchError(stageResult: 'FAILURE') {
                    sh '''
                    docker rm -f prod-payment || true

                    docker run -d \
                      --name prod-payment \
                      -p 8083:8082 \
                      -e DB_HOST=host.docker.internal \
                      $IMAGE
                    '''
                }
            }
        }

        // 8. VERIFY (REAL CHECK)
        stage('Verify') {
            steps {
                catchError(stageResult: 'FAILURE') {
                    sh '''
                    echo "VERIFY API"

                    sleep 5

                    RESPONSE=$(curl -s -X POST http://localhost:8083/payment \
                      -H "Content-Type: application/json" \
                      -d '{
                        "order_id":2,
                        "amount":10000,
                        "paid":10000,
                        "payment_method":"BANK_TRANSFER"
                      }')

                    echo "Response: $RESPONSE"

                    echo "$RESPONSE" | grep PAID
                    '''
                }
            }
        }
    }

    post {
        always {
            echo 'PIPELINE SELESAI (REAL TEST)'
        }
    }
}