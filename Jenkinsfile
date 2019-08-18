pipeline {
    agent any
    stages{ 
        stage('Clone Game Client Repository'){
            steps {
                dir('client') {
                    checkout([$class: 'GitSCM', branches: [[name: '*/master']], doGenerateSubmoduleConfigurations: false, extensions: [], submoduleCfg: [], userRemoteConfigs: [[credentialsId: '76b15be1-6ed4-4a8d-b52b-4c7bff391b1b', url: 'https://github.com/BryJamin/PipeDream-Godot-Client.git']]])
                }
            }
        }
        stage('Export HTML5 Game into Cloned Game Client Repository'){
            steps {
                dir('client') {
                    //Build Godot HTML5 client and export in Docker
                    sh "docker rm -f pipedream-godot-client || true"
                    sh "docker build -t pipedream-godot-client ."

                    //Copy Exported Game From Docker Container to Host Directory Using a Temporary Container
                    sh "docker create --name temporary-export-container pipedream-godot-client"
                    sh "docker cp temporary-export-container:/pipedream-godot-client/exports ."
                    sh "docker rm temporary-export-container"
                }
            }
        }

        stage('Clone PipeDream Website Service Repository'){
            steps {
                dir('website') {
                    checkout([$class: 'GitSCM', branches: [[name: '*/master']], doGenerateSubmoduleConfigurations: false, extensions: [], submoduleCfg: [], userRemoteConfigs: [[credentialsId: '76b15be1-6ed4-4a8d-b52b-4c7bff391b1b', url: 'https://github.com/BryJamin/PipeDream-Website-Service.git']]])
                }
            }
        }

        stage('Copy Exported HTML5 Game into Static Website Folder'){
            steps {
                sh "ls -l"
                script {
                    if (fileExists('/client/exports')){
                        echo 'client exports exists'
                    } else {
                         echo 'client exports no existss'
                    }
                    if (fileExists('/website/static')){
                        echo '/website/static exists'
                    }
                    if (fileExists('/website/static/')){
                        echo '/website/static/ exists'
                    }
                }
                //Remove All contents of the Directory
                sh "rm /website/static/*"

                //Copy contents of export directory into Static directory
                sh "cp -a client/exports/. website/static"
            }
        }
        stage("Build PipeDream Website Image"){
            steps{
                dir('website') {
                sh "docker rm -f pipedream-website || true"
                sh "docker build -t pipedream-website ."
                }
            }
        }
        stage("Run PipeDream Website Image at port 80"){
            steps{
                sh "docker run -p 80:80 -d --name pipedream-website pipedream-website"
            }
        }
    }
    post { 
        always { 
            cleanWs()
        }
    }
}