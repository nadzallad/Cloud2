        steps {
                catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                    sh '''
                    echo "Run DB container..."
                    docker rm -f postgres-test || true
                    docker run -d \
                    --name postgres-test \
                    -e POSTGRES_PASSWORD=123 \
                    -e POSTGRES_DB=testdb \
                    --health-cmd="pg_isready" \
                    --health-interval=2s \
                    postgres

                    echo "Waiting DB ready..."
                    until [ "$(docker inspect -f {{.State.Health.Status}} postgres-test)" = "healthy" ]; do
                    sleep 1
                    done

                    echo "Run app container..."
                    docker rm -f test-payment || true
                    docker run -d -p 8082:8082 --name test-payment $IMAGE

                    echo "Waiting API..."
                    until curl -s http://localhost:8082; do
                    sleep 1
                    done

                    echo "Run Functional Test..."
                    cd PaymentService
                    go test -run TestPaymentAPI_Success

                    echo "Cleanup..."
                    docker rm -f test-payment
                    docker rm -f postgres-test
                    '''
                }
            }
        }
