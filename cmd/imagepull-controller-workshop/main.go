package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/oklog/run"
	"github.com/sirupsen/logrus"
	koopercontroller "github.com/spotahome/kooper/v2/controller"
	kooperlog "github.com/spotahome/kooper/v2/log/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	controllernamespace "github.com/slok/imagepull-controller-workshop/internal/controller/namespace"
	controllersecretcache "github.com/slok/imagepull-controller-workshop/internal/controller/secretcache"
	loglogrus "github.com/slok/imagepull-controller-workshop/internal/log/logrus"
	storagekubernetes "github.com/slok/imagepull-controller-workshop/internal/storage/kubernetes"
)

// Run runs the main application.
func Run(ctx context.Context, stdin io.Writer, stdout, stderr io.Reader) error {
	// Load command flags and arguments.
	cmdCfg, err := NewCmdConfig()
	if err != nil {
		return fmt.Errorf("could not load command configuration: %w", err)
	}

	// Set up logger.
	logrusLog := logrus.New()
	logrusLogEntry := logrus.NewEntry(logrusLog)
	kooperLogger := kooperlog.New(logrusLogEntry.WithField("lib", "kooper"))
	logger := loglogrus.NewLogrus(logrusLogEntry)
	if cmdCfg.Debug {
		logrusLog.SetLevel(logrus.DebugLevel)
	}

	// Load Kubernetes clients.
	logger.Infof("loading Kubernetes configuration...")
	kcfg, err := loadKubernetesConfig(*cmdCfg)
	if err != nil {
		return fmt.Errorf("could not load K8S configuration: %w", err)
	}

	kcli, err := kubernetes.NewForConfig(kcfg)
	if err != nil {
		return fmt.Errorf("could not create Kubernetes client: %w", err)
	}

	// Create dependencies
	k8sRepo := storagekubernetes.NewRepository(kcli)
	cachedSecretK8sRepo := storagekubernetes.NewSecretCachedRepository(k8sRepo)

	// Prepare our run entrypoints.
	var g run.Group

	// OS signals.
	{
		sigC := make(chan os.Signal, 1)
		exitC := make(chan struct{})
		signal.Notify(sigC, syscall.SIGTERM, syscall.SIGINT)

		g.Add(
			func() error {
				select {
				case s := <-sigC:
					logger.Infof("signal %s received", s)
					return nil
				case <-exitC:
					return nil
				}
			},
			func(_ error) {
				close(exitC)
			},
		)
	}

	// Main controller for namespaces.
	{
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		handler, err := controllernamespace.NewHandler(controllernamespace.HandlerConfig{
			RunningNamespace:      cmdCfg.NamespaceRunning,
			ImagePullSecretName:   cmdCfg.SecretName,
			SaImagePullSecretName: cmdCfg.SaSecretName,
			K8sRepo:               cachedSecretK8sRepo,
			Logger:                logger,
		})
		if err != nil {
			return fmt.Errorf("could not create namespace controller handler: %w", err)
		}

		retriever, err := controllernamespace.NewRetriever(k8sRepo)
		if err != nil {
			return fmt.Errorf("could not create namespace controller retriever: %w", err)
		}

		ctrl, err := koopercontroller.New(&koopercontroller.Config{
			Handler:              handler,
			Retriever:            retriever,
			Logger:               kooperLogger,
			Name:                 "imagepull-workshop-namespace",
			ConcurrentWorkers:    cmdCfg.Workers,
			ProcessingJobRetries: 2,
			ResyncInterval:       cmdCfg.ResyncInterval,
		})
		if err != nil {
			return fmt.Errorf("could not create namespace controller: %w", err)
		}

		g.Add(
			func() error {
				return ctrl.Run(ctx)
			},
			func(_ error) {
				cancel()
			},
		)
	}

	// Secret cache controller optimization.
	{
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		handler, err := controllersecretcache.NewHandler(cachedSecretK8sRepo, logger)
		if err != nil {
			return fmt.Errorf("could not create secret cache controller handler: %w", err)
		}

		retriever, err := controllersecretcache.NewRetriever(k8sRepo, cmdCfg.NamespaceRunning, cmdCfg.SecretName)
		if err != nil {
			return fmt.Errorf("could not create secret cache controller retriever: %w", err)
		}

		ctrl, err := koopercontroller.New(&koopercontroller.Config{
			Handler:              handler,
			Retriever:            retriever,
			Logger:               kooperLogger,
			Name:                 "imagepull-workshop-secret-cache",
			ConcurrentWorkers:    1,
			ProcessingJobRetries: 1,
			ResyncInterval:       5 * time.Minute,
		})
		if err != nil {
			return fmt.Errorf("could not create secret cache controller: %w", err)
		}

		g.Add(
			func() error {
				return ctrl.Run(ctx)
			},
			func(_ error) {
				cancel()
			},
		)
	}

	err = g.Run()
	if err != nil {
		return err
	}

	return nil
}

// loadKubernetesConfig loads kubernetes configuration based on flags.
func loadKubernetesConfig(cmdCfg CmdConfig) (*rest.Config, error) {
	var cfg *rest.Config

	// If devel mode then use configuration flag path.
	if cmdCfg.Development {
		config, err := clientcmd.BuildConfigFromFlags("", cmdCfg.KubeConfig)
		if err != nil {
			return nil, fmt.Errorf("could not load configuration: %w", err)
		}
		cfg = config
	} else {
		config, err := rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("error loading kubernetes configuration inside cluster, check app is running outside kubernetes cluster or run in development mode: %w", err)
		}
		cfg = config
	}

	// Set better cli rate limiter.
	cfg.QPS = 100
	cfg.Burst = 100

	return cfg, nil
}

func main() {
	ctx := context.Background()

	err := Run(ctx, os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "app error: %s", err)
		os.Exit(1)
	}
}
