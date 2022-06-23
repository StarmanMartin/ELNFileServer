package main

import "fmt"

var nginx_location = `

 location /%s/projects {

  proxy_pass http://127.0.0.1:%d;
  proxy_set_header Host $host;
  proxy_set_header X-Real-IP $remote_addr;
  proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
  proxy_set_header X-Forwarded-Proto https;
}
`

func get_nginx_location(project string) string {
	return fmt.Sprintf(nginx_location, project, port)
}
