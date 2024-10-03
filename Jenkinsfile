node('go1.23') {
	container('run'){
		def tag = ''

		try {
			stage('Checkout'){
					checkout scm
					notifyBitbucket()
					tag = sh(script: 'git tag -l --contains HEAD', returnStdout: true).trim()
			}

			stage('Run test'){
				sh('go test -v ./...')
			}

			if( tag != ''){
				strippedTag = tag.replaceFirst('v', '')
				stage('Build the application'){
					echo "Building with docker tag ${strippedTag}"
					sh('CGO_ENABLED=0 GOOS=linux go build')
				}

				stage('Generate docker image'){
					image = docker.build('fortnox/bitbucket-slack-bot:'+strippedTag, '--pull .')
				}

				stage('Push docker image'){
					docker.withRegistry("https://quay.io", 'docker-registry') {
						image.push()
					}
				}
			}

			currentBuild.result = 'SUCCESS'
		} catch(err) {
			currentBuild.result = 'FAILED'
			notifyBitbucket()
			throw err
		}

		notifyBitbucket()
	}
}

