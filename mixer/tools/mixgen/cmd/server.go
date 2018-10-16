package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"github.com/spf13/cobra"

	"istio.io/istio/mixer/cmd/shared"
	"istio.io/istio/mixer/pkg/protobuf/descriptor"
	"istio.io/istio/mixer/tools/codegen/pkg/modelgen"
)

func serverGenCmd(rawArgs []string, printf, fatalf shared.FormatFn) *cobra.Command {
	var outServerFile string
	var adapterName string
	var templateFiles []string
	var adapterPackage string
	var configPackage string
	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "creates server code from the given templates for an out of process adapter",
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			outServerFile, err = filepath.Abs(outServerFile)
			if err != nil {
				fatalf("Invalid path %s: %v", outServerFile, err)
			}

			generator := ServerGenerator{
				OutServerDir:         outServerFile,
				AdapterName:          adapterName,
				AdapterPackage:       adapterPackage,
				AdapterConfigPackage: configPackage,
				TemplatePaths:        templateFiles,
				TemplatePackages:     make([]string, 0, len(templateFiles)),
				Models:               make([]*modelgen.Model, 0, len(templateFiles)),
			}
			for _, t := range generator.TemplatePaths {
				fds, err := getFileDescSet(t)
				if err != nil {
					fatalf("cannot parse file '%s' as a FileDescriptorSetProto: %v", t, err)
				}

				parser := descriptor.CreateFileDescriptorSetParser(fds, map[string]string{}, "")
				model, err := modelgen.Create(parser)
				if err != nil {
					fatalf("cannot create model for '%s': %v", t, err)
				}
				generator.Models = append(generator.Models, model)
				generator.TemplatePackages = append(generator.TemplatePackages, filepath.Dir(stripGoPath(t)))
			}
			if err := generator.Generate(); err != nil {
				fatalf("%v", err)
			}
		},
	}
	serverCmd.PersistentFlags().StringArrayVarP(&templateFiles, "templates", "t", nil,
		"paths to the descriptor files for all the templates that the adapter supports.")
	serverCmd.PersistentFlags().StringVar(&adapterName, "adapter_name", "",
		"name of the adapter.")
	serverCmd.PersistentFlags().StringVar(&outServerFile, "out_dir", "./",
		"output directory for out of process adapter server code.")
	serverCmd.PersistentFlags().StringVar(&adapterPackage, "adapter_package", "",
		"adapter package, e.g. istio.io/mixer/adapter/prometheus")
	serverCmd.PersistentFlags().StringVar(&configPackage, "config_package", "",
		"adapter config package, e.g. istio.io/mixer/adapter/prometheus/config")
	return serverCmd
}

// ServerGenerator generates server code for out of process adapter.
type ServerGenerator struct {
	OutServerDir         string
	AdapterName          string
	AdapterPackage       string
	AdapterConfigPackage string
	TemplatePaths        []string
	TemplatePackages     []string
	Models               []*modelgen.Model
}

// Generate generates server and main go code for out of process adapter.
func (sg *ServerGenerator) Generate() error {
	if _, err := os.Stat(sg.OutServerDir + "/server"); os.IsNotExist(err) {
		err = os.MkdirAll(sg.OutServerDir+"/server", 0755)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	f1, err := os.Create(sg.OutServerDir + "/server/nosession.go")
	if err != nil {
		return err
	}
	serverFileData, err := sg.getNoSessionServer()
	if _, err = f1.Write(serverFileData); err != nil {
		_ = f1.Close()
		_ = os.Remove(f1.Name())
		return err
	}

	f2, err := os.Create(sg.OutServerDir + "/main.go")
	if err != nil {
		return err
	}

	mainFileData, err := sg.getMainGoContent()
	if _, err = f2.Write(mainFileData); err != nil {
		_ = f2.Close()
		_ = os.Remove(f2.Name())
		return err
	}
	return nil
}

func (sg *ServerGenerator) getNoSessionServer() ([]byte, error) {
	importProto := false
	serverTmpl, err := template.New("ProcServer").Funcs(
		template.FuncMap{
			"Capitalize": strings.Title,
			"FindInterface": func(in modelgen.MessageInfo) string {
				for _, m := range sg.Models {
					if reflect.DeepEqual(m.TemplateMessage, in) {
						return m.InterfaceName
					}
				}
				return ""
			},
			"ConstructDecodeFunc": func(pkg string, goType string) string {
				return "decode" + strings.Title(pkg) + goType[1:]
			},
			"AddProtoToImpt": func() string {
				importProto = true
				return ""
			},
		}).Parse(noSessionServerTempl)
	if err != nil {
		return nil, fmt.Errorf("cannot load template: %v", err)
	}
	serverBuf := new(bytes.Buffer)
	err = serverTmpl.Execute(serverBuf, sg)
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

func (sg *ServerGenerator) getMainGoContent() ([]byte, error) {
	type cmdMain struct {
		AdapterName string
		PackagePath string
	}
	cm := cmdMain{AdapterName: sg.AdapterName, PackagePath: stripGoPath(sg.OutServerDir)}
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

func stripGoPath(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return ""
	}
	goPath := os.Getenv("GOPATH") + "/src/"
	if strings.HasPrefix(absPath, goPath) {
		return strings.TrimPrefix(absPath, goPath)
	}
	return ""
}
