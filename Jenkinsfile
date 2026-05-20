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
                sh '''
                docker build -t $IMAGE ./PaymentService
                '''
            }
        }

        stage('Functional Test') {
            steps {
                script {
                    sh '''
                    echo "🧹 Cleanup container lama"
                    docker rm -f test-payment || true

                    echo "🚀 Run container baru"
                    docker run -d --name test-payment \
                      -e DB_HOST=host.docker.internal \
                      -e DB_NAME=payment_db \
                      -e DB_PASS=admin123 \
                      -p 8082:8082 \
                      $IMAGE

                    echo "⏳ Waiting for app..."

                    sleep 3
                    docker logs test-payment

                    READY=0

                    for i in 1 2 3 4 5
                    do
                      STATUS=$(curl -s -o /dev/null -w "%{http_code}" \
                        -X POST http://host.docker.internal:8082/payment \
                        -H "Content-Type: application/json" \
                        -d '{"amount":1,"paid":1}')

                      echo "Attempt $i → Status: $STATUS"

                      if [ "$STATUS" = "200" ]; then
                        READY=1
                        break
                      fi

                      sleep 2
                    done

                    if [ $READY -eq 0 ]; then
                      echo "❌ APP FAILED"
                      docker logs test-payment
                      exit 1
                    fi

                    echo "✅ RUN GO TEST"

                    cd PaymentService
                    go test -v ./... || exit 1

                    echo "🧹 Cleanup"
                    docker rm -f test-payment
                    '''
                }
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
                  -e DB_NAME=payment_db \
                  -e DB_PASS=admin123 \
                  $IMAGE
                '''
            }
        }

        stage('Verify') {
            steps {
                sh '''
                echo "VERIFY API"

                sleep 3

                RESPONSE=$(curl -s -X POST http://localhost:8083/payment \
                  -H "Content-Type: application/json" \
                  -d '{"amount":10000,"paid":10000}')

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

    post {
        always {
            sh 'docker rm -f test-payment || true'
        }
    }
}