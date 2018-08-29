def jobnameparts = JOB_NAME.tokenize('/') as String[]
def jobconsolename = jobnameparts[0]
def gopackage = "${jobconsolename}"

pipeline {
  agent any

  tools {
    go 'Go 1.10'
  }
  
  environment {
    GOPATH = "${WORKSPACE}"
    GOROOT = tool name: 'Go 1.10', type: 'go'
    PATH = "${GOPATH}/bin:$PATH"
  }

  stages {
    stage('Checkout') {
      steps {
        withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'Githbu', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD']]) {
          echo 'Checking out SCM'
          checkout scm

          sh 'git config credential.helper \'!f() { sleep 1; echo "username=${USERNAME}\npassword=${PASSWORD}"; }; f\''
          sh 'git fetch'
        }
      }
    }
      
    stage('Install dependencies') {
      steps {
        sh 'go version'
        sh "rm -rf src && mkdir -p src/${gopackage}"
        sh "ln -sf ${WORKSPACE}/* src/${gopackage}"
        sh "cd src/${gopackage} && go get -v"
      }
    }

    stage('Run tests') {
      steps {
        sh "cd src/${gopackage} && go test"
      }
    }
  
    stage('Build') {
      steps {
        sh "cd src/${gopackage} && go build -o ${jobconsolename}"
      }
    }

    stage('Deploy') {
      steps {
        script {
          tag = sh(returnStdout: true, script: "git tag --contains | head -1").trim()
          if (env.BRANCH_NAME != 'master' || !tag) {
            return
          }

          sh "mkdir -p /var/lib/jenkins/bindeploy/${jobconsolename}/${tag}"
          sh "cp -v src/${gopackage}/${jobconsolename} /var/lib/jenkins/bindeploy/${jobconsolename}/${tag}"
        }
      }
    }
  }

  post {
    always {
      sendSlackNotification(currentBuild)
      // cleanup current workspace
      deleteDir() 
    }
  }
}
