package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Vaziria/coin_cli/lib/commonlib"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

type Compose struct {
	Version  string              `yaml:"version"`
	Services map[string]*Service `yaml:"services"`
	Networks map[string]*Network `yaml:"networks"`
}
type Build struct {
	Context    string `yaml:"context"`
	Dockerfile string `yaml:"dockerfile"`
}

type Config struct {
	Subnet string `yaml:"subnet"`
}
type Ipam struct {
	Driver string   `yaml:"driver"`
	Config []Config `yaml:"config"`
}
type Network struct {
	EnableIpv6 bool `yaml:"enable_ipv6"`
	Ipam       Ipam `yaml:"ipam"`
}

type Service struct {
	ContainerName string                     `yaml:"container_name"`
	Build         *Build                     `yaml:"build"`
	Image         string                     `yaml:"image"`
	Restart       string                     `yaml:"restart"`
	Volumes       []string                   `yaml:"volumes"`
	DependsOn     []string                   `yaml:"depends_on,omitempty"`
	Expose        []int                      `yaml:"expose"`
	Environment   []string                   `yaml:"environment"`
	Networks      map[string]*ServiceNetwork `yaml:"networks"`
}

type ServiceNetwork struct {
	Ipv4Address string `yaml:"ipv4_address"`
}

func CreateCompose() (*Compose, error) {

	config, err := GetConfiguration()

	if err != nil {
		return nil, err
	}

	ip, err := GettingWslInterface("vEthernet (WSL)")

	if err != nil {
		return nil, err
	}

	ghToken := config.GithubCoin.Token

	envToken := fmt.Sprintf("GH_TOKEN=%s", ghToken)

	displayIp := fmt.Sprintf("DISPLAY=%s:0.0", ip.String())

	port_coin := config.Port
	coinDir := fmt.Sprintf("./coindir/virtuoso_deployer:/root/.%score", config.GithubCoin.Name)
	entryPoint := fmt.Sprintf("COIN_ENTRYPOINT=qt/%sd", config.GithubCoin.Name)

	compose := Compose{
		Version: "3.1",
		Services: map[string]*Service{
			"deployer": {
				ContainerName: "virtuoso_deployer",
				Build: &Build{
					Context:    "./",
					Dockerfile: "Dockerfile",
				},
				Image:   "virtuedeployer",
				Restart: "always",
				Volumes: []string{
					"./virtuosocoin:/root/coin",
					"./workcoin:/root/workcoin",
					"./tool:/root/tool",
					"./secret_key:/root/secretkey",
					"./cache_build:/root/.ccache",
					coinDir,
				},
				Expose: []int{
					port_coin,
				},

				Environment: []string{
					displayIp,
					envToken,
					"QT_XKB_CONFIG_ROOT=/usr/share/X11/xkb",
					entryPoint,
				},

				Networks: map[string]*ServiceNetwork{
					"deployer_network": {
						Ipv4Address: "192.168.100.11",
					},
				},
			},
		},
		Networks: map[string]*Network{
			"deployer_network": {
				EnableIpv6: false,
				Ipam: Ipam{
					Driver: "default",
					Config: []Config{
						{
							Subnet: "192.168.100.0/24",
						},
					},
				},
			},
		},
	}

	// inisiasi workdir coin
	for _, dir := range config.WorkCoinVolume {
		pat := fmt.Sprintf("./workcoin/%s:/root/%s", dir, dir)
		compose.Services["deployer"].Volumes = append(compose.Services["deployer"].Volumes, pat)
	}

	return &compose, nil
}

func GenerateCompose(ctx *cli.Context) error {
	dirctx := commonlib.NewBaseLocation(ctx)

	copys := map[string]func(dirctx *commonlib.BaseLocation) error{
		"Copy Docker File":             CopyDockerFile,
		"Copy Shell Daemon":            CopyDaemonSh,
		"Copy Script Berkeley Install": ScriptBerkeley,
		"Copy script tools": func(dirctx *commonlib.BaseLocation) error {
			names := []string{
				"tool/installdepend.sh",
				"tool/renamefile.py",
				"tool/charpreffix.py",
				"tool/create_release.sh",
				"tool/generate_block_info.py",
				"tool/path_private_key.py",
			}

			for _, name := range names {
				err := CopyFile(dirctx, name)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	for key, callback := range copys {
		err := callback(dirctx)
		if err != nil {
			log.Println(key + " gagal")
			return err
		}
	}

	fname := dirctx.Path("docker-compose.yaml")

	file, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	compose, err := CreateCompose()
	if err != nil {
		return err
	}

	err = yaml.NewEncoder(file).Encode(compose)

	if err != nil {
		return err
	}

	return nil

}

func CopyFile(dirctx *commonlib.BaseLocation, fname string) error {
	data, err := Templates.ReadFile("templates/" + fname)
	if err != nil {
		return err
	}

	content := strings.ReplaceAll(string(data), "\r", "")

	dstname := dirctx.Path(fname)
	file, err := os.OpenFile(dstname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)

	return err

}

func CopyDockerFile(dirctx *commonlib.BaseLocation) error {
	data, err := Templates.ReadFile("templates/Dockerfile")
	if err != nil {
		return err
	}
	fname := dirctx.Path("Dockerfile")
	file, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(string(data))

	return err
}

func CopyDaemonSh(dirctx *commonlib.BaseLocation) error {
	data, err := Templates.ReadFile("templates/start_daemon.sh")
	if err != nil {
		return err
	}
	fname := dirctx.Path("tool", "start_daemon.sh")
	file, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(string(data))

	return err
}

func ScriptBerkeley(dirctx *commonlib.BaseLocation) error {
	data, err := Templates.ReadFile("templates/scripts/install-berkeley.sh")
	if err != nil {
		return err
	}

	os.MkdirAll(dirctx.Path("scripts"), os.ModeDir)

	fname := dirctx.Path("scripts/install-berkeley.sh")
	file, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(string(data))

	return err
}
