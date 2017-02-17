
# Introduction

This is a fork from https://github.com/blablacar/phpfpm-prometheus-exporter to add PHP-FPM process status.

On the original version, if PHP-FPM process is not running the metrics are not updated and you can't detect it. 

In this fork, there is a new metric called "phpfpm_up" (1 PHP-FPM is running, 0 PHP-FPM is not running).

# Metrics

Metrics returned by the exporters:

	# HELP _process_cpu_seconds_total Total user and system CPU time spent in seconds.
	# TYPE _process_cpu_seconds_total counter
	_process_cpu_seconds_total 0
	# HELP _process_max_fds Maximum number of open file descriptors.
	# TYPE _process_max_fds gauge
	_process_max_fds 1024
	# HELP _process_open_fds Number of open file descriptors.
	# TYPE _process_open_fds gauge
	_process_open_fds 8
	# HELP _process_resident_memory_bytes Resident memory size in bytes.
	# TYPE _process_resident_memory_bytes gauge
	_process_resident_memory_bytes 7.335936e+06
	# HELP _process_start_time_seconds Start time of the process since unix epoch in seconds.
	# TYPE _process_start_time_seconds gauge
	_process_start_time_seconds 1.48734144512e+09
	# HELP _process_virtual_memory_bytes Virtual memory size in bytes.
	# TYPE _process_virtual_memory_bytes gauge
	_process_virtual_memory_bytes 1.26308352e+08
	# HELP go_gc_duration_seconds A summary of the GC invocation durations.
	# TYPE go_gc_duration_seconds summary
	go_gc_duration_seconds{quantile="0"} 0
	go_gc_duration_seconds{quantile="0.25"} 0
	go_gc_duration_seconds{quantile="0.5"} 0
	go_gc_duration_seconds{quantile="0.75"} 0
	go_gc_duration_seconds{quantile="1"} 0
	go_gc_duration_seconds_sum 0
	go_gc_duration_seconds_count 0
	# HELP go_goroutines Number of goroutines that currently exist.
	# TYPE go_goroutines gauge
	go_goroutines 13
	# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
	# TYPE go_memstats_alloc_bytes gauge
	go_memstats_alloc_bytes 1.130096e+06
	# HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
	# TYPE go_memstats_alloc_bytes_total counter
	go_memstats_alloc_bytes_total 1.130096e+06
	# HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
	# TYPE go_memstats_buck_hash_sys_bytes gauge
	go_memstats_buck_hash_sys_bytes 1.44262e+06
	# HELP go_memstats_frees_total Total number of frees.
	# TYPE go_memstats_frees_total counter
	go_memstats_frees_total 468
	# HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
	# TYPE go_memstats_gc_sys_bytes gauge
	go_memstats_gc_sys_bytes 131072
	# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
	# TYPE go_memstats_heap_alloc_bytes gauge
	go_memstats_heap_alloc_bytes 1.130096e+06
	# HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
	# TYPE go_memstats_heap_idle_bytes gauge
	go_memstats_heap_idle_bytes 172032
	# HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
	# TYPE go_memstats_heap_inuse_bytes gauge
	go_memstats_heap_inuse_bytes 1.564672e+06
	# HELP go_memstats_heap_objects Number of allocated objects.
	# TYPE go_memstats_heap_objects gauge
	go_memstats_heap_objects 8142
	# HELP go_memstats_heap_released_bytes_total Total number of heap bytes released to OS.
	# TYPE go_memstats_heap_released_bytes_total counter
	go_memstats_heap_released_bytes_total 0
	# HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
	# TYPE go_memstats_heap_sys_bytes gauge
	go_memstats_heap_sys_bytes 1.736704e+06
	# HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
	# TYPE go_memstats_last_gc_time_seconds gauge
	go_memstats_last_gc_time_seconds 0
	# HELP go_memstats_lookups_total Total number of pointer lookups.
	# TYPE go_memstats_lookups_total counter
	go_memstats_lookups_total 32
	# HELP go_memstats_mallocs_total Total number of mallocs.
	# TYPE go_memstats_mallocs_total counter
	go_memstats_mallocs_total 8610
	# HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
	# TYPE go_memstats_mcache_inuse_bytes gauge
	go_memstats_mcache_inuse_bytes 2400
	# HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
	# TYPE go_memstats_mcache_sys_bytes gauge
	go_memstats_mcache_sys_bytes 16384
	# HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
	# TYPE go_memstats_mspan_inuse_bytes gauge
	go_memstats_mspan_inuse_bytes 22240
	# HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
	# TYPE go_memstats_mspan_sys_bytes gauge
	go_memstats_mspan_sys_bytes 32768
	# HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
	# TYPE go_memstats_next_gc_bytes gauge
	go_memstats_next_gc_bytes 4.194304e+06
	# HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
	# TYPE go_memstats_other_sys_bytes gauge
	go_memstats_other_sys_bytes 804284
	# HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
	# TYPE go_memstats_stack_inuse_bytes gauge
	go_memstats_stack_inuse_bytes 360448
	# HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
	# TYPE go_memstats_stack_sys_bytes gauge
	go_memstats_stack_sys_bytes 360448
	# HELP go_memstats_sys_bytes Number of bytes obtained by system. Sum of all system allocations.
	# TYPE go_memstats_sys_bytes gauge
	go_memstats_sys_bytes 4.52428e+06
	# HELP phpfpm_accepted_conn The number of requests accepted by the pool
	# TYPE phpfpm_accepted_conn counter
	phpfpm_accepted_conn{pool_name="app1"} 2
	phpfpm_accepted_conn{pool_name="www"} 2
	# HELP phpfpm_active_processes The number of active processes
	# TYPE phpfpm_active_processes gauge
	phpfpm_active_processes{pool_name="app1"} 1
	phpfpm_active_processes{pool_name="www"} 1
	# HELP phpfpm_exporter_build_info A metric with a constant '1' value labeled by version, revision, branch, and goversion from which phpfpm_exporter was built.
	# TYPE phpfpm_exporter_build_info gauge
	phpfpm_exporter_build_info{branch="",goversion="go1.7.4",revision="",version=""} 1
	# HELP phpfpm_idle_processes The number of idle processes
	# TYPE phpfpm_idle_processes gauge
	phpfpm_idle_processes{pool_name="app1"} 1
	phpfpm_idle_processes{pool_name="www"} 1
	# HELP phpfpm_listen_queue The number of requests in the queue of pending connections
	# TYPE phpfpm_listen_queue gauge
	phpfpm_listen_queue{pool_name="app1"} 0
	phpfpm_listen_queue{pool_name="www"} 0
	# HELP phpfpm_listen_queue_len The size of the socket queue of pending connections
	# TYPE phpfpm_listen_queue_len gauge
	phpfpm_listen_queue_len{pool_name="app1"} 0
	phpfpm_listen_queue_len{pool_name="www"} 0
	# HELP phpfpm_max_active_processes The maximum number of active processes since FPM has started
	# TYPE phpfpm_max_active_processes counter
	phpfpm_max_active_processes{pool_name="app1"} 1
	phpfpm_max_active_processes{pool_name="www"} 1
	# HELP phpfpm_max_children_reached The number of times, the process limit has been reached, when pm tries to start more children (works only for pm 'dynamic' and 'ondemand')
	# TYPE phpfpm_max_children_reached counter
	phpfpm_max_children_reached{pool_name="app1"} 0
	phpfpm_max_children_reached{pool_name="www"} 0
	# HELP phpfpm_max_listen_queue The maximum number of requests in the queue of pending connections since FPM has started
	# TYPE phpfpm_max_listen_queue counter
	phpfpm_max_listen_queue{pool_name="app1"} 0
	phpfpm_max_listen_queue{pool_name="www"} 0
	# HELP phpfpm_process_cpu_seconds_total Total user and system CPU time spent in seconds.
	# TYPE phpfpm_process_cpu_seconds_total counter
	phpfpm_process_cpu_seconds_total 0.01
	# HELP phpfpm_process_max_fds Maximum number of open file descriptors.
	# TYPE phpfpm_process_max_fds gauge
	phpfpm_process_max_fds 1024
	# HELP phpfpm_process_resident_memory_bytes Resident memory size in bytes.
	# TYPE phpfpm_process_resident_memory_bytes gauge
	phpfpm_process_resident_memory_bytes 2.0668416e+07
	# HELP phpfpm_process_start_time_seconds Start time of the process since unix epoch in seconds.
	# TYPE phpfpm_process_start_time_seconds gauge
	phpfpm_process_start_time_seconds 1.48734142537e+09
	# HELP phpfpm_process_virtual_memory_bytes Virtual memory size in bytes.
	# TYPE phpfpm_process_virtual_memory_bytes gauge
	phpfpm_process_virtual_memory_bytes 1.78786304e+08
	# HELP phpfpm_slow_requests The number of requests that exceeded your request_slowlog_timeout value
	# TYPE phpfpm_slow_requests counter
	phpfpm_slow_requests{pool_name="app1"} 0
	phpfpm_slow_requests{pool_name="www"} 0
	# HELP phpfpm_start_since Number of seconds since FPM has started
	# TYPE phpfpm_start_since counter
	phpfpm_start_since{pool_name="app1"} 20
	phpfpm_start_since{pool_name="www"} 20
	# HELP phpfpm_total_processes The number of idle + active processes
	# TYPE phpfpm_total_processes gauge
	phpfpm_total_processes{pool_name="app1"} 2
	phpfpm_total_processes{pool_name="www"} 2
	# HELP phpfpm_up Whether the PHP-FPM process is up.
	# TYPE phpfpm_up counter
	phpfpm_up{pool_name="app1"} 1
	phpfpm_up{pool_name="www"} 1

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
