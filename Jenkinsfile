#!/usr/bin/env groovy

pipeline {
    agent any

    stages {
        stage('Build') {
            
            environment {
                GOPATH='/data/opt/server/jenkins/jenkins/.gvm/pkgsets/go1.16/global'
                GOROOT='/data/opt/server/jenkins/jenkins/.gvm/gos/go1.16'
                GOBIN='/data/opt/server/jenkins/jenkins/.gvm/gos/go1.16/bin'
                GVM_ROOT='/data/opt/server/jenkins/jenkins/.gvm'
                GVM_VERSION='1.0.22'
            }
    
            steps {
                sh 'rm -r build'
                sh 'make config'
                sh 'make build'
                archiveArtifacts artifacts: 'build/**/*', fingerprint: true
            }
        }
        stage('Deploy') {
            when {
                expression {
                    currentBuild.result == null || currentBuild.result == 'SUCCESS'
                }
            }
            steps {
                sh 'make publish'
                sh 'make restart'
            }
        }
    }
}
