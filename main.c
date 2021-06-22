#include <signal.h>

#include "libs/log/log.h"
#include "src/k8deployment.h"

void sigHandler (){
    log_info("Exiting programs!!!");
    exit(0);
}

int main() {


    int deploy;

    signal(SIGINT,sigHandler);
    signal(SIGTERM,sigHandler);

    char *url = getenv("URL");
    const  char *git_user = getenv("GIT_USER");
    const  char *token = getenv("TOKEN");
    char *branch = getenv("BRANCH");
    char *wait = getenv("WAIT");

    if (check_environment(url, token, branch)) {
        while (1) {
            deploy = check_clone_or_pull_repository("deployment", url, git_user, token, branch, wait);
            if (deploy) {
                listFilesRecursively("deployment");
                //log_info("%s", deploy);
            }
            if(wait) {
                sleep(atoi(wait));
            }else{
                exit(0);
            }
        }
    }
}
