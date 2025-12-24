pipeline {
    agent {
        label 'master'
    }

    environment {
        TARGET_BRANCH = 'master'
    }

    stages {
        stage('Run Tests') {
            steps {
                sh "docker pull golang"
                sh """docker run -t --rm -v \${WORKSPACE}:/go/src/cli -w /go/src/cli golang bash -c '
                    go install github.com/jstemmer/go-junit-report/v2@latest && \\
                    go test -v ./cmd/... 2>&1 | tee /tmp/test-output.txt | /go/bin/go-junit-report -set-exit-code > test-results.xml
                '"""
            }
            post {
                always {
                    junit 'test-results.xml'
                }
            }
        }

        stage('Build') {
            steps {
                sh "docker run -t --rm -v ${WORKSPACE}:/go/src/cli -w /go/src/cli golang make"
            }
        }

        stage('Release?') {
            agent none
            when {
                branch env.TARGET_BRANCH
            }
            steps {
                timeout(time: 1, unit: 'HOURS') {
                    input 'Release to PROD?'
                }
            }
        }

        stage('Upload to public S3') {
            when {
                branch env.TARGET_BRANCH
            }
            steps {
                sh "aws-profile connectors-staging aws s3 cp ${WORKSPACE}/bin s3://smartling-connectors-releases/cli/ --recursive --acl public-read"
            }
        }

        stage('Generate Packages') {
            when {
                branch env.TARGET_BRANCH
            }
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
