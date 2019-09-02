package config

type App struct {
	Name             string `json:"name,omitempty"`
	Version          string `json:"version,omitempty"`
	RegisterTTL      int    `json:"register_ttl,omitempty"`
	RegisterInterval int    `json:"register_interval,omitempty"`
	MaxConcurrent    int    `json:"max_concurrent,omitempty"`
	RpsLimit         int    `json:"rps_limit,omitempty"`
	TraceAddr        string `json:"trace_addr,omitempty"`
	BrokerAddr       string `json:"broker_addr,omitempty"`
}

type Mysql struct {
	DataSource string `json:"data_source,omitempty"`
	MaxIdle    int    `json:"max_idle,omitempty"`
	MaxOpen    int    `json:"max_open,omitempty"`
}

type Redis struct {
	Addr     string `json:"addr,omitempty"`
	Password string `json:"password,omitempty"`
	MaxIdle  int    `json:"max_idle,omitempty"`
}

type Hystrix struct {
	Timeout                int `json:"timout,omitempty"`
	MaxConcurrentRequests  int `json:"max_concurrent_requests,omitempty"`
	RequestVolumeThreshold int `json:"request_volume_threshold,omitempty"`
	SleepWindow            int `json:"sleep_window,omitempty"`
	ErrorPercentThreshold  int `json:"error_percent_threshold,omitempty"`
}

type Consul struct {
	Addrs []string `json:"addrs,omitempty"`
}

type Logger struct {
	Level      string `json:"level,omitempty"`
	Filename   string `json:"filename,omitempty"`
	MaxSize    int    `json:"max_size,omitempty"`
	MaxBackups int    `json:"max_backups,omitempty"`
	MaxAge     int    `json:"max_age,omitempty"`
	Compress   bool   `json:"compress,omitempty"`
}

type Broker struct {
	Addrs       []string `json:"addrs,omitempty"`
	ClusterID   string   `json:"cluster_id,omitempty"`
	DurableName string   `json:"durable_name,omitempty"`
	Queue       string   `json:"queue,omitempty"`
}

