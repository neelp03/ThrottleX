pipeline {
    agent any
    
    environment {
        GO111MODULE = 'on'
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Install Dependencies') {
            steps {
                sh 'go mod tidy'
            }
        }

        stage('Run Tests') {
            steps {
                sh 'go test ./... -v'
            }
        }

        stage('Build') {
            steps {
                sh 'go build -o throttlex cmd/throttlex/main.go'
            }
        }
    }

    post {
        always {
            archiveArtifacts artifacts: 'throttlex', allowEmptyArchive: true
            junit 'test-results.xml'
        }
    }
}
