//
// Created by burak on 5.06.2021.
//

#ifndef K8DEPLOYMENT_K8DEPLOYMENT_H
#define K8DEPLOYMENT_K8DEPLOYMENT_H

#include "../libs/log/log.h"
#include "../libs/kubernetes/libkubernetes.h"

#include <git2.h>
#include <git2/common.h>
#include <string.h>
#include <stdio.h>
#include <dirent.h>
#include <errno.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <stdbool.h>

bool check_environment(char *url, const char *token, char *branch);


//git
void check_error(int error_code, const char *action);
void use_remote(git_repository *repo,git_credential_userpass_payload user_pass,char *name,char *branch, char oid[41]);
int check_clone_or_pull_repository(char *dir_name, char *url, const char *token, char *branch, char *wait);

//deploy
int isFile(const char *path);
void listFilesRecursively(char *basePath);

#endif //K8DEPLOYMENT_K8DEPLOYMENT_H
