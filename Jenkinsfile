pipeline {
    agent any

    environment {
        PAYMENT_IMAGE = "payment-service"
        ORDER_IMAGE = "order-service"
        PICKUP_IMAGE = "pickup-service"
        SHIPMENT_IMAGE = "shipment-service"
        DELIVERY_IMAGE = "delivery-service"
        NOTIFICATION_IMAGE = "notification-service"
        TRACKING_IMAGE = "tracking-service"
        TAG = "latest"
    }

    stages {

        stage('Checkout Repo') {
            steps {
                deleteDir()
                git branch: 'main', url: 'https://github.com/nadzallad/Cloud2.git'
            }
        }

        // ========================
        // UNIT TEST (FIXED)
        // ========================
        stage('Unit Test') {
            steps {
                bat '''
                set TMP=C:\\Windows\\Temp
                set TEMP=C:\\Windows\\Temp
        
                go list ./... ^
                | findstr /V functional ^
                | findstr /V tests > packages.txt
        
                for /f %%i in (packages.txt) do go test %%i
                '''
            }
        }

        stage('Vet') {
            steps {
                bat 'go vet ./...'
            }
        }

        // ========================
        // BUILD DOCKER
        // ========================
        stage('Build Docker Images') {
            steps {
                bat '''
                cd PaymentService
                docker build -t payment-service:latest .
                cd ..

                cd OrderService
                docker build -t order-service:latest .
                cd ..

                cd PickupService
                docker build -t pickup-service:latest .
                cd ..

                cd ShipmentService
                docker build -t shipment-service:latest .
                cd ..

                cd DeliveryService
                docker build -t delivery-service:latest .
                cd ..

                cd NotificationService
                docker build -t notification-service:latest .
                cd ..

                cd TrackingService
                docker build -t tracking-service:latest .
                cd ..
                '''
            }
        }

        // ========================
        // FUNCTIONAL TEST
        // ========================
        stage('Functional Test') {
            steps {
                bat '''
                cd PaymentService
                start /b go run .
                cd ..

                cd OrderService
                start /b go run .
                cd ..

                cd PickupService
                start /b go run .
                cd ..

                cd ShipmentService
                start /b go run .
                cd ..

                cd DeliveryService
                start /b go run .
                cd ..

                cd NotificationService
                start /b go run .
                cd ..

                cd TrackingService
                start /b go run .
                cd ..

                timeout /t 10

                curl -X POST http://localhost:8081/payment ^
                -H "Content-Type: application/json" ^
                -d "{\\"amount\\":10000,\\"paid\\":10000}"

                curl -X POST http://localhost:8080/order ^
                -H "Content-Type: application/json" ^
                -d "{\\"user_id\\":1,\\"weight_kg\\":2,\\"distance_km\\":5,\\"base_price\\":10000}"

                curl -X POST http://localhost:8082/pickup ^
                -H "Content-Type: application/json" ^
                -d "{\\"order_id\\":\\"ORD1\\",\\"payment_status\\":\\"paid\\",\\"weight\\":2}"

                curl -X POST http://localhost:8083/notify ^
                -H "Content-Type: application/json" ^
                -d "{\\"message\\":\\"order created\\"}"

                curl -X POST http://localhost:8084/track ^
                -H "Content-Type: application/json" ^
                -d "{\\"order_id\\":\\"ORD1\\",\\"status\\":\\"shipped\\"}"
                '''
            }
        }

        // ========================
        // PUSH IMAGE
        // ========================
        stage('Push Images') {
            steps {
                bat '''
                docker tag order-service:latest ghryalvrt/order-service:latest
                docker push ghryalvrt/order-service:latest

                docker tag payment-service:latest ghryalvrt/payment-service:latest
                docker push ghryalvrt/payment-service:latest

                docker tag pickup-service:latest ghryalvrt/pickup-service:latest
                docker push ghryalvrt/pickup-service:latest

                docker tag shipment-service:latest selikakanajmi/shipment-service:latest
                docker push selikakanajmi/shipment-service:latest
                
                docker tag delivery-service:latest selikakanajmi/delivery-service:latest
                docker push selikakanajmi/delivery-service:latest

                docker tag notification-service:latest yourdockerhub/notification-service:latest
                docker push yourdockerhub/notification-service:latest

                docker tag tracking-service:latest yourdockerhub/tracking-service:latest
                docker push yourdockerhub/tracking-service:latest
                '''
            }
        }

        // ========================
        // DEPLOY
        // ========================
        stage('Deploy to Kubernetes') {
            steps {
                bat 'kubectl apply -f k8s/'
            }
        }

        stage('Verify Deployment') {
            steps {
                bat 'kubectl get pods && kubectl get svc'
            }
        }
    }
}
