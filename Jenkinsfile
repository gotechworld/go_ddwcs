pipeline {
    agent any

    environment {
        APP_IMAGE_NAME = "go_ddwcs"
        APP_IMAGE_VERSION = "$BUILD_VERSION.$BUILD_NUMBER"
        DEPLOY_JOB = 'DDWCS-deploy'
    }

    parameters {
        string(name: 'BUILD_BRANCH', defaultValue: 'master', description: 'Build branch')
        string(name: 'BUILD_VERSION', defaultValue: '1', description: 'Build version')
        booleanParam(name: 'EXECUTE_TESTS', defaultValue: false, description: 'Execute unit test')
        booleanParam(name: 'EXECUTE_DEPLOY', defaultValue: false, description: 'Execute deploy')
    }

    stages {
        stage('Preparation') {
            steps{
                dir ('clone') {
                    // Get some code from a GitHub repository
                    git branch: params.BUILD_BRANCH, credentialsId: env.SSH_KEY_ID, url: scm.getUserRemoteConfigs()[0].getUrl()
                }

              sh "cp '${env.PIPELINE_HOME}/build/${env.JOB_NAME}/app.json' 'clone/config/app.json'"
            }
            post {
                success {
                    echo 'Successfully cloned!'
                }
                failure {
                    echo 'Failed to clone the repository.'
                }
            }
        }

        stage('Build') {
            steps {
                script {
                    dir ('clone') {
                        // Run the docker build
                        image = docker.build(env.APP_IMAGE_NAME + ':' + env.APP_IMAGE_VERSION, ' --build-arg SSH_KEY="$(cat ~/.ssh/'+ env.SSH_KEY_NAME +')" . -f Dockerfile --no-cache')
                    }
                }
            }
            post {
                success {
                    echo 'The image was successfully build.'
                }

                failure {
                    echo 'Failed to build the image'
                }
            }
        }

        stage('Test') {
            when {
                beforeAgent true
    			expression {
                    params.EXECUTE_TESTS
                }
            }
            steps {
                echo 'Execute Unit Tests'
            }
        }

        stage('Publish') {
            steps {
                echo 'Registry: ' + env.AMAZON_REGISTRY
                script {
                    docker.withRegistry('https://' + env.AMAZON_REGISTRY, env.AMAZON_REGISTRY_CREDENTIALS_ID) {
                        image.push(env.APP_IMAGE_VERSION)
                        image.push('latest')
                    }
                }
            }
            post {
                success {
                  echo 'The image with version ' + env.APP_IMAGE_VERSION + 'was pushed with success.'
                }

                failure {
                    echo 'Failed push the new image'
                }
            }
        }

        stage('Clean up') {
            steps {
                sh 'rm -rf clone scm' // clean repository folder 
                sh 'docker rmi ' + env.AMAZON_REGISTRY + '/' + env.APP_IMAGE_NAME + ':' + env.APP_IMAGE_VERSION + ' ' + env.AMAZON_REGISTRY + '/' + env.APP_IMAGE_NAME + ':latest' + ' ' + image.id + ' -f'
            }
            post {
                success {
                    echo 'Successfully cleaned up.'
                }

                failure {
                    echo 'Failed to clean the mess :('
                }
            }
        }

        stage('Deploy') {
            when {
                beforeAgent true
                expression {
                    params.EXECUTE_DEPLOY
                }
            }
            steps {
                build job: env.DEPLOY_JOB, parameters: [string(name: 'APP_IMAGE_VERSION', value: 'latest'), booleanParam(name: 'APP_FISRT_DEPLOY', value: false)]
            }

            post {
                success {
                    echo 'Successfully deployed.'
                }

                failure {
                    echo 'Failed to deploy.'
                }
            }
        }
    }
}
