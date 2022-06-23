package main

import "fmt"

var service = `[Unit]
Description = file server instance %s

[Service]
WorkingDirectory = /home/%s/server
ExecStart = /home/%s/server/eln_file_server

[Install]
WantedBy = multi-user.target`

func get_service(user string) string {
	return fmt.Sprintf(service, user, user, user)
}
