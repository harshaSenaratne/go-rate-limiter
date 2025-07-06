pipeline {
    agent any
    
    tools {
        go '1.24.4'
    }
    
    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }
        
        stage('Build Projects') {
            parallel {
                stage('Tollbooth Project') {
                    steps {
                        script {
                            dir('tollbooth') {
                                sh 'go mod tidy'
                                sh 'go test -v ./...'
                                sh 'go build -o tollbooth-app .'
                            }
                        }
                    }
                }
                
                stage('Token-Bucket Project') {
                    steps {
                        script {
                            dir('token-bucket') {
                                sh 'go mod tidy'
                                sh 'go test -v ./...'
                                sh 'go build -o token-bucket-app .'
                            }
                        }
                    }
                }
            }
        }
    }
    
    post {
        always {
            cleanWs()
        }
    }
}