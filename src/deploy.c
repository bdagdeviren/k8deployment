#include "k8deployment.h"

void listFilesRecursively(char *basePath)
{
    char path[1000];
    struct dirent *dp;
    DIR *dir = opendir(basePath);

    // Unable to open directory stream
    if (!dir)
        return;

    while ((dp = readdir(dir)) != NULL)
    {
        if (strcmp(dp->d_name, ".") != 0 && strcmp(dp->d_name, "..") != 0 && strcmp(dp->d_name, ".git") != 0) {
            // Construct new path from our base path
            strcpy(path, basePath);
            strcat(path, "/");
            strcat(path, dp->d_name);

            if (isFile(path)) {
                if (strstr(path, ".yml") != NULL || strstr(path, ".yaml") != NULL) {
                    struct apply_yaml_return applyYamlReturn;
                    applyYamlReturn = apply_yaml(path);
                    if (applyYamlReturn.r0 != 0) {
                        log_error(applyYamlReturn.r2);
                    } else {
                        log_info(applyYamlReturn.r1);
                    }
                }
            }

            listFilesRecursively(path);
        }
    }

    closedir(dir);
}

int isFile(const char *path) {
    struct stat statbuf;
    if (stat(path, &statbuf) != 0)
        return 0;
    return S_ISREG(statbuf.st_mode);
}