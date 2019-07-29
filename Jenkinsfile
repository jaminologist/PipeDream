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
                sh "docker rm -f pipedream-website || true"
                sh "docker build -t pipedream-website ."
            }
        }
        stage("Run Image"){
            steps{
                sh "docker run -p 17700:17700 -d --name pipedream-website pipedream-website"
            }
        }
    }
}