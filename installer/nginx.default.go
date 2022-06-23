package main

import "fmt"

var nginx_default = `server {
  listen 80;
  server_name %s;
  # enforce https
  return 301 https://$server_name$request_uri;
}

server {
  listen 443;
  server_name %s;

  ssl on;
  ssl_certificate /var/eln_file_server/server.crt;
  ssl_certificate_key /var/eln_file_server/server.key;

  access_log /var/log/nginx/file_server.log;
  error_log /var/log/nginx/err_file_server.log;

#NEW INSTANCES
  location / {
    root /usr/share/nginx/html/;
  }
}`

func get_nginx_default() string {

	return fmt.Sprintf(nginx_default, get_ip(), get_ip())
}
