package main

import (
	"io/ioutil"
	"os"
	"strings"
)

func run_int_setup() {
	InfoLogger.Println("mkdir /var/eln_file_server")
	cmd("mkdir /var/eln_file_server")

	handle_error(os.Chdir("/var/eln_file_server"))
	InfoLogger.Println("Change dir /var/eln_file_server")
	cmd("apt-get update -y")
	cmd("apt-get install openssh-server nginx wget -y")
	cmd("systemctl enable ssh")
	cmd("systemctl start ssh")

	cmd("ufw allow 'NGINX full'")
	cmd("ufw allow 'OpenSSH'")
	cmd("echo 'y' | ufw enable")
	cmd("ufw reload")

	download()

	rewrite_file("/etc/ssh/sshd_config", func(line *string, eof bool) bool {
		if eof || strings.Index(*line, "AllowUsers") >= 0 {
			if equal := strings.Index(*line, "#"); equal == 0 {
				*line = string((*line)[1:])
			}

			all_user := remove_duplicate_str(append(get_all_user(), strings.Split(*line, " ")...))
			*line = strings.Join(append([]string{"AllowUsers"}, all_user...), " ")
			return false
		}

		return true
	})
	InfoLogger.Println("Update /etc/ssh/sshd_config")

	handle_error(ioutil.WriteFile("/var/eln_file_server/san.cnf", []byte(get_san_cnf()), 0644))

	cmd("openssl req -x509 -nodes -days 730 -newkey rsa:2048 -keyout server.key -out server.crt -config san.cnf")
	handle_error(os.Rename("/etc/nginx/sites-available/default", "/etc/nginx/sites-available/default.bck"))
	InfoLogger.Println("Update /etc/nginx/sites-available/default")
	defer func() {
		if !InitDone {
			InfoLogger.Println("Revert /etc/nginx/sites-available/default")
			_ = os.Remove("/etc/nginx/sites-available/default")
			_ = os.Rename("/etc/nginx/sites-available/default.bck", "/etc/nginx/sites-available/default")
		}
	}()

	handle_error(ioutil.WriteFile("/etc/nginx/sites-available/default", []byte(get_nginx_default()), 0644))

	handle_error(os.Rename("/usr/share/nginx/html/index.html", "/usr/share/nginx/html/index.html.bck"))
	InfoLogger.Println("Update /usr/share/nginx/html/index.html")
	defer func() {
		if !InitDone {
			InfoLogger.Println("Revert /usr/share/nginx/html/index.html")
			_ = os.Remove("/usr/share/nginx/html/index.html")
			_ = os.Rename("/usr/share/nginx/html/index.html.bck", "/usr/share/nginx/html/index.html")
		}
	}()

	handle_error(os.Rename("/var/eln_file_server/src/index.html", "/usr/share/nginx/html/index.html"))
	cmd("chmod 644 /usr/share/nginx/html/index.html")

	handle_error(ioutil.WriteFile("/var/eln_file_server/port_list.txt", []byte("8081"), 0764))
	handle_error(ioutil.WriteFile("/var/eln_file_server/project_list.txt", []byte(""), 0764))

	cmd("systemctl restart nginx.service")
	InfoLogger.Println("------INIT DONE!!!!------")
	InitDone = true
}
