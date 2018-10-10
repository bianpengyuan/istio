package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"path/filepath"

	"github.com/spf13/cobra"
	"istio.io/istio/mixer/cmd/shared"
	descriptor2 "istio.io/istio/mixer/pkg/protobuf/descriptor"
	"istio.io/istio/mixer/tools/codegen/pkg/modelgen"
)

func serverGenCmd(rawArgs []string, printf, fatalf shared.FormatFn) *cobra.Command {
	var outServerFile string
	var adapterName string
	var packages []string
	var templateFiles []string
	var adapterPackage string
	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "creates server code from the given templates",
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			outServerFile, err = filepath.Abs(outServerFile)
			if err != nil {
				fatalf("Invalid path %s: %v", outServerFile, err)
			}

			generator := ServerGenerator{OutServerPath: outServerFile, AdapterName: adapterName, AdapterPackage: adapterPackage, TemplatesPath: templateFiles, Packages: packages}
			if err := generator.Generate(); err != nil {
				fatalf("%v", err)
			}
		},
	}
	serverCmd.PersistentFlags().StringArrayVarP(&templateFiles, "templates", "t", nil,
		"supported template names")
	serverCmd.PersistentFlags().StringVar(&outServerFile, "go_out", "./generated.go", "Output "+
		"file path for generated template based go types and interfaces.")
	serverCmd.PersistentFlags().StringArrayVarP(&packages, "packages", "p", []string{},
		"")
	serverCmd.PersistentFlags().StringVar(&adapterPackage, "adapter_package", "", "")
	serverCmd.PersistentFlags().StringVar(&adapterName, "adapter_name", "", "Name of the adapter.")
	return serverCmd
}

// ServerGenerator ...
type ServerGenerator struct {
	OutServerPath  string
	AdapterName    string
	AdapterPackage string
	TemplatesPath  []string
	Packages       []string
}

// Generate ...
func (sg *ServerGenerator) Generate() error {
	var models []*modelgen.Model
	for _, tp := range sg.TemplatesPath {
		fds, err := getFileDescSet(tp)
		if err != nil {
			return fmt.Errorf("cannot parse file '%s' as a FileDescriptorSetProto: %v", tp, err)
		}

		parser := descriptor2.CreateFileDescriptorSetParser(fds, map[string]string{}, "")

		model, err := modelgen.Create(parser)
		models = append(models, model)
		if err != nil {
			return err
		}
	}

	for _, m := range models {
		fmt.Printf("%v %v\n", m.TemplateName, m.InterfaceName)
	}
	f1, err := os.Create(sg.OutServerPath)
	if err != nil {
		return err
	}
	serverFileData, err := sg.getServerGoContent(models, sg.AdapterName, sg.AdapterPackage, noSessionServerTempl, sg.Packages)
	if _, err = f1.Write(serverFileData); err != nil {
		_ = f1.Close()
		_ = os.Remove(f1.Name())
		return err
	}
	return nil
}

func (sg *ServerGenerator) getServerGoContent(models []*modelgen.Model, adapterName, adapterPackage, noSessionTmpl string, packages []string) ([]byte, error) {
	type MS struct {
		AdapterName    string
		AdapterPackage string
		Packages       []string
		Models         []*modelgen.Model
	}
	ms := MS{AdapterName: adapterName, AdapterPackage: adapterPackage, Models: models, Packages: packages}
	fmt.Printf("%+v\n", ms)
	serverTmpl, err := template.New("ProcServer").Funcs(
		template.FuncMap{
			"Capitalize": strings.Title,
			"FindMessage": func(model string, message string) *modelgen.MessageInfo {
				for _, m := range ms.Models {
					if m.TemplateName == model {
						for _, r := range m.ResourceMessages {
							if r.Name == message {
								return &r
							}
						}
					}
				}
				return nil
			},
		}).Parse(noSessionTmpl)
	if err != nil {
		return nil, fmt.Errorf("cannot load template: %v", err)
	}
	serverBuf := new(bytes.Buffer)
	err = serverTmpl.Execute(serverBuf, ms)
	fmt.Printf("%v!\n", err)
	if err != nil {
		return nil, fmt.Errorf("cannot execute the template with the given data: %v", err)
	}
	return serverBuf.Bytes(), nil
}
