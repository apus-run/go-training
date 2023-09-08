package conf

type Conf struct {
	Server struct {
		Addr    string `json:"addr"`
		Timeout string `json:"timeout"`
	} `json:"server"`

	Data struct {
		Database struct {
			Driver string `json:"driver,omitempty"`
			Dsn    string `json:"dsn,omitempty"`
		} `json:"database,omitempty"`

		Redis struct {
			Addr         string `json:"addr,omitempty"`
			Password     string `json:"password,omitempty"`
			Db           int    `json:"db,omitempty"`
			ReadTimeout  string `json:"read_timeout,omitempty"`
			WriteTimeout string `json:"write_timeout,omitempty"`
		} `json:"redis,omitempty"`

		Memory struct {
			Size int `json:"size,omitempty"`
		} `json:"memory,omitempty"`
	} `json:"data,omitempty"`

	Jwt struct {
		Secret string `json:"secret,omitempty"`
		Expire int    `json:"expire,omitempty"`
	} `json:"jwt,omitempty"`
}
