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
                    sh 'go test -short ./...'
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
                echo "START DB"
                docker rm -f postgres-test || true
                docker run -d \
                  --name postgres-test \
                  -e POSTGRES_PASSWORD=admin123 \
                  -e POSTGRES_DB=payment_db \
                  postgres

                sleep 3

                echo "INIT DB"
                docker exec -i postgres-test psql -U postgres -d payment_db <<EOF
                CREATE TABLE IF NOT EXISTS payments (
                  id INT PRIMARY KEY,
                  amount INT,
                  status TEXT
                );
                DELETE FROM payments;
                INSERT INTO payments VALUES (1, 10000, 'PAID');
                EOF

                echo "START APP"
                docker rm -f test-payment || true
                docker run -d \
                  --name test-payment \
                  --link postgres-test \
                  -e DB_HOST=postgres-test \
                  -e DB_NAME=payment_db \
                  -e DB_PASS=admin123 \
                  -p 8082:8082 \
                  $IMAGE

                sleep 3

                echo "RUN TEST"
                cd PaymentService
                go test -run TestPaymentAPI_Success
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
                sh 'echo "DEPLOY OK"'
            }
        }

        stage('Verify') {
            steps {
                sh 'echo "PIPELINE SUCCESS"'
            }
        }
    }
}