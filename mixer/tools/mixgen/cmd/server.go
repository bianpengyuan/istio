package cmd

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
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
	var configPackage string
	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "creates server code from the given templates",
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			outServerFile, err = filepath.Abs(outServerFile)
			if err != nil {
				fatalf("Invalid path %s: %v", outServerFile, err)
			}

			generator := ServerGenerator{
				OutServerDir:   outServerFile,
				AdapterName:    adapterName,
				AdapterPackage: adapterPackage,
				TemplatesPath:  templateFiles,
				Packages:       packages,
				ConfigPackage:  configPackage,
			}
			if err := generator.Generate(); err != nil {
				fatalf("%v", err)
			}
		},
	}
	serverCmd.PersistentFlags().StringArrayVarP(&templateFiles, "templates", "t", nil,
		"supported template names")
	serverCmd.PersistentFlags().StringVar(&outServerFile, "go_out", "./generated.go", "Output "+
		"file path for generated template based go types and interfaces.")
	serverCmd.PersistentFlags().StringVar(&configPackage, "config_package", "", "")
	serverCmd.PersistentFlags().StringArrayVarP(&packages, "packages", "p", []string{},
		"")
	serverCmd.PersistentFlags().StringVar(&adapterPackage, "adapter_package", "", "")
	serverCmd.PersistentFlags().StringVar(&adapterName, "adapter_name", "", "Name of the adapter.")
	return serverCmd
}

// ServerGenerator ...
type ServerGenerator struct {
	OutServerDir      string
	AdapterName       string
	AdapterPackage    string
	AdapterConfigPath string
	TemplatesPath     []string
	Packages          []string
	ConfigPackage     string
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
	if _, err := os.Stat(sg.OutServerDir + "/nosession"); os.IsNotExist(err) {
		err = os.MkdirAll(sg.OutServerDir+"/nosession", os.ModePerm)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	f1, err := os.Create(sg.OutServerDir + "/nosession/nosession.go")
	if err != nil {
		return err
	}
	serverFileData, err := sg.getServerGoContent(models, sg.AdapterName, sg.AdapterPackage, noSessionServerTempl, sg.Packages)
	if _, err = f1.Write(serverFileData); err != nil {
		_ = f1.Close()
		_ = os.Remove(f1.Name())
		return err
	}

	f2, err := os.Create(sg.OutServerDir + "/main.go")
	if err != nil {
		return err
	}

	mainFileData, err := sg.getMainGoContent(sg.AdapterName, sg.OutServerDir, oopMainTempl)
	if _, err = f2.Write(mainFileData); err != nil {
		_ = f2.Close()
		_ = os.Remove(f2.Name())
		return err
	}
	return nil
}

func (sg *ServerGenerator) getMainGoContent(adapterName, outServerPath, oopMainTempl string) ([]byte, error) {
	type cmdMain struct {
		AdapterName string
		PackagePath string
	}
	pp := strings.TrimPrefix(outServerPath, os.Getenv("GOPATH")+"/src/")
	cm := cmdMain{AdapterName: adapterName, PackagePath: pp}
	mainTmpl, err := template.New("ProcMain").Funcs(
		template.FuncMap{
			"Capitalize": strings.Title,
		}).Parse(oopMainTempl)
	mainBuf := new(bytes.Buffer)
	err = mainTmpl.Execute(mainBuf, cm)
	if err != nil {
		return nil, fmt.Errorf("cannot execute the template with the given data: %v", err)
	}
	return mainBuf.Bytes(), nil
}

func (sg *ServerGenerator) getServerGoContent(models []*modelgen.Model, adapterName, adapterPackage, noSessionTmpl string, packages []string) ([]byte, error) {
	type MS struct {
		AdapterName    string
		AdapterPackage string
		Packages       []string
		ConfigPath     string
		Models         []*modelgen.Model
	}
	ms := MS{AdapterName: adapterName, AdapterPackage: adapterPackage, Models: models, Packages: packages, ConfigPath: sg.ConfigPackage}
	importProto := false
	serverTmpl, err := template.New("ProcServer").Funcs(
		template.FuncMap{
			"Capitalize": strings.Title,
			"FindInterface": func(in modelgen.MessageInfo) string {
				for _, m := range ms.Models {
					if reflect.DeepEqual(m.TemplateMessage, in) {
						return m.InterfaceName
					}
				}
				return ""
			},
			"ConstructDecodeFunc": func(pkg string, goType string) string {
				return "decode" + strings.Title(pkg) + goType[1:]
			},
			"FindMessage": func(in modelgen.MessageInfo, message string) modelgen.MessageInfo {
				for _, m := range ms.Models {
					if reflect.DeepEqual(m.TemplateMessage, in) {
						for _, r := range m.ResourceMessages {
							if "*"+r.Name == message {
								return r
							}
						}
					}
					for _, r := range m.ResourceMessages {
						if reflect.DeepEqual(r, in) {
							for _, r := range m.ResourceMessages {
								if "*"+r.Name == message {
									return r
								}
							}
						}
					}
				}
				return modelgen.MessageInfo{}
			},
			"AddProtoToImpt": func() string {
				importProto = true
				return ""
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
	var retBytes []byte
	if importProto {
		retBytes = bytes.Replace(serverBuf.Bytes(), []byte("$$additional_imports$$"), []byte("proto \"github.com/gogo/protobuf/types\""), 1)
	} else {
		retBytes = bytes.Replace(serverBuf.Bytes(), []byte("$$additional_imports$$"), []byte(""), 1)
	}
	return retBytes, nil
}
