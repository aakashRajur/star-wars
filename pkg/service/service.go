package service

type Service struct {
	Id          string
	Name        string
	Scheme      string
	Hostname    string
	Port        int
	Healthcheck Healthcheck
}
