package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"github.com/tomasen/fcgi_client"
	"golang.org/x/net/context"
	"gopkg.in/ini.v1"
)

const (
	namespace = "phpfpm"
)

type FpmPoolMetrics struct {
	StartTime          int `json:"start time"`
	StartSince         int `json:"start since"`
	AcceptedConn       int `json:"accepted conn"`
	ListenQueue        int `json:"listen queue"`
	MaxListenQueue     int `json:"max listen queue"`
	ListenQueueLen     int `json:"listen queue len"`
	IdleProcesses      int `json:"idle processes"`
	ActiveProcesses    int `json:"active processes"`
	TotalProcesses     int `json:"total processes"`
	MaxActiveProcesses int `json:"max active processes"`
	MaxChildrenReached int `json:"max children reached"`
	SlowRequests       int `json:"slow requests"`
}

type PhpFpmPool struct {
	Name        string
	Endpoint    string
	StatusUri   string
	lastMetrics FpmPoolMetrics
	mu          sync.RWMutex
}

type PhpFpmPoolExporter struct {
	poolsToMonitor                                                                                 []*PhpFpmPool
	listenQueue, listenQueueLen, idleProcesses, activeProcesses, totalProcesses                    *prometheus.GaugeVec
	startSince, acceptedConn, maxListenQueue, maxActiveProcesses, maxChildrenReached, slowRequests *prometheus.CounterVec
}

func (e *PhpFpmPoolExporter) resetMetrics() {
	e.listenQueue.Reset()
	e.listenQueueLen.Reset()
	e.idleProcesses.Reset()
	e.activeProcesses.Reset()
	e.totalProcesses.Reset()
	e.startSince.Reset()
	e.acceptedConn.Reset()
	e.maxListenQueue.Reset()
	e.maxActiveProcesses.Reset()
	e.maxChildrenReached.Reset()
	e.slowRequests.Reset()
}

func (p *PhpFpmPool) GetSyncedCopy() PhpFpmPool {
	p.mu.Lock()
	pfp := p
	p.mu.Unlock()

	return *pfp
}

func (p *PhpFpmPool) PushSyncedLastMetrics(fpm *FpmPoolMetrics) {
	p.mu.Lock()
	p.lastMetrics = *fpm
	p.mu.Unlock()
}

func (p *PhpFpmPool) GetSyncedLastMetricsCopy() FpmPoolMetrics {
	p.mu.Lock()
	lm := &(p).lastMetrics
	p.mu.Unlock()

	return *lm
}

func (e *PhpFpmPoolExporter) Describe(ch chan<- *prometheus.Desc) {
	e.listenQueue.Describe(ch)
	e.listenQueueLen.Describe(ch)
	e.idleProcesses.Describe(ch)
	e.activeProcesses.Describe(ch)
	e.totalProcesses.Describe(ch)
	e.startSince.Describe(ch)
	e.acceptedConn.Describe(ch)
	e.maxListenQueue.Describe(ch)
	e.maxActiveProcesses.Describe(ch)
	e.maxChildrenReached.Describe(ch)
	e.slowRequests.Describe(ch)
}

func (e *PhpFpmPoolExporter) Collect(ch chan<- prometheus.Metric) {
	e.resetMetrics()
	for _, p := range e.poolsToMonitor {
		lastMetric := p.GetSyncedLastMetricsCopy()

		(e.listenQueue.WithLabelValues(p.Name)).Set(float64(lastMetric.ListenQueue))
		(e.listenQueueLen.WithLabelValues(p.Name)).Set(float64(lastMetric.ListenQueueLen))
		(e.idleProcesses.WithLabelValues(p.Name)).Set(float64(lastMetric.IdleProcesses))
		(e.activeProcesses.WithLabelValues(p.Name)).Set(float64(lastMetric.ActiveProcesses))
		(e.totalProcesses.WithLabelValues(p.Name)).Set(float64(lastMetric.TotalProcesses))
		(e.startSince.WithLabelValues(p.Name)).Add(float64(lastMetric.StartSince))
		(e.acceptedConn.WithLabelValues(p.Name)).Add(float64(lastMetric.AcceptedConn))
		(e.maxListenQueue.WithLabelValues(p.Name)).Add(float64(lastMetric.MaxListenQueue))
		(e.maxActiveProcesses.WithLabelValues(p.Name)).Add(float64(lastMetric.MaxActiveProcesses))
		(e.maxChildrenReached.WithLabelValues(p.Name)).Add(float64(lastMetric.MaxChildrenReached))
		(e.slowRequests.WithLabelValues(p.Name)).Add(float64(lastMetric.SlowRequests))

		(e.listenQueue.WithLabelValues(p.Name)).Collect(ch)
		(e.listenQueueLen.WithLabelValues(p.Name)).Collect(ch)
		(e.idleProcesses.WithLabelValues(p.Name)).Collect(ch)
		(e.activeProcesses.WithLabelValues(p.Name)).Collect(ch)
		(e.totalProcesses.WithLabelValues(p.Name)).Collect(ch)
		(e.startSince.WithLabelValues(p.Name)).Collect(ch)
		(e.acceptedConn.WithLabelValues(p.Name)).Collect(ch)
		(e.maxListenQueue.WithLabelValues(p.Name)).Collect(ch)
		(e.maxActiveProcesses.WithLabelValues(p.Name)).Collect(ch)
		(e.maxChildrenReached.WithLabelValues(p.Name)).Collect(ch)
		(e.slowRequests.WithLabelValues(p.Name)).Collect(ch)

		log.Debugln("Metrics collection completed!")
	}
}

func NewPhpFpmPoolExporter(pools []*PhpFpmPool) *PhpFpmPoolExporter {
	poolLabelNames := []string{"pool_name"}

	return &PhpFpmPoolExporter{
		poolsToMonitor: pools,
		startSince: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "start_since",
				Help:      "Number of seconds since FPM has started",
			},
			poolLabelNames,
		),
		acceptedConn: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "accepted_conn",
				Help:      "The number of requests accepted by the pool",
			},
			poolLabelNames,
		),
		listenQueue: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "listen_queue",
				Help:      "The number of requests in the queue of pending connections",
			},
			poolLabelNames,
		),
		maxListenQueue: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "max_listen_queue",
				Help:      "The maximum number of requests in the queue of pending connections since FPM has started",
			},
			poolLabelNames,
		),
		listenQueueLen: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "listen_queue_len",
				Help:      "The size of the socket queue of pending connections",
			},
			poolLabelNames,
		),
		idleProcesses: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "idle_processes",
				Help:      "The number of idle processes",
			},
			poolLabelNames,
		),
		activeProcesses: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "active_processes",
				Help:      "The number of active processes",
			},
			poolLabelNames,
		),
		totalProcesses: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "total_processes",
				Help:      "The number of idle + active processes",
			},
			poolLabelNames,
		),
		maxActiveProcesses: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "max_active_processes",
				Help:      "The maximum number of active processes since FPM has started",
			},
			poolLabelNames,
		),
		maxChildrenReached: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "max_children_reached",
				Help:      "The number of times, the process limit has been reached, when pm tries to start more children (works only for pm 'dynamic' and 'ondemand')",
			},
			poolLabelNames,
		),
		slowRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "slow_requests",
				Help:      "The number of requests that exceeded your request_slowlog_timeout value",
			},
			poolLabelNames,
		),
	}
}

func GetFilesIn(dirPath string) []string {
	var poolFiles []string

	if strings.HasSuffix(dirPath, "/") {
		dirPath = strings.TrimRight(dirPath, "/")
	}

	dir, err := os.Open(dirPath)

	if err != nil {
		fmt.Errorf("%s", err)
		return nil
	}
	defer dir.Close()

	filesInfo, err := dir.Readdir(-1)

	if err != nil {
		fmt.Errorf("%s", err)
		return nil
	}

	for i := 0; i < len(filesInfo); i++ {
		if filesInfo[i].Mode().IsRegular() {
			poolFiles = append(poolFiles, dirPath+"/"+filesInfo[i].Name())
		}
	}

	return poolFiles
}

func PollFpmStatusMetrics(p *PhpFpmPool, fetcher func() (string, error), pollInterval int, mustQuit chan bool, done chan bool) {

	var mts FpmPoolMetrics
	var res string
	var err error

	for i := 0; i < 1; {
		res, err = fetcher()

		log.Debugln(p.Name + " - End of fetch logic")

		if err != nil {
			log.Errorln(err.Error())
		} else {
			err = json.Unmarshal([]byte(res), &mts)

			if err != nil {
				log.Errorln(err.Error())
			} else {

				log.Debugln(p.Name + " - StartTime read on status: " + strconv.Itoa(mts.StartTime))
				log.Debugln(p.Name + " - StartSince read on status: " + strconv.Itoa(mts.StartSince))
				log.Debugln(p.Name + " - AcceptedConn read on status: " + strconv.Itoa(mts.AcceptedConn))
				log.Debugln(p.Name + " - ListenQueue read on status: " + strconv.Itoa(mts.ListenQueue))
				log.Debugln(p.Name + " - MaxListenQueue read on status: " + strconv.Itoa(mts.MaxListenQueue))
				log.Debugln(p.Name + " - ListenQueueLen read on status: " + strconv.Itoa(mts.ListenQueueLen))
				log.Debugln(p.Name + " - IdleProcesses read on status: " + strconv.Itoa(mts.IdleProcesses))
				log.Debugln(p.Name + " - ActiveProcesses read on status: " + strconv.Itoa(mts.ActiveProcesses))
				log.Debugln(p.Name + " - TotalProcesses read on status: " + strconv.Itoa(mts.TotalProcesses))
				log.Debugln(p.Name + " - MaxActiveProcesses read on status: " + strconv.Itoa(mts.MaxActiveProcesses))
				log.Debugln(p.Name + " - MaxChildrenReached read on status: " + strconv.Itoa(mts.MaxChildrenReached))
				log.Debugln(p.Name + " - SlowRequests read on status: " + strconv.Itoa(mts.SlowRequests))

				p.PushSyncedLastMetrics(&mts)
				log.Debugln(p.Name + " - Metrics pushed to pool structure")
			}
		}

		time.Sleep(time.Duration(pollInterval * int(time.Second)))
		select {
		case <-mustQuit:
			i = 1
			log.Infoln("Goroutine received signal asking to quit")
			done <- true
		default:
			continue
		}
	}
	return
}

func NativeClientFcgiStatusFetcher(p *PhpFpmPool, fcgiConnectTimeout int) func() (string, error) {
	poolCpy := p.GetSyncedCopy()
	endpoint := poolCpy.Endpoint

	env := make(map[string]string)
	env["SCRIPT_NAME"] = poolCpy.StatusUri
	env["SCRIPT_FILENAME"] = poolCpy.StatusUri
	env["QUERY_STRING"] = "json"
	env["SERVER_SOFTWARE"] = "go/fcgiclient"

	return func() (string, error) {

		var netType string
		fileInfo, err := os.Stat(endpoint)
		if err != nil {
			netType = "tcp"
		} else {
			if fileInfo.Mode()&os.ModeSocket != 0 {
				netType = "unix"
			} else {
				netType = "tcp"
			}
		}

		log.Debugln(endpoint + " has been identified as " + netType + " network type")

		fcgi, err := fcgiclient.DialTimeout(netType, endpoint, time.Duration(fcgiConnectTimeout*int(time.Millisecond)))
		if err != nil {
			return "", err
		}

		defer fcgi.Close()

		resp, err := fcgi.Get(env)
		if err != nil {
			//fcgi.Close()
			return "", err
		}

		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			//fcgi.Close()
			return "", err
		}

		//fcgi.Close()
		return string(content), nil
	}
}

func CgiFcgiFcgiStatusFetcher(p *PhpFpmPool, pollTimeout int, cgiFastCgiPath string, cgiFastCgiLdLibPath string) func() (string, error) {
	poolCpy := p.GetSyncedCopy()
	endpoint := poolCpy.Endpoint

	env := os.Environ()

	if cgiFastCgiLdLibPath != "" {
		env = append(env, fmt.Sprintf("LD_LIBRARY_PATH=%s", cgiFastCgiLdLibPath))
	}

	env = append(env, fmt.Sprintf("SCRIPT_NAME=%s", poolCpy.StatusUri))
	env = append(env, fmt.Sprintf("SCRIPT_FILENAME=%s", poolCpy.StatusUri))
	env = append(env, "QUERY_STRING=json")
	env = append(env, "REQUEST_METHOD=GET")

	return func() (string, error) {
		var data []byte
		var err error
		var strData []string

		ctx := context.TODO()
		ctxWithCancel, cancel := context.WithTimeout(ctx, time.Duration(pollTimeout*int(time.Second)))
		defer cancel()

		cmd := exec.CommandContext(ctxWithCancel, cgiFastCgiPath, "-bind", "-connect", endpoint)
		cmd.Env = env
		data, err = cmd.Output()
		if err != nil {
			return "", err
		}

		strData = strings.SplitAfter(string(data), "\r\n\r\n")
		if len(strData) < 2 {
			return "", errors.New("Unexpected cgi-fcgi response")
		}

		return strData[1], nil
	}
}

func main() {
	var (
		listenAddress       = flag.String("web.listen-address", ":9101", "Address to listen on for web interface and telemetry.")
		metricsPath         = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
		phpfpmPidFile       = flag.String("phpfpm.pid-file", "/var/run/php5-fpm.pid", "Path to phpfpm's pid file.")
		configDir           = flag.String("phpfpm.config", "/etc/php5/fpm/pool.d/", "Pools conf dir")
		pollInterval        = flag.Int("phpfpm.poll-interval", 10, "Poll interval in seconds")
		useNativeClient     = flag.Bool("phpfpm.use-native-client", true, "Use a native go client to get status or use external cgi-fcgi command")
		ncConnectTimeout    = flag.Int("nc.connect-timeout", 2, "Native client connect timeout")
		pollTimeout         = flag.Int("cgi-fcgi.poll-timeout", 2, "Poll timeout in seconds")
		cgiFastCgiPath      = flag.String("cgi-fcgi.path", "/usr/bin/cgi-fcgi", "cgi-fcgi program path")
		cgiFastCgiLdLibPath = flag.String("cgi-fcgi.ld-library-path", "", "LD_LIBRARY_PATH value to run cgi-fcgi")
		showVersion         = flag.Bool("version", false, "Print version information.")
	)

	flag.Parse()

	if *showVersion {
		fmt.Fprintln(os.Stdout, version.Print("phpfpm_prometheus_exporter"))
		os.Exit(0)
	}

	log.Infoln("Starting phpfpm_exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	if *phpfpmPidFile != "" {
		log.Debugln("Export master process metrics enabled")

		procExporter := prometheus.NewProcessCollectorPIDFn(
			func() (int, error) {
				content, err := ioutil.ReadFile(*phpfpmPidFile)
				if err != nil {
					return 0, fmt.Errorf("Can't read pid file: %s", err)
				}
				value, err := strconv.Atoi(strings.TrimSpace(string(content)))
				if err != nil {
					return 0, fmt.Errorf("Can't parse pid file: %s", err)
				}
				return value, nil
			}, namespace)
		prometheus.MustRegister(procExporter)
	} else {
		log.Debugln("Export master process metrics disabled")
	}

	sigs := make(chan os.Signal)
	mustQuit := make(chan bool)
	done := make(chan bool)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	phpFpmPools := []*PhpFpmPool{}

	confFiles := GetFilesIn(*configDir)
	cfg := ini.Empty()

	for _, cf := range confFiles {
		log.Infoln("We will parse: ", cf)
		err := cfg.Append(cf)

		if err != nil {
			fmt.Errorf("%f")
		}
	}

	sections := cfg.SectionStrings()
	sectionsCount := 0

	for _, sect := range sections {
		statusKey, err := cfg.Section(sect).GetKey("pm.status_path")

		if err != nil {
			continue
		}

		listenKey, err := cfg.Section(sect).GetKey("listen")

		if err != nil {
			continue
		}

		pool := PhpFpmPool{Name: sect, Endpoint: listenKey.String(), StatusUri: statusKey.String()}

		mts := FpmPoolMetrics{
			StartTime:          0,
			StartSince:         0,
			AcceptedConn:       0,
			ListenQueue:        0,
			MaxListenQueue:     0,
			ListenQueueLen:     0,
			IdleProcesses:      0,
			ActiveProcesses:    0,
			TotalProcesses:     0,
			MaxActiveProcesses: 0,
			MaxChildrenReached: 0,
			SlowRequests:       0,
		}

		pool.PushSyncedLastMetrics(&mts)

		sectionsCount++

		var fetcher func() (string, error)

		if *useNativeClient {
			fetcher = NativeClientFcgiStatusFetcher(&pool, *ncConnectTimeout)
		} else {
			fetcher = CgiFcgiFcgiStatusFetcher(&pool, *pollTimeout, *cgiFastCgiPath, *cgiFastCgiLdLibPath)
		}

		go PollFpmStatusMetrics(&pool, fetcher, *pollInterval, mustQuit, done)

		phpFpmPools = append(phpFpmPools, &pool)
	}

	log.Infoln("We will monitor ", sectionsCount, " phpfpm pool(s)")

	phpFpmExporter := NewPhpFpmPoolExporter(phpFpmPools)

	prometheus.MustRegister(phpFpmExporter)
	prometheus.MustRegister(version.NewCollector("phpfpm_exporter"))

	log.Infoln("Listening on", *listenAddress)
	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
      <head><title>PhpFpm Exporter</title></head>
      <body>
      <h1>PhpFpm Exporter</h1>
      <p><a href='` + *metricsPath + `'>Metrics</a></p>
      </body>
      </html>`))
	})
	//log.Fatal(http.ListenAndServe(*listenAddress, nil))
	go http.ListenAndServe(*listenAddress, nil)

	log.Infoln("Awaiting quit signal")

	<-sigs

	for j := 0; j < sectionsCount; j++ {
		mustQuit <- true
	}

	log.Infoln("Awaiting all done signals")

	for j := 0; j < sectionsCount; j++ {
		<-done
	}
	close(mustQuit)
	close(done)
	log.Infoln("Clean shutdown!")
}
