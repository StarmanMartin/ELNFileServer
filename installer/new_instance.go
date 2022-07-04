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
	InfoLogger.Printf("| Instance Name: %s\n", user)

	cmd(fmt.Sprintf("useradd %s -m", user))
	on_cleanup(func() {
		if !NewDone {
			InfoLogger.Printf("delete user: %s\n", user)
			cmd(fmt.Sprintf("userdel %s", user))
			cmd(fmt.Sprintf("rm -R /home/%s", user))
		}
	})
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

	copy_to_user(user)

	handle_error(ioutil.WriteFile(fmt.Sprintf("/home/%s/server/config.yml", user), []byte(get_config(project, user)), 0764))

	cmd(fmt.Sprintf("chown %s:%s -R /home/%s", user, user, user))
	cmd(fmt.Sprintf("chmod +x -R /home/%s/server", user))
	cmd("/sbin/service sshd restart")

	InfoLogger.Printf("| SSH user: %s \n", user)
	cmd("cp /etc/nginx/sites-available/default /etc/nginx/sites-available/default.ist.bck")

	rewrite_file("/etc/nginx/sites-available/default", func(line *string, eof bool) bool {
		if strings.Index(*line, "#NEW INSTANCES") >= 0 {
			*line += get_nginx_location(project)
			return false
		}

		return true
	})

	on_cleanup(func() {
		if NewDone {
			cmd("rm /etc/nginx/sites-available/default.ist.bck")
		} else {
			InfoLogger.Printf("remove location from nginx\n")
			cmd("rm /etc/nginx/sites-available/default")
			cmd("cp /etc/nginx/sites-available/default.ist.bck /etc/nginx/sites-available/default")
			cmd("systemctl restart nginx")
		}

	})

	cmd("systemctl restart nginx")

	handle_error(ioutil.WriteFile(fmt.Sprintf("/etc/systemd/system/eln_instance_%s.service", user), []byte(get_service(user)), 766))
	on_cleanup(func() {
		if !NewDone {

			InfoLogger.Printf("remove eln_instance_%s.service\n", user)
			cmd(fmt.Sprintf("rm /etc/systemd/system/eln_instance_%s.service", user))
		}
	})
	InfoLogger.Printf("| systemctl restart eln_instance_%s.service\n", user)
	cmd(fmt.Sprintf("systemctl enable eln_instance_%s.service", user))
	on_cleanup(func() {
		if !NewDone {
			InfoLogger.Printf("systemctl disable & stop eln_instance_%s.service\n", user)
			cmd(fmt.Sprintf("systemctl stop eln_instance_%s.service", user))
			cmd(fmt.Sprintf("systemctl disable eln_instance_%s.service", user))
		}
	})
	cmd(fmt.Sprintf("systemctl start eln_instance_%s.service", user))

	InfoLogger.Printf("| Server address: https://%s/%s/projects \n", get_ip(), project)
	InfoLogger.Println("------NEW INSTANCE DONE!!!!------")
	NewDone = true
}
