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

        // UNIT TEST HARUS PASS
        stage('Unit Test') {
            steps {
                dir('PaymentService') {
                    sh '''
                    echo "Running Unit Test..."
                    go test ./...
                    '''
                }
            }
        }

        // VET HARUS PASS
        stage('Lint / Vet') {
            steps {
                dir('PaymentService') {
                    sh 'go vet ./...'
                }
            }
        }

        // BUILD HARUS BERHASIL
        stage('Build Image') {
            steps {
                sh 'docker build -t $IMAGE ./PaymentService'
            }
        }

        // FUNCTIONAL TEST REAL DB
        stage('Functional Test') {
            steps {
                sh '''
                echo "START DB"
                docker rm -f postgres-test || true
                docker run -d \
                  --name postgres-test \
                  -e POSTGRES_PASSWORD=123 \
                  -e POSTGRES_DB=testdb \
                  postgres

                echo "WAIT DB READY"
                until docker exec postgres-test pg_isready; do
                  sleep 1
                done

                echo "INIT DB"
                docker exec -i postgres-test psql -U postgres -d testdb <<EOF
                CREATE TABLE IF NOT EXISTS payments (
                  id INT PRIMARY KEY,
                  amount INT,
                  status TEXT
                );
                DELETE FROM payments;
                INSERT INTO payments VALUES (1, 50000, 'SUCCESS');
                EOF

                echo "START APP"
                docker rm -f test-payment || true
                docker run -d \
                  --name test-payment \
                  --link postgres-test \
                  -e DB_HOST=postgres-test \
                  -p 8082:8082 \
                  $IMAGE

                echo "WAIT API"
                for i in {1..15}; do
                  if curl -s http://localhost:8082; then
                    break
                  fi
                  sleep 1
                done

                echo "RUN TEST"
                cd PaymentService
                go test -run TestPaymentAPI_Success
                '''
            }
        }

        // PUSH HARUS BERHASIL
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
                sh 'echo "Deploy OK"'
            }
        }

        stage('Verify') {
            steps {
                sh 'echo "PIPELINE SUCCESS"'
            }
        }
    }
}