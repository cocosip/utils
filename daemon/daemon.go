package daemon

import (
	"github.com/kardianos/service"
	"os"
)

type Service interface {
	Name() string
	Run() error
	HandleError(error)
}

type Daemon struct {
	config   *service.Config
	services []Service
	errs     chan error
}

func NewDaemon(conf *service.Config, services ...Service) *Daemon {
	return &Daemon{
		config:   conf,
		errs:     make(chan error, 100),
		services: services,
	}
}

func (d *Daemon) Start(_ service.Service) error {
	go d.run()
	return nil
}

func (d *Daemon) run() {
	for _, svc := range d.services {
		if err := svc.Run(); err != nil {
			svc.HandleError(err)
			return
		}
	}
}

func (d *Daemon) Stop(_ service.Service) error {
	if service.Interactive() {
		os.Exit(0)
	}
	return nil
}

type Controller struct {
	service service.Service
}

func NewController(conf *service.Config, services ...Service) (*Controller, error) {
	d := NewDaemon(conf, services...)
	if service.Platform() == "linux-systemd" {
		d.config.Option = service.KeyValue{
			"LimitNOFILE": 40960,
		}
	}
	s, err := service.New(d, d.config)
	if err != nil {
		return nil, err
	}

	return &Controller{
		service: s,
	}, nil
}

func (d *Controller) Install() error {
	return d.service.Install()
}

func (d *Controller) Uninstall() error {
	if err := d.service.Stop(); err != nil {
		return err
	}
	return d.service.Uninstall()
}

func (d *Controller) Start() error {
	return d.service.Start()
}

func (d *Controller) Stop() error {
	return d.service.Stop()
}

func (d *Controller) Restart() error {
	return d.service.Restart()
}

func (d *Controller) Run() error {
	return d.service.Run()
}
