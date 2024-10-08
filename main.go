package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/utils/ptr"

	"github.com/pigeonligh/kube-indexer/pkg/dataprocessor"
	"github.com/pigeonligh/kube-indexer/pkg/server"
)

//go:embed all:web/build
var assets embed.FS

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func main() {
	var kubeconfig string
	var templateFile string
	var restfulPort int

	cmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			template, err := readTemplate(templateFile)
			if err != nil {
				panic(err)
			}
			backend := gin.New()
			kcm := getKubeConfigMap(kubeconfig)

			clusterGroup := backend.Group("api").Group("cluster")
			clusterGroup.GET("", func(ctx *gin.Context) {
				names := make([]string, 0)
				for name := range kcm {
					names = append(names, name)
				}
				ctx.JSON(http.StatusOK, names)
			})
			for name, config := range kcm {
				go runServer(cmd.Context(), name, config, clusterGroup, template)
			}

			frontend := gin.New()
			frontend.StaticFS("/", http.FS(subpathFS{assets, "web/build"}))

			r := gin.Default()
			r.Any("*path", func(ctx *gin.Context) {
				path := ctx.Param("path")
				if strings.HasPrefix(path, "/api/") {
					backend.ServeHTTP(ctx.Writer, ctx.Request)
				} else {
					frontend.ServeHTTP(ctx.Writer, ctx.Request)
				}
			})
			r.Run(fmt.Sprintf(":%v", restfulPort))
		},
	}

	cmd.Flags().StringVar(&kubeconfig, "kubeconfig", os.Getenv("KUBECONFIG"), "kubeconfig")
	cmd.Flags().StringVar(&templateFile, "template", "", "template")
	cmd.Flags().IntVar(&restfulPort, "port", 8080, "port")
	_ = cmd.Execute()
}

func readTemplate(templateFile string) (*dataprocessor.Template, error) {
	template := &dataprocessor.Template{}
	data, err := os.ReadFile(templateFile)
	if err != nil {
		return nil, fmt.Errorf("reading template file: %w", err)
	}
	err = yaml.Unmarshal(data, template)
	if err != nil {
		return nil, fmt.Errorf("parse template file: %w", err)
	}
	return template, nil
}

func getKubeConfigMap(kubeconfig string) map[string]*genericclioptions.ConfigFlags {
	if kubeconfig == "" {
		return map[string]*genericclioptions.ConfigFlags{
			"default": {},
		}
	}

	configflags := &genericclioptions.ConfigFlags{
		KubeConfig: &kubeconfig,
	}
	config, err := configflags.ToRawKubeConfigLoader().RawConfig()
	if err != nil {
		return map[string]*genericclioptions.ConfigFlags{
			"default": {},
		}
	}
	ret := make(map[string]*genericclioptions.ConfigFlags)
	for name := range config.Contexts {
		ret[name] = &genericclioptions.ConfigFlags{
			KubeConfig: &kubeconfig,
			Context:    ptr.To(name),
		}
	}
	return ret
}

func runServer(ctx context.Context, name string, config *genericclioptions.ConfigFlags, router gin.IRouter, template *dataprocessor.Template) {
	for {
		select {
		case <-ctx.Done():
			return

		default:
			s := server.New(config, router.Group(name), template)
			if err := s.Init(ctx); err != nil {
				fmt.Fprintf(os.Stderr, "Init %v failed: %v\n", name, err)
				time.Sleep(time.Second * 10)
				continue
			}
			if err := s.Run(ctx); err != nil {
				fmt.Fprintf(os.Stderr, "Run %v failed: %v\n", name, err)
				time.Sleep(time.Second * 10)
				continue
			}
		}
	}
}

type subpathFS struct {
	fs.FS
	subpath string
}

func (f subpathFS) Open(name string) (fs.File, error) {
	path := filepath.Join(f.subpath, name)
	if file, err := f.FS.Open(path); err == nil {
		return file, err
	}
	return f.FS.Open(filepath.Join(f.subpath, "index.html"))
}
