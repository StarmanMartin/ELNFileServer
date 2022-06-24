package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func run_update() {
	download()
	project_list, err := ioutil.ReadFile("/var/eln_file_server/project_list.txt")
	handle_error(err)
	file_text := strings.Split(string(project_list), "\n")
	for _, user := range file_text {
		if user != "" {
			cmd(fmt.Sprintf("systemctl stop eln_instance_%s.service", user))
			copy_to_user(user)
			cmd(fmt.Sprintf("systemctl start eln_instance_%s.service", user))
		}
	}
}
