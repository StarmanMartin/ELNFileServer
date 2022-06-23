package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func check_project(project *string) error {
	project_list, err := ioutil.ReadFile("/var/eln_file_server/project_list.txt")
	handle_error(err)
	file_text := strings.Split(string(project_list), "\n")
	for _, line := range file_text {
		if line == *project {
			return errors.New("Project already exists!")
		}
	}
	file_text = append(file_text, *project)
	handle_error(ioutil.WriteFile("/var/eln_file_server/project_list.txt", []byte(strings.Join(file_text, "\n")), 764))
	return nil
}

func run_insatnce_setup() {
	var pk, project string
	fmt.Print("Please enter instance name (only letters):")
	_, err := fmt.Scan(&project)
	handle_error(err)
	handle_error(check_project(&project))
	user := strings.ToLower(project)
	project = user

	cmd(fmt.Sprintf("useradd %s -m", user))
	handle_error(os.Mkdir(fmt.Sprintf("/home/%s/.ssh", user), 600))

	fmt.Print("Enter your the SSH public key (Hint nano id_dsa.pub): ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	pk = scanner.Text()

	handle_error(err)

	handle_error(ioutil.WriteFile(fmt.Sprintf("/home/%s/.ssh/authorized_keys", user), []byte(pk), 0600))

	handle_error(os.Mkdir(fmt.Sprintf("/home/%s/server", user), 0764))
	handle_error(os.Mkdir(fmt.Sprintf("/home/%s/server/views", user), 0764))
	handle_error(os.Mkdir(fmt.Sprintf("/home/%s/data", user), 764))

	rewrite_file("/etc/ssh/sshd_config", func(line *string, eof bool) bool {
		if eof || strings.Index(*line, "AllowUsers") >= 0 {
			if equal := strings.Index(*line, "#"); equal == 0 {
				*line = string((*line)[1:])
			}

			all_user := remove_duplicate_str(append(get_all_user(), strings.Split(*line, " ")[1:]...))
			*line = strings.Join(append([]string{"AllowUsers"}, all_user...), " ")
			return false
		}

		return true
	})

	_, err = copy_file("/var/eln_file_server/eln_file_server", fmt.Sprintf("/home/%s/server/eln_file_server", user))
	handle_error(err)

	_, err = copy_file("/var/eln_file_server/views/new_user.gtpl", fmt.Sprintf("/home/%s/server/views/new_user.gtpl", user))
	handle_error(err)

	handle_error(ioutil.WriteFile(fmt.Sprintf("/home/%s/server/config.yml", user), []byte(get_config(project, user)), 0764))

	cmd(fmt.Sprintf("chown %s:%s -R /home/%s", user, user, user))
	cmd("/sbin/service sshd restart")

	rewrite_file("/etc/nginx/sites-available/default", func(line *string, eof bool) bool {
		if strings.Index(*line, "#NEW INSTANCES") >= 0 {
			*line += get_nginx_location(project)
			return false
		}

		return true
	})

	cmd("systemctl restart nginx")

	handle_error(ioutil.WriteFile(fmt.Sprintf("/etc/systemd/system/eln_instance_%s.service", user), []byte(get_service(user)), 766))
	InfoLogger.Printf("systemctl enable eln_instance_%s.service\n", user)
	cmd(fmt.Sprintf("systemctl enable eln_instance_%s.service", user))
	cmd(fmt.Sprintf("systemctl start eln_instance_%s.service", user))

	InfoLogger.Println("------NEW INSTANCE DONE!!!!------")
}
