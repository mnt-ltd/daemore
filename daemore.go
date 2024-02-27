package daemore

import (
	"os"

	"github.com/kardianos/service"
	"go.uber.org/zap"
)

type Daemon struct {
	config   *service.Config
	errs     chan error
	logger   *zap.Logger
	callback func()
	Service  service.Service
}

type DaemonOption struct {
	Name        string
	DisplayName string
	Description string
	Logger      *zap.Logger
	RunCallback func()
}

func NewDaemon(opt DaemonOption, username ...string) (d *Daemon, err error) {
	config := &service.Config{
		Name:        opt.Name,        //服务显示名称
		DisplayName: opt.DisplayName, //服务名称
		Description: opt.Description, //服务描述
	}
	if len(username) > 0 {
		config.UserName = username[0]
	}

	if opt.Logger == nil {
		opt.Logger, _ = zap.NewDevelopment()
	}

	d = &Daemon{
		config:   config,
		errs:     make(chan error, 100),
		logger:   opt.Logger.Named("daemore"),
		callback: opt.RunCallback,
	}

	d.Service, err = service.New(d, d.config)
	if err != nil {
		d.logger.Fatal("new", zap.Error(err))
	}
	return
}

func (d *Daemon) Config() *service.Config {
	return d.config
}

func (d *Daemon) Start(s service.Service) error {
	go d.Run()
	return nil
}

func (d *Daemon) Run() {
	if d.callback == nil {
		d.logger.Fatal("run", zap.String("callback", "callback is nil"))
	}
	d.logger.Info("run...")
	d.callback()
}

func (d *Daemon) Stop(s service.Service) error {
	if service.Interactive() {
		os.Exit(0)
	}
	return nil
}

// RunAsService 运行为服务
func (d *Daemon) Install(args ...string) {
	d.logger.Debug("install", zap.Any("args", args))
	d.config.Arguments = args
	s, err := service.New(d, d.config)
	if err != nil {
		d.logger.Fatal("install", zap.Error(err))
	}
	err = s.Install()
	if err != nil {
		d.logger.Fatal("install", zap.Error(err))
	}
	d.logger.Info("install success")
	os.Exit(0)
}

// RunAsService 运行为服务
func (d *Daemon) Uninstall() {
	s, err := service.New(d, d.config)
	if err != nil {
		d.logger.Fatal("uninstall", zap.Error(err))
	}
	err = s.Uninstall()
	if err != nil {
		d.logger.Fatal("uninstall", zap.Error(err))
	}
	d.logger.Info("uninstall success")
	os.Exit(0)
}

// RunAsService 运行为服务
func (d *Daemon) Restart() {
	s, err := service.New(d, d.config)
	if err != nil {
		d.logger.Fatal("restart", zap.Error(err))
	}
	err = s.Restart()
	if err != nil {
		d.logger.Fatal("restart", zap.Error(err))
	}
	d.logger.Info("restart success")
	os.Exit(0)
}
