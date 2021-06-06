package main

import (
	"k8deployment/src"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	url := "https://github.com/bdagdeviren/k8s-yaml.git"
	token := "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"
	branch := "main"

	deploy := src.CloneOrPullRepository(url,branch,token)
	if deploy {
		err := filepath.Walk("deployment", func(path string, info os.FileInfo, err error) error {
			if !strings.Contains(path,".git") && strings.Contains(path,".yml") || strings.Contains(path,".yaml") {
				log.Println(path)
			}
			return nil
		})
		if err != nil {
			log.Fatalln(err.Error())
		}
	}


	//commitId,err := GetRemoteCommitId(url,branch,token)
	//if err != nil {
	//	log.Fatalln(commitId)
	//}

	//log.Printf("Last Commit Id: %s\n", commitId)
}
