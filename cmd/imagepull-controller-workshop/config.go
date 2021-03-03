package main

import (
	"os"
	"path/filepath"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
	"k8s.io/client-go/util/homedir"
)

// CmdConfig represents the configuration of the command.
type CmdConfig struct {
	Development      bool
	Debug            bool
	Workers          int
	KubeConfig       string
	NamespaceRunning string
	ResyncInterval   time.Duration
	SecretName       string
}

// NewCmdConfig returns a new command configuration.
func NewCmdConfig() (*CmdConfig, error) {
	kubeHome := filepath.Join(homedir.HomeDir(), ".kube", "config")

	c := &CmdConfig{}
	app := kingpin.New("imagepullsecret-controller", "A Kubernetes controller to spread imagepullsecrets on namespaces.")

	app.Flag("debug", "Enable debug mode.").BoolVar(&c.Debug)
	app.Flag("development", "Enable development mode.").BoolVar(&c.Development)
	app.Flag("kube-config", "kubernetes configuration path, only used when development mode enabled.").Default(kubeHome).Short('c').StringVar(&c.KubeConfig)
	app.Flag("namespace-running", "kubernetes namespace where the controller is running.").Short('r').Required().StringVar(&c.NamespaceRunning)
	app.Flag("workers", "concurrent processing workers for each kubernetes controller.").Default("5").Short('w').IntVar(&c.Workers)
	app.Flag("resync-interval", "the duration between resync the controllers resources.").Default("5m").DurationVar(&c.ResyncInterval)
	app.Flag("secret-name", "the secret name in the running ns that has the image pull credentials.").Default("image-pull-credentials").StringVar(&c.SecretName)

	_, err := app.Parse(os.Args[1:])
	if err != nil {
		return nil, err
	}

	return c, nil
}
