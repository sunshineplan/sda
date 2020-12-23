#! /bin/bash

installSoftware() {
    apt -qq -y install nginx
}

installSDA() {
    curl -Lo /var/www/sda/build/bundle.js https://github.com/sunshineplan/sda/releases/download/v1.0/bundle.js --create-dirs
    cd /var/www/sda
    curl -LO https://raw.githubusercontent.com/sunshineplan/sda/main/public/style.css
    curl -LO https://raw.githubusercontent.com/sunshineplan/sda/main/public/index.html
    curl -LO https://raw.githubusercontent.com/sunshineplan/sda/main/scripts/sda.conf
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
    cp -s /var/www/sda/sda.conf /etc/nginx/conf.d
    sed -i "s/\$domain/$domain/" /var/www/sda/sda.conf
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
