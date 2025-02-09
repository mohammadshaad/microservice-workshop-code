pipeline {
    agent any
    
    environment {
        DOCKER_REGISTRY = 'mohammadshaad'
    }
    
    stages {
        stage('Build User Service') {
            steps {
                dir('user-service') {
                    sh 'mvn clean package'
                    sh 'docker build -t ${DOCKER_REGISTRY}/user-service:${BUILD_NUMBER} .'
                }
            }
        }
        
        stage('Build Task Service') {
            steps {
                dir('task-service') {
                    sh 'docker build -t ${DOCKER_REGISTRY}/task-service:${BUILD_NUMBER} .'
                }
            }
        }
        
        stage('Build Frontend') {
            steps {
                dir('user-interface') {
                    sh 'docker build -t ${DOCKER_REGISTRY}/user-interface:${BUILD_NUMBER} .'
                }
            }
        }
        
        stage('Push Images') {
            steps {
                sh '''
                    docker push ${DOCKER_REGISTRY}/user-service:${BUILD_NUMBER}
                    docker push ${DOCKER_REGISTRY}/task-service:${BUILD_NUMBER}
                    docker push ${DOCKER_REGISTRY}/user-interface:${BUILD_NUMBER}
                '''
            }
        }
        
        stage('Deploy to Kubernetes') {
            steps {
                withKubeConfig([credentialsId: 'kubernetes-config']) {
                    sh '''
                        kubectl apply -f k8s/
                        kubectl set image deployment/user-service user-service=${DOCKER_REGISTRY}/user-service:${BUILD_NUMBER}
                        kubectl set image deployment/task-service task-service=${DOCKER_REGISTRY}/task-service:${BUILD_NUMBER}
                        kubectl set image deployment/user-interface user-interface=${DOCKER_REGISTRY}/user-interface:${BUILD_NUMBER}
                    '''
                }
            }
        }
    }
}
