pipeline {
    agent {
        label 'master'
    }

    stages {
        stage('Build') {
            steps {
                sh "docker pull golang"
                sh "docker run -t --rm -v ${WORKSPACE}:/go/src/cli -w /go/src/cli golang make"
                sh "aws-profile connectors-staging aws s3 cp ${WORKSPACE}/bin s3://smartling-connectors-releases/cli/ --recursive --acl public-read"
            }
        }

        stage('Generate Packages') {
            steps {
                sh "docker run -t --rm -v ${WORKSPACE}:/go/src/cli -w /go/src/cli gvangool/rpmbuilder:centos7 bash -c 'make rpm'"
                // TODO : Replace with special docker image
                sh "docker run -t --rm -v ${WORKSPACE}:/go/src/cli -w /go/src/cli debian bash -c 'apt-get update && apt-get install -y make git && make deb'"
            }
        }
    }

    post {
        unstable {
            slackSend (
                    channel: "#emergency-connectors",
                    color: 'bad',
                    message: "Tests failed: <${env.RUN_DISPLAY_URL}|${env.JOB_NAME} #${env.BUILD_NUMBER}>"
            )
        }

        failure {
            slackSend (
                    channel: "#emergency-connectors",
                    color: 'bad',
                    message: "Build of <${env.RUN_DISPLAY_URL}|${env.JOB_NAME} #${env.BUILD_NUMBER}> is failed!"
            )
        }
    }
}
