node {
  def root = tool name: 'Go 1.10', type: 'go'

  def jobnameparts = JOB_NAME.tokenize('/') as String[]
  def jobconsolename = jobnameparts[0]

  ws("${JENKINS_HOME}/jobs/${jobconsolename}/builds/${BUILD_ID}/") {
    withEnv([
      "GOPATH=${JENKINS_HOME}/jobs/${jobconsolename}/builds",
      "GOROOT=${root}",
      "PATH+GO=${root}/bin"
      ]) {
      env.PATH="${GOPATH}/bin:$PATH"
      
      withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'Githbu', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD']]) {
        stage('Checkout') {
          echo 'Checking out SCM'
          checkout scm

          sh 'git config credential.helper \'!f() { sleep 1; echo "username=${USERNAME}\npassword=${PASSWORD}"; }; f\''
          sh 'git fetch'
        }
      }
      
      stage('Pre Test') {
        echo 'Pulling Dependencies'

        sh "rm -rf ../src/${jobconsolename} && mkdir -p ../src/${jobconsolename} && cp -vr * ../src/${jobconsolename}"
    
        sh 'go version'
        sh "cd ../src/${jobconsolename} && go get -v"
      }
  
      stage('Build') {
        echo 'Building Executable'
      
        sh "go build -o ${jobconsolename}"
      }

      def tag = sh(returnStdout: true, script: "git tag --contains | head -1").trim()

      if (tag) {
        stage('Deploy') {
          if (env.BRANCH_NAME == 'master') {
            sh "mkdir -p /var/lib/jenkins/bindeploy/${jobconsolename}/${tag}"
            sh "cp -v ${jobconsolename} /var/lib/jenkins/bindeploy/${jobconsolename}/${tag}"
          }
        }
      }
    }
  }
}
