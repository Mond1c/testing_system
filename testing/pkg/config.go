package pkg

type WorkerConfig struct {
	Type    string `json:"type"`
	Address string `json:"address"`
}

type Config struct {
	Workers []WorkerConfig `json:"workers"`
}
