<h1>How to setup Ubuntu server:</h1>

Download and run the [Installer](https://github.com/StarmanMartin/ELNFileServer/raw/main/installer/bin/installer) as root (hint: sudo su)<br>

```
mkdir ~/eln_installer
cd ~/eln_installer
wget https://github.com/StarmanMartin/ELNFileServer/raw/main/installer/bin/installer
sudo chmod +x installer
sudo su -c ./installer
```

<h1>How to setup by hand!</h1>
<p>The following part is only fo the case it is the first instance of the ELN file receiver on this machine:</p>
<p>Hint: [A-z] are variables and have to be set by you.</p>
<h2 id="section1">First Setup</h2>
<p>If one ELN file receiver is already running on your system jump to <a href="#section2">Add ELN file receiver instance</a></p>
<p>The following first steps are for a ubuntu system with a fixed IP (without URL). It also shows how to quickly create a self-signed TSL certificate. </p>
<ol>
    <li>Install nginx: <br>sudo apt-get update & sudo apt-get install nginx</li>
    <li>Check nginx: <br>nginx -v<br>sudo systemctl status nginx </li>
    <li>Start nginx: <br>sudo systemctl start nginx<br>sudo systemctl enable nginx</li>
</ol>
    <p>[Optional] add a firewall:</p>
<ol>
    <li>Run to check: <br> sudo ufw app list</li>
    <li>Make sure that: <br> sudo ufw allow 'NGINX full' <br> sudo ufw allow 'OpenSSH'</li>
    <li>Reload firewall: <br> sudo ufw enable & sudo ufw reload </li>
</ol>
    <p>[Optional] add TSL</p>
<ol>
    <li>Make a directory for certifications: <br> mkdir ~/certs & cd ~/certs </li>
    <li>Make a config file: <br> nano san.cnf <br>
        and add the following content:
        <hr>
        [req]<br>
        default_bits  = 2048<br>
        distinguished_name = req_distinguished_name<br>
        req_extensions = req_ext<br>
        x509_extensions = v3_req<br>
        prompt = no<br>
        [req_distinguished_name]<br>
        countryName = []<br>
        stateOrProvinceName = []<br>
        localityName = []<br>
        organizationName = []<br>
        OU=[]<br>
        commonName = [IP_ADDRESS]<br>
        [req_ext]<br>
        subjectAltName = @alt_names<br>
        [v3_req]<br>
        subjectAltName = @alt_names<br>
        [alt_names]<br>
        IP.1 = [IP_ADDRESS]<hr>
        FILE: ~/certs/san.cnf<hr>
    </li>
    <li>Then run: <br> openssl req -x509 -nodes -days 730 -newkey rsa:2048 -keyout key.pem -out cert.pem -config san.cnf</li>
</ol>
    <p>Edit nginx routing</p>
<ol>
    <li>Make backup: <br> sudo cp /etc/nginx/sites-enabled/default /etc/nginx/sites-enabled/default.back</li>
    <li>Edit /etc/nginx/sites-enabled/default: <br> sudo nano /etc/nginx/sites-enabled/default <hr>
        server {<br>
        listen        80;<br>
        server_name   [IP_ADDRESS];<br>
        # enforce https<br>
        return        301 https://$server_name$request_uri;<br>
        }<br><br>
        server {<br>
        listen                   443;<br>
        server_name              [IP_ADDRESS];<br>
            <br>
        ssl                     on;<br>
        ssl_certificate         /home/[LINUX_USER]/certs/server.crt;<br>
        ssl_certificate_key     /home/[LINUX_USER]/certs/server.key;<br>
                    <br>
        access_log              /var/log/nginx/file_server.log;<br>
        error_log               /var/log/nginx/err_file_server.log;<br>
                    <br>
        #NEW INSTANCES <br>
        location / {<br>
        root /usr/share/nginx/html/;<br>
        }<br>
        }<hr>
        FILE: /etc/nginx/sites-enabled/default
        <hr>
    </li>
    <li>Replace the /usr/share/nginx/html/index.html with the content of the index html of the project.</li>
</ol>
<h2 id="section2">Add ELN file receiver instance</h2>
<p>If no ELN file receiver is already running on your system jump to <a href="#section1">First Setup</a></p>
<ol>
    <li>Add a new Linux user: <br> https://docs.oracle.com/en/cloud/cloud-at-customer/occ-get-started/add-ssh-enabled-user.html</li>
    <li>change user: <br> sudo su [USERNAME] & cd ~ </li>
    <li>Make a server directory and a data directory <br> mkdir server & mkdir data & cd server</li>
    <li>Copy from <a href="https://github.com/StarmanMartin/ELNFileServer/releases/download/Prerelease/ELNFileServer_v0.1.tar.gz">Prerelease/ELNFileServer</a><br>
        curl -fsSL https://github.com/StarmanMartin/ELNFileServer/releases/download/Prerelease/ELNFileServer_v0.1.tar.gz --output ./ELNFileServer_v0.1.tar.gz <br>
        tar -xf ELNFileServer_v0.1.tar.gz<br>
    </li>
    <li>Make a config file: <br> nano config.yml<hr>
        root_dir: "/home/[NEW_USER]/data"<br>
        webdav_prefix_url: "/[NEW_ORGANIZATION_NAME]/projects"<br>
        port: [UNUSED_PORT]<br>
        logfile: "server.log"<br>
        host: "https://[IP_ADDRESS]"<br>
        admin_password: [NEW_ADMIN_PASSOWRD]<hr>
        FILE: ~/config.yml
        <hr>
    </li>
    <li>Switch to sudo user: <br> exit </li>
    <li>sudo chmod +x /home/[NEW_USER]/server/eln_file_server</li>
    <li>Edit /etc/nginx/sites-enabled/default. Add new routing below the commend: #NEW INSTANCES<br>
        sudo nano /etc/nginx/sites-enabled/default
        <hr>
        ...<br>
        #NEW INSTANCES <br>
        <br>
        location /[NEW_ORGANIZATION_NAME]/projects {<br><br>
        proxy_pass              http://127.0.0.1:[UNUSED_PORT];<br>
        proxy_set_header        Host $host;<br>
        proxy_set_header        X-Real-IP $remote_addr;<br>
        proxy_set_header        X-Forwarded-For $proxy_add_x_forwarded_for;<br>
        proxy_set_header        X-Forwarded-Proto https;<br><br>
        }<br>
        ...<hr>
        FILE: /etc/nginx/sites-enabled/default
        <hr>
    </li>
    <li>Restart nginx: <br> sudo systemctr restart nginx</li>
    <li>Finally create systemd service: <br>sudo nano /etc/systemd/system/[SERVICE_NAME].service<hr>
        [Unit]<br>
        Description = file server instance [SERVICE_NAME]<br>
        <br>
        [Service]<br>
        WorkingDirectory = /home/[NEW_USER]/server<br>
        ExecStart = /home/[NEW_USER]/server/eln_file_server<br>
        <br>
        [Install]<br>
        WantedBy = multi-user.target<hr>
        FILE: /etc/systemd/system/[SERVICE_NAME].service
        <hr>
    </li>
    <li>Enable and start server: <br>sudo systemctl enable [SERVICE_NAME].service & sudo systemctl start [SERVICE_NAME].service</li>
    <li>Check server: <br>sudo systemctl statu [SERVICE_NAME].service</li>
</ol>