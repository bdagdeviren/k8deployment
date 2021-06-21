#include "k8deployment.h"

bool check_environment(char *url, const char *token, char *branch) {
    bool run = false;

    if (url == NULL){
        log_error("Set URL environment!");
    }else{
        log_info("Url: %s",url);
    }

    if (token == NULL){
        log_error("Set TOKEN environment!");
    }else{
        log_info("Token: ************");
    }

    if (branch == NULL){
        log_error("Set BRANCH environment!");
    }else{
        log_info("Branch: %s",branch);
    }

    if (url && token && branch){
        run = true;
    }

    return run;
}

void check_error(int error_code, const char *action)
{
    const git_error *error = git_error_last();
    if (!error_code)
        return;
    log_error("Error %d %s - %s\n", error_code, action, (error && error->message) ? error->message : "???");
    exit(1);
}