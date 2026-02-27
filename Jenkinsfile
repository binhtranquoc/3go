def IMAGE_NAME = "octotechvn/gogogo-api"
def GITHUB_URL = "https://github.com/teodevlor/3go.git"

def BRANCH = env.BRANCH_NAME

def NODE_LABEL
def STACK_NAME
def ENVIRONMENT

if (BRANCH == "dev") {
    ENVIRONMENT = "dev"
    STACK_NAME = "gogogo-dev"
    NODE_LABEL = "dev-server"
} 
else if (BRANCH == "master") {
    ENVIRONMENT = "prod"
    STACK_NAME = "gogogo-prod"
    NODE_LABEL = "prod-server"
} 
else {
    error("Branch ${BRANCH} is not deployable!")
}

node(NODE_LABEL) {

    stage("Check Environment") {
        sh """
            echo "Branch: ${BRANCH}"
            echo "Environment: ${ENVIRONMENT}"
            echo "Stack: ${STACK_NAME}"
            echo "Running on node:"
            hostname
            whoami
            pwd
        """
    }

}