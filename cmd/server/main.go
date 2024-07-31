package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/pigeonligh/kube-indexer/pkg/cache"
	"github.com/pigeonligh/kube-indexer/pkg/dataprocessor"
	"github.com/pigeonligh/kube-indexer/pkg/server"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func main() {
	configFlags := &genericclioptions.ConfigFlags{
		KubeConfig: new(string),
	}
	var templateFile string
	var restfulPort int

	cmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			template, err := readTemplate(templateFile)
			if err != nil {
				panic(err)
			}

			c, err := cache.New(configFlags, template.ForList()...)
			if err != nil {
				panic(err)
			}
			c.Init()
			go func() {
				_ = c.Run(cmd.Context())
			}()

			c.WaitForCacheSync(cmd.Context())

			s := server.New(c, restfulPort, template)
			if err = s.Run(cmd.Context()); err != nil {
				panic(err)
			}
		},
	}

	configFlags.AddFlags(cmd.Flags())
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
