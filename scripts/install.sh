#! /bin/bash

installSoftware() {
    apt -qq -y install nginx
}

installSDA() {
    mkdir -p /var/www/sda
    curl -Lo- https://github.com/sunshineplan/sda/releases/download/v1.0/release.tar.gz | tar zxC /var/www/sda
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
    cp -s /var/www/sda/scripts/sda.conf /etc/nginx/conf.d
    sed -i "s/\$domain/$domain/" /var/www/sda/scripts/sda.conf
    service nginx reload
}

main() {
    read -p 'Please enter domain:' domain
    installSoftware
    installSDA
    writeLogrotateScrip
    setupNGINX
}

main
