#! /bin/bash

installSoftware() {
    apt -qq -y install nginx
    apt -qq -y -t $(lsb_release -sc)-backports install golang-go
}

installSDA() {
    curl -Lo- https://github.com/sunshineplan/sda/archive/v1.0.tar.gz | tar zxC /var/www
    mv /var/www/sda* /var/www/sda
    cd /var/www/sda/webapp
    go build -ldflags "-s -w" -o sda
}

configSDA() {
    read -p 'Please enter unix socket(default: /run/sda.sock): ' unix
    [ -z $unix ] && unix=/run/sda.sock
    read -p 'Please enter host(default: 127.0.0.1): ' host
    [ -z $host ] && host=127.0.0.1
    read -p 'Please enter port(default: 12345): ' port
    [ -z $port ] && port=12345
    read -p 'Please enter log path(default: /var/log/app/sda.log): ' log
    [ -z $log ] && log=/var/log/app/sda.log
    mkdir -p $(dirname $log)
    sed "s,\$unix,$unix," /var/www/sda/webapp/config.ini.default > /var/www/sda/webapp/config.ini
    sed -i "s,\$log,$log," /var/www/sda/webapp/config.ini
    sed -i "s/\$host/$host/" /var/www/sda/webapp/config.ini
    sed -i "s/\$port/$port/" /var/www/sda/webapp/config.ini
}

setupsystemd() {
    cp -s /var/www/sda/webapp/sda.service /etc/systemd/system
    systemctl enable sda
    service sda start
}

writeLogrotateScrip() {
    if [ ! -f '/etc/logrotate.d/app' ]; then
	cat >/etc/logrotate.d/app <<-EOF
		/var/log/app/*.log {
		    copytruncate
		    rotate 12
		    compress
		    delaycompress
		    missingok
		    notifempty
		}
		EOF
    fi
}

setupNGINX() {
    cp -s /var/www/sda/webapp/sda.conf /etc/nginx/conf.d
    sed -i "s/\$domain/$domain/" /var/www/sda/webapp/sda.conf
    sed -i "s,\$unix,$unix," /var/www/sda/webapp/sda.conf
    service nginx reload
}

main() {
    read -p 'Please enter domain:' domain
    installSoftware
    installSDA
    configSDA
    setupsystemd
    writeLogrotateScrip
    setupNGINX
}

main