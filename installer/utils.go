package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
)

func handle_error(err error) {
	if err != nil {
		ErrorLogger.Println(err)
		runtime.Goexit()
	}
}

func get_all_user() []string {
	files, err := ioutil.ReadDir("/home")
	handle_error(err)
	all_user := []string{}

	for _, file := range files {
		if file.IsDir() {
			all_user = append(all_user, file.Name())
		}
	}
	return all_user

}

func cmd(cmd string) {
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	InfoLogger.Println(string(out))
	handle_error(err)
}

func is_root() bool {
	currentUser, err := user.Current()
	handle_error(err)
	return currentUser.Username == "root"
}

func rewrite_file(file_path string, process func(x *string, eof bool) bool) {
	input, err := ioutil.ReadFile(file_path)
	handle_error(err)

	file_text := strings.Split(string(input), "\n")
	defer func() {
		err = ioutil.WriteFile(file_path, []byte(strings.Join(file_text, "\n")), 0644)
		handle_error(err)
	}()

	for i, line := range file_text {
		is_done := !process(&line, false)
		file_text[i] = line
		if is_done {
			return
		}

	}
	line := ""
	process(&line, true)
	if line != "" {
		file_text = append(file_text, line)
	}

}

func remove_duplicate_str(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func copy_file(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer func(source *os.File) {
		handle_error(source.Close())
	}(source)

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer func(destination *os.File) {
		handle_error(destination.Close())
	}(destination)
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func get_ip() string {
	if ip == "" {
		fmt.Print("Please enter Public reachable IP:")
		_, err := fmt.Scan(&ip)
		handle_error(err)
	}

	return ip
}
