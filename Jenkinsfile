pipeline {
    agent any
    stages{ 
        stage('Clone Repository'){
            steps {
                checkout([$class: 'GitSCM', branches: [[name: '*/master']], doGenerateSubmoduleConfigurations: false, extensions: [], submoduleCfg: [], userRemoteConfigs: [[credentialsId: '76b15be1-6ed4-4a8d-b52b-4c7bff391b1b', url: 'https://github.com/BryJamin/PipeDream-Lobby-Service']]])
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
                sh "docker run -p 5080:5080 -d --name pipedream-test pipedream-test"
            }
        }
    }
}