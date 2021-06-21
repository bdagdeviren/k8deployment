#include "k8deployment.h"

void use_remote(git_repository *repo,git_credential_userpass_payload user_pass,char *name,char *branch, char oid[41])
{
    git_remote *remote = NULL;
    int error;
    const git_remote_head **refs;
    size_t refs_len, i;

    git_remote_callbacks callbacks = GIT_REMOTE_CALLBACKS_INIT;
    callbacks.credentials = git_cred_userpass;
    callbacks.payload = &user_pass;

    /* Find the remote by name */
    error = git_remote_lookup(&remote, repo, name);
    if (error < 0) {
        error = git_remote_create_anonymous(&remote, repo, name);
        if (error < 0)
            goto cleanup;
    }

    error = git_remote_connect(remote, GIT_DIRECTION_FETCH, &callbacks, NULL, NULL);
    if (error < 0)
        goto cleanup;

    if (git_remote_ls(&refs, &refs_len, remote) < 0)
        goto cleanup;

    for (i = 0; i < refs_len; i++) {
        if(strstr(refs[i]->name, branch) != NULL) {
            git_oid_fmt(oid, &refs[i]->oid);
        }
    }

    cleanup:
    git_remote_free(remote);
}

int check_clone_or_pull_repository(char *dir_name, char *url, const char *token, char *branch, char *wait){
    int error;

    int deploy = 0;

    char oid[GIT_OID_HEXSZ + 1] = {0};
    char local_oid[GIT_OID_HEXSZ + 1] = {0};

    git_reference *ref = NULL;
    git_repository *repo = NULL;
    git_remote *remote = NULL;

    git_credential_userpass_payload user_pass = {
            token, ""
    };

    const char *path = "deployment";

    git_libgit2_init();
    DIR* dir = opendir(dir_name);
    if (dir) {
        closedir(dir);
        // Getting Remote Commit Id
        error = git_repository_open(&repo, path);
        check_error(error, "opening repository");
        use_remote(repo,user_pass,url,branch,oid);
        // Getting Local Commit Id
        error = git_reference_dwim(&ref, repo, branch);
        check_error(error, "getting local commit id");
        git_oid_fmt(local_oid, git_reference_target(ref));
        if (strcmp(local_oid,oid) != 0){
            // Printing Local and Remote Commit Id
            log_info("Remote commit id:%s",oid);
            log_info("Local commit id:%s",local_oid);
            //Fetching Repository
            log_info("Fetching repository!");
            error = git_remote_lookup(&remote, repo, "origin");
            check_error(error, "remote lookup error");
            git_fetch_options opts = GIT_FETCH_OPTIONS_INIT;
            opts.callbacks.credentials = git_cred_userpass;
            opts.callbacks.payload = &user_pass;
            error = git_remote_fetch(remote,NULL,&opts,NULL);
            check_error(error, "fetch error");
            log_info("Fetching repository!");

            deploy = 1;
        } else {
            if(wait) {
                log_info("Nothing to change. Waiting %s second!", wait);
            }
        }

    } else if (ENOENT == errno) {
        log_info("Cloning git repository!");
        git_clone_options opts = GIT_CLONE_OPTIONS_INIT;
        opts.checkout_branch = branch;
        opts.fetch_opts.callbacks.credentials = git_cred_userpass;
        opts.fetch_opts.callbacks.payload = &user_pass;
        error = git_clone(&repo,url, "deployment", &opts);
        check_error(error, "clone error");
        log_info("Successfully cloned git repository!");

        deploy = 1;
    } else {
        log_info("Opening directory error!!");
    }

    git_repository_free(repo);
    git_reference_free(ref);
    git_remote_free(remote);
    git_libgit2_shutdown();
    repo = NULL; ref = NULL; remote = NULL;

    return deploy;
}