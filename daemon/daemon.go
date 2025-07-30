// Package daemon provides utilities for creating and managing system daemons/services.
// It leverages the 'github.com/kardianos/service' package to provide cross-platform daemon capabilities.
package daemon

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/kardianos/service"
)

// Service defines the interface that any service managed by the Daemon must implement.
type Service interface {
	// Name returns the unique name of the service.
	Name() string
	// Run contains the main logic of the service. It should block until the service is stopped
	// or an unrecoverable error occurs. It should return an error if the service fails to start or encounters a critical issue.
	Run() error
	// HandleError is called when an error occurs within the service's Run method.
	// Implementations can use this to log the error, attempt recovery, or signal a shutdown.
	HandleError(error)
}

// Daemon is the core structure that manages multiple services.
// It implements the 'service.Service' interface from the kardianos/service package.
type Daemon struct {
	config   *service.Config // config holds the service configuration.
	services []Service       // services is a slice of individual services to be managed.
	errs     chan error      // errs is a channel to receive errors from running services.
}

// NewDaemon creates and returns a new Daemon instance.
// It initializes the daemon with a service configuration and a list of services to run.
func NewDaemon(conf *service.Config, services ...Service) *Daemon {
	return &Daemon{
		config:   conf,
		errs:     make(chan error, 100),
		services: services,
	}
}

// Start is part of the 'service.Service' interface.
// It is called by the service manager when the daemon is started.
// This method launches a goroutine for each managed service and a goroutine to listen for errors.
func (d *Daemon) Start(_ service.Service) error {
	// Start a goroutine to listen for errors from individual services.
	go d.handleServiceErrors()

	// Start each service in its own goroutine.
	for _, svc := range d.services {
		go func(s Service) {
			slog.Info(fmt.Sprintf("Starting service: %s", s.Name()))
			if err := s.Run(); err != nil {
				s.HandleError(err)
				d.errs <- fmt.Errorf("service %s failed: %w", s.Name(), err)
			}
			slog.Info(fmt.Sprintf("Service %s stopped.", s.Name()))
		}(svc)
	}
	return nil
}

// handleServiceErrors listens for errors sent by individual services on the errs channel.
// It logs these errors. In a more complex scenario, it might trigger service restarts or other recovery actions.
func (d *Daemon) handleServiceErrors() {
	for err := range d.errs {
		slog.Error("Daemon service error", "error", err)
		// Here, you could add logic to restart a failed service,
		// or decide to stop the entire daemon if a critical service fails.
	}
}

// Stop is part of the 'service.Service' interface.
// It is called by the service manager when the daemon is stopped.
// This method should perform any necessary cleanup before the daemon exits.
func (d *Daemon) Stop(_ service.Service) error {
	// Close the error channel to signal handleServiceErrors to stop.
	close(d.errs)

	// In a real-world scenario, you might want to gracefully stop each running service here.
	// For example, by sending a signal to their goroutines or calling a Stop() method on them.
	// The current implementation relies on the service manager to terminate the process.

	// If running interactively (e.g., from command line), exit the process.
	// This is typically handled by the service package itself, but added for clarity.
	if service.Interactive() {
		slog.Info("Daemon running interactively, exiting.")
		os.Exit(0) // Consider removing this for more graceful shutdown in production daemons.
	}

	return nil
}

type Controller struct {
	service service.Service // service is the underlying service.Service instance.
}

// NewController creates and returns a new Controller instance.
// It sets up the daemon and registers it with the system's service manager.
func NewController(conf *service.Config, services ...Service) (*Controller, error) {
	d := NewDaemon(conf, services...)

	// Apply platform-specific options.
	// For systemd on Linux, increase the file descriptor limit.
	if service.Platform() == "linux-systemd" {
		d.config.Option = service.KeyValue{
			"LimitNOFILE": 40960, // Set a higher limit for open file descriptors.
		}
	}

	// Create a new service instance with the daemon and its configuration.
	s, err := service.New(d, d.config)
	if err != nil {
		return nil, fmt.Errorf("failed to create new service: %w", err)
	}

	return &Controller{
		service: s,
	}, nil
}

// Install installs the daemon as a system service.
func (d *Controller) Install() error {
	slog.Info(fmt.Sprintf("Installing service: %s", d.service.String()))
	if err := d.service.Install(); err != nil {
		return fmt.Errorf("failed to install service: %w", err)
	}
	slog.Info(fmt.Sprintf("Service %s installed successfully.", d.service.String()))
	return nil
}

// Uninstall uninstalls the daemon from the system services.
// It attempts to stop the service before uninstalling.
func (d *Controller) Uninstall() error {
	slog.Info(fmt.Sprintf("Uninstalling service: %s", d.service.String()))
	// Attempt to stop the service first to ensure a clean uninstall.
	if err := d.service.Stop(); err != nil {
		slog.Warn(fmt.Sprintf("Failed to stop service %s during uninstall: %v", d.service.String(), err))
		// Do not return error here, try to uninstall anyway.
	}
	// Proceed with uninstalling the service.
	if err := d.service.Uninstall(); err != nil {
		return fmt.Errorf("failed to uninstall service: %w", err)
	}
	slog.Info(fmt.Sprintf("Service %s uninstalled successfully.", d.service.String()))
	return nil
}

// Start starts the installed daemon service.
func (d *Controller) Start() error {
	slog.Info(fmt.Sprintf("Starting service: %s", d.service.String()))
	if err := d.service.Start(); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}
	slog.Info(fmt.Sprintf("Service %s started successfully.", d.service.String()))
	return nil
}

// Stop stops the running daemon service.
func (d *Controller) Stop() error {
	slog.Info(fmt.Sprintf("Stopping service: %s", d.service.String()))
	if err := d.service.Stop(); err != nil {
		return fmt.Errorf("failed to stop service: %w", err)
	}
	slog.Info(fmt.Sprintf("Service %s stopped successfully.", d.service.String()))
	return nil
}

// Restart restarts the installed daemon service.
func (d *Controller) Restart() error {
	slog.Info(fmt.Sprintf("Restarting service: %s", d.service.String()))
	if err := d.service.Restart(); err != nil {
		return fmt.Errorf("failed to restart service: %w", err)
	}
	slog.Info(fmt.Sprintf("Service %s restarted successfully.", d.service.String()))
	return nil
}

// Run runs the daemon service in the foreground.
// This is typically used for debugging or when running the service directly without installation.
func (d *Controller) Run() error {
	slog.Info(fmt.Sprintf("Running service: %s in foreground.", d.service.String()))
	if err := d.service.Run(); err != nil {
		return fmt.Errorf("failed to run service in foreground: %w", err)
	}
	slog.Info(fmt.Sprintf("Service %s finished running.", d.service.String()))
	return nil
}
