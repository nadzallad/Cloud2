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

        stage('Unit Test') {
            steps {
                dir('PaymentService') {
                    sh 'go test -v -run TestValidatePayment ./...'
                }
            }
        }

        stage('Lint / Vet') {
            steps {
                dir('PaymentService') {
                    sh 'go vet ./...'
                }
            }
        }

        stage('Build Image') {
            steps {
                sh 'docker build -t $IMAGE ./PaymentService'
            }
        }

        stage('Functional Test') {
            steps {
                sh '''
                docker rm -f test-payment || true

                docker run -d \
                  --name test-payment \
                  -e DB_HOST=host.docker.internal \
                  -e DB_NAME=payment_db \
                  -e DB_PASS=admin123 \
                  -p 8082:8082 \
                  $IMAGE

                echo "WAIT APP"
                until curl -s http://localhost:8082/payment; do
                  sleep 1
                done

                cd PaymentService
                go test -v -run TestPaymentAPI ./...
                '''
            }
        }

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

        stage('Deploy') {
            steps {
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

        stage('Verify') {
            steps {
                sh '''
                echo "VERIFY API"

                until curl -s http://localhost:8083/payment; do
                  sleep 1
                done

                RESPONSE=$(curl -s -X POST http://localhost:8083/payment \
                  -H "Content-Type: application/json" \
                  -d '{
                    "order_id":2,
                    "amount":10000,
                    "paid":10000,
                    "payment_method":"BANK_TRANSFER"
                  }')

                echo "Response: $RESPONSE"

                if echo "$RESPONSE" | grep -q PAID; then
                  echo "SUCCESS"
                else
                  echo "FAILED"
                  exit 1
                fi
                '''
            }
        }
    }
}