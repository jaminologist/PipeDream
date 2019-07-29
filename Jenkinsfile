pipeline {
    agent any
    stages{ 
        stage('Clone Repository'){
            steps {
                checkout([$class: 'GitSCM', branches: [[name: '*/master']], doGenerateSubmoduleConfigurations: false, extensions: [], submoduleCfg: [], userRemoteConfigs: [[credentialsId: '76b15be1-6ed4-4a8d-b52b-4c7bff391b1b', url: 'https://github.com/BryJamin/PipeDream-Website-Service']]])
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