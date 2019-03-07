def slack = new org.daisho.Slack()

def packageName = "golang-template"
def appsList = "example-app other-app"
def packagePath = "src/github.com/nalej/${packageName}"

pipeline {
    agent { node { label 'golang' } }
    options {
        checkoutToSubdirectory("${packagePath}")
        buildDiscarder(logRotator(numToKeepStr: '10'))
    }

    stages {
        stage("Variable initialization") {
            steps { stepVariableInitialization packagePath }
        }
        stage("Git setup") {
            steps { container("golang") { stepGitSetup() } }
        }
        stage("Dependency download") {
            steps { container("golang") { stepGolangDependencyDownload packagePath } }
        }
        stage("Unit tests") {
            steps { container("golang") { stepGolangUnitTests packagePath } }
        }
    }

    post {
        success { script { slack.sendBuildNotification("success", "good") } }
        failure { script { slack.sendBuildNotification("failure", "danger") } }
        aborted { script { slack.sendBuildNotification("aborted", "warning") } }
    }
}
