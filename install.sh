#! /bin/bash

installSoftware() {
    apt -qq -y install nginx
    apt -qq -y -t $(lsb_release -sc)-backports install golang-go
}

installSDA() {
    curl -Lo- https://github.com/sunshineplan/sda-go/archive/v1.0.tar.gz | tar zxC /var/www
    mv /var/www/sda-go* /var/www/sda-go
    cd /var/www/sda-go
    go build
}

configSDA() {
    read -p 'Please enter unix socket(default: /var/www/sda-go/sda-go.sock): ' unix
    [ -z $unix ] && unix=/var/www/sda-go/sda-go.sock
    read -p 'Please enter host(default: 127.0.0.1): ' host
    [ -z $host ] && host=127.0.0.1
    read -p 'Please enter port(default: 12345): ' port
    [ -z $port ] && port=12345
    read -p 'Please enter log path(default: /var/log/app/sda-go.log): ' log
    [ -z $log ] && log=/var/log/app/sda-go.log
    mkdir -p $(dirname $log)
    sed "s,\$unix,$unix," /var/www/sda-go/config.ini.default > /var/www/sda-go/config.ini
    sed -i "s,\$log,$log," /var/www/sda-go/config.ini
    sed -i "s/\$host/$host/" /var/www/sda-go/config.ini
    sed -i "s/\$port/$port/" /var/www/sda-go/config.ini
}

setupsystemd() {
    cp -s /var/www/sda-go/sda-go.service /etc/systemd/system
    systemctl enable sda-go
    service sda-go start
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
    cp -s /var/www/sda-go/sda-go.conf /etc/nginx/conf.d
    sed -i "s/\$domain/$domain/" /var/www/sda-go/sda-go.conf
    sed -i "s,\$unix,$unix," /var/www/sda-go/sda-go.conf
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