
# Introduction

This is a fork from https://github.com/blablacar/phpfpm-prometheus-exporter to add PHP-FPM process status.

On the original version, if PHP-FPM process is not running the metrics are not updated and you can't detect it. 

In this fork, there is a new metric called "phpfpm_up" (1 PHP-FPM is running, 0 PHP-FPM is not running).

# Compilation

## Requirements

- GO >= 1.5
- GIT

## Debian Jessie

For Debian Jessie (current Debian stable version), you need to use Debian backports to use GO > 1.3.

Add debian-backports repository:

	cat > /etc/apt/sources.list.d/debian-backports.list << EOF
	deb http://http.debian.net/debian jessie-backports main contrib non-free
	EOF

Update Debian packages:

	apt-get update

Install GO 1.7:

	apt-get install -t jessie-backports golang

### Compilation process

To compile sources:

	export GOPATH=~/go
	mkdir $GOPATH
	cd $GOPATH
	mkdir -p src/github.com/vmercierfr
	git clone https://github.com/vmercierfr/phpfpm-prometheus-exporter.git src/github.com/vmercierfr/phpfpm-prometheus-exporter
	cd src/github.com/vmercierfr/phpfpm-prometheus-exporter
	go get
	go build

# Configuration

## PHP-FPM

You must enable status page on PHP-FPM pool.

For example, add the following line in /etc/php5/fpm/pool.d/www.conf to enable status page on **/status**:

    pm.status_path = /status

## Webserver

You must configure your webserver to forward status page URL to PHP-FPM process.

For example, Nginx configuration to forward **/status** requests to PHP-FPM process:

	location ~ ^/(status|ping)$ {
	     allow 127.0.0.1;
	     deny all;
	     include fastcgi_params;
	     fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
	     fastcgi_pass unix:/var/run/php5-fpm.sock;
	}

## Prometheus rule

If you use Prometheus Alert Manager, you can create a rule to raise an alert when PHP-FPM process is not running.

You can use the following rule:

	ALERT ProcessPhpfpmDown
	  IF phpfpm_up == 0
	  FOR 1m
	  LABELS { severity = "php" }
	  ANNOTATIONS {
	    summary = "{{ $labels.instance }} - PHP-FPM is down",
	    description = "Process PHP-FPM on {{ $labels.instance }} is down!",
	  }

# Running

If you want manage PHP-FPM exporter process with Systemd, you can use the following Systemd commands:

Copy builded exporter:

	mkdir /opt/prometheus_exporters
	cp $GOPATH/src/github.com/vmercierfr/phpfpm-prometheus-exporter/phpfpm-prometheus-exporter /opt/prometheus_exporters/phpfpm-prometheus-exporter-$(git rev-parse HEAD)
	ln -s /opt/prometheus_exporters/phpfpm-prometheus-exporter-$(git rev-parse HEAD) /opt/prometheus_exporters/phpfpm-prometheus-exporter

Create Systemd file:

	vim /etc/systemd/system/promotheus_phpfpm_exporter.service

Add the following configuration:

	[Unit]
	Description=Promotd heus PHP-FPM Exporter
	
	[Service]
	User=www-data
	ExecStart=/opt/prometheus_exporters/phpfpm-prometheus-exporter 
	
	[Install]
	WantedBy=default.target
	
Start and enable service:

	systemctl daemon-reload
	systemctl start promotheus_phpfpm_exporter.service
	systemctl enable promotheus_phpfpm_exporter.service

Check the process is running:

	systemctl status promotheus_phpfpm_exporter.service
