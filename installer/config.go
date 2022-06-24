package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var config_file = `root_dir: "/home/%s/data"
webdav_prefix_url: "/%s/projects"
port: %d
logfile: "/home/%s/server.log"
host: "https://%s"
admin_password: %s`

func get_config(project, user string) string {
	project_list, err := ioutil.ReadFile("/var/eln_file_server/port_list.txt")
	handle_error(err)
	first_line := strings.Split(string(project_list), "\n")[0]

	port, err = strconv.ParseInt(first_line, 10, 64)
	handle_error(err)
	port += 1
	project_list = []byte(fmt.Sprintf("%d\n%s", port, project_list))
	handle_error(ioutil.WriteFile("/var/eln_file_server/port_list.txt", project_list, 764))

	var pass string
	fmt.Print("Please enter the new Admin password:")
	_, err = fmt.Scan(&pass)
	handle_error(err)
	return fmt.Sprintf(config_file, user, project, port, user, get_ip(), pass)

}
