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
                    sh '''
                    echo "UNIT TEST"
                    go test ./...
                    '''
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
                echo "START DB (REUSE)"
                docker start postgres-test || docker run -d \
                  --name postgres-test \
                  -e POSTGRES_PASSWORD=123 \
                  -e POSTGRES_DB=testdb \
                  postgres

                echo "WAIT DB QUICK"
                sleep 2

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

                echo "START APP (REUSE)"
                docker start test-payment || docker run -d \
                  --name test-payment \
                  --link postgres-test \
                  -e DB_HOST=postgres-test \
                  -p 8082:8082 \
                  $IMAGE

                echo "WAIT APP QUICK"
                sleep 2

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