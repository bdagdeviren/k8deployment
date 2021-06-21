package main

import "C"
import (
	"github.com/johandry/klient"
	"io/ioutil"
)

//export create_yaml
func create_yaml(path *C.char) (rc int,result *C.char,errStr *C.char) {
	file, err := ioutil.ReadFile(C.GoString(path))
	if err != nil {
		return -1, nil, C.CString(err.Error())
	}

	c := klient.New("", "")
	err = c.Create(file)
	if err != nil {
		return -1, nil, C.CString(err.Error())
	}

	resultFile := "Successfully create "+C.GoString(path)+" file!"
	return 0,C.CString(resultFile),nil
}

//export update_yaml
func update_yaml(path *C.char) (rc int,result *C.char,errStr *C.char) {
	file, err := ioutil.ReadFile(C.GoString(path))
	if err != nil {
		return -1, nil, C.CString(err.Error())
	}

	c := klient.New("", "")
	err = c.Replace(file)
	if err != nil {
		return -1, nil, C.CString(err.Error())
	}

	resultFile := "Successfully update "+C.GoString(path)+" file!"
	return 0,C.CString(resultFile),nil
}

//export delete_yaml
func delete_yaml(path *C.char) (rc int,result *C.char,errStr *C.char) {
	file, err := ioutil.ReadFile(C.GoString(path))
	if err != nil {
		return -1, nil, C.CString(err.Error())
	}

	c := klient.New("", "")
	err = c.Delete(file)
	if err != nil {
		return -1, nil, C.CString(err.Error())
	}

	resultFile := "Successfully delete "+C.GoString(path)+" file!"
	return 0,C.CString(resultFile),nil
}

//export apply_yaml
func apply_yaml(path *C.char) (rc int,result *C.char,errStr *C.char) {
	file, err := ioutil.ReadFile(C.GoString(path))
	if err != nil {
		return -1, nil, C.CString(err.Error())
	}

	c := klient.New("", "")
	err = c.Apply(file)
	if err != nil {
		return -1, nil, C.CString(err.Error())
	}

	resultFile := "Successfully apply "+C.GoString(path)+" file!"
	return 0,C.CString(resultFile),nil
}

func main() {}