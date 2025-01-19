// Constants
final String GIT_BRANCH = env.BRANCH_NAME
final String BLDSRV = 'bldsrv_prod'
final String PROJ_NAME = 'cube-cos-api'
final String SLACK_CHANNEL = "#${PROJ_NAME}"

// Environment setup
String BLDPTH = "/home/jenkins/workspace/${JOB_NAME}/${PROJ_NAME}"

env.getEnvironment().each { name, value ->
    println "Name: $name -> Value $value"
}

// Lock resources for exclusive access
timeout(time: 5, unit: 'MINUTES') {
    lock("${JOB_NAME}-${BLDSRV}") {
        node("${BLDSRV}") {
            docker.image('golang:1.22-alpine').inside('--user root -w /app') {
                ansiColor('xterm') {
                    stage('Checkout') {
                        echo "Using GIT_BRANCH: ${GIT_BRANCH}"

                        try {
                            checkoutCode(GIT_BRANCH, PROJ_NAME)
                        } catch (Exception e) {
                            echo 'Failed to download repository. Cleaning up and retrying...'
                            sh "sudo rm -rf ${PROJ_NAME}"
                            checkoutCode(GIT_BRANCH, PROJ_NAME)
                        }
                    }

                    dir(PROJ_NAME) {
                        stage('Prepare') {
                            echo 'Preparing the build environment...'
                            sh 'apk add --no-cache git openssh go-task'
                            sh "git config --global --add safe.directory \$(pwd)"
                        }

                        stage('Check') {
                            echo 'Running checks...'
                            sh 'go-task check'
                        }

                        stage('Build') {
                            echo 'Building the project...'
                            sh 'go-task build'
                        }
                    }
                }
            }
        }
    }
}

def checkoutCode(branchName, projectName) {
    checkout([
        $class: 'GitSCM',
        branches: [[name: branchName]],
        browser: [$class: 'GithubWeb', repoUrl: "https://github.com/bigstack-oss/${projectName}/"],
        doGenerateSubmoduleConfigurations: false,
        extensions: [
            [$class: 'GitLFSPull'],
            [$class: 'RelativeTargetDirectory', relativeTargetDir: projectName]
        ],
        userRemoteConfigs: [[
            url: "git@github.com:bigstack-oss/${projectName}.git",
            credentialsId: 'arashi-github-ssh-key'
        ]]
    ])
}
