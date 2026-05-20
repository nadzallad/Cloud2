pipeline {
    agent any

    environment {
        DOCKER_IMAGE = "nadzalla/payment-service"
        IMAGE_TAG = "${BUILD_NUMBER}"
    }

    stages {

        stage('Checkout') {
            steps {
                git 'https://github.com/USERNAME/REPO.git'
            }
        }

        stage('Unit Test') {
            steps {
                sh '''
                cd PaymentService
                go test -v -run TestValidatePayment
                '''
            }
        }

        stage('Lint / Vet') {
            steps {
                sh '''
                cd PaymentService
                go vet ./...
                '''
            }
        }

        stage('Build Image') {
            steps {
                sh '''
                docker build -t $DOCKER_IMAGE:$IMAGE_TAG ./PaymentService
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
                      $DOCKER_IMAGE:$IMAGE_TAG

                    echo "⏳ Waiting for app to be ready..."

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
                      echo "❌ APP FAILED TO START"
                      docker logs test-payment
                      exit 1
                    fi

                    echo "✅ APP READY, RUN FUNCTIONAL TEST"

                    cd PaymentService
                    go test -v -run TestPaymentAPI || exit 1

                    echo "🧹 Cleanup"
                    docker rm -f test-payment
                    '''
                }
            }
        }

        stage('Push Image') {
            steps {
                withCredentials([string(credentialsId: 'dockerhub-pass', variable: 'PASS')]) {
                    sh '''
                    echo $PASS | docker login -u nadzalla --password-stdin
                    docker push $DOCKER_IMAGE:$IMAGE_TAG
                    '''
                }
            }
        }

        stage('Deploy') {
            steps {
                sh '''
                docker rm -f prod-payment || true

                docker run -d --name prod-payment \
                  -p 8083:8082 \
                  $DOCKER_IMAGE:$IMAGE_TAG
                '''
            }
        }

        stage('Verify') {
            steps {
                sh '''
                curl -X POST http://localhost:8083/payment \
                -H "Content-Type: application/json" \
                -d '{"amount":10000,"paid":10000}'
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