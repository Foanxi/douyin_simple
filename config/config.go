package config

type Config struct {
	Mysql    Mysql
	Ipconfig Ipconfig
}
type Mysql struct {
	Username string `mapstucture:"username"`
	Password string `mapstucture:"username"`
	Url      string `mapstucture:"url"`
}

type Ipconfig struct {
	Ip_url string `mapstucture:"ip_Url"`
}
