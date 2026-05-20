pipeline {
    agent any

    environment {
        IMAGE = "nadzalla/payment-service:${env.BUILD_NUMBER}"
    }

    stages {

        stage('Checkout Repo') {
            steps {
                deleteDir()
                git branch: 'main', url: 'https://github.com/nadzallad/Cloud2.git'
            }
        }

        // UNIT TEST
        stage('Unit Test') {
            steps {
                dir('PaymentService') {
                    catchError(stageResult: 'FAILURE') {
                        sh 'go test -run TestValidatePayment ./...'
                    }
                }
            }
        }

        // LINT
        stage('Lint / Vet') {
            steps {
                dir('PaymentService') {
                    sh 'go vet ./...'
                }
            }
        }

        // BUILD
        stage('Build Image') {
            steps {
                sh 'docker build -t $IMAGE ./PaymentService'
            }
        }

        // FUNCTIONAL (PAKE DB LU LANGSUNG)
      stage('Functional Test') {
            steps {
                catchError(stageResult: 'FAILURE') {
                    sh '''
                    docker rm -f test-payment || true

                    docker run -d \
                    --name test-payment \
                    -e DB_HOST=host.docker.internal \
                    -e DB_NAME=payment_db \
                    -e DB_PASS=admin123 \
                    -p 8082:8082 \
                    $IMAGE

                    sleep 3

                    cd PaymentService
                    go test -run TestFunctional ./...
                    '''
                }
            }
        }
        
        // PUSH
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

        // DEPLOY
        stage('Deploy') {
            steps {
                catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
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

        // VERIFY (REAL HIT API)
        stage('Verify') {
            steps {
                catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                    sh '''
                    sleep 3

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
            echo 'PIPELINE SELESAI'
        }
    }
}