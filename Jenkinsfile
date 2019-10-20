pipeline {
    agent any
    stages{ 
        stage('Clone Repository'){
            steps {
                checkout scm
            }
        }
        stage("Build Image"){
            steps{
                sh "docker rm -f pipedream-test || true"
                sh "docker build -t pipedream-test ."
            }
        }
        stage("Run Image"){
            steps{
                sh "docker run -p 5080:5080 -p 80:80 -p 443:443 -e ENVIRONMENT=-env=production -d --name pipedream-test pipedream-test"
            }
        }
    }
}