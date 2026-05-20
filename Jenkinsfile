pipeline {
    agent {
        label 'master'
    }

    environment {
        TARGET_BRANCH = 'master'
    }

    stages {
        stage('Build') {
            steps {
                sh "docker pull goreleaser/goreleaser:v2.15.4"
                sh """
                  docker run -t --rm \\
                    -v ${WORKSPACE}:/go/src/cli -w /go/src/cli \\
                    --entrypoint sh \\
                    goreleaser/goreleaser:v2.15.4 \\
                    -c 'apk add --no-cache make && make build'
                """
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
                sh '''
                  aws-profile connectors-staging aws s3 cp ${WORKSPACE}/bin s3://smartling-connectors-releases/cli/ --acl public-read --exclude "*" --include "smartling-cli*" --include "smartling_*.deb" --include "smartling_*.rpm" --include "checksums.txt" --recursive
                '''
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
