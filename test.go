package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

func FindAllString(yamlFile []byte) []byte {
	regex := "\\$\\{(.*?)\\}"
	re := regexp.MustCompile(regex)
	res := re.FindAllStringSubmatch(string(yamlFile),-1)
	for i := range res {
		fmt.Printf("Key: %s, Value: %s \n", res[i][0],res[i][1])
		env := os.Getenv(res[i][1])
		yamlFile = []byte(strings.ReplaceAll(string(yamlFile), res[i][0], env))
	}
	return yamlFile
}

func main()  {

	yamlFile, err := ioutil.ReadFile("test.yml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	yamlFile = FindAllString(yamlFile)

	f, err := os.Create("test2.yml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	write, err := f.Write(yamlFile)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(write)

	err = f.Close()
	if err != nil {
		fmt.Println(err)
	}

}
