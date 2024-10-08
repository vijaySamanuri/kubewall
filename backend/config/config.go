package config

import (
	"os"
	"path/filepath"
)

var (
	K8SQPS   = 50
	K8SBURST = 50
)

const (
	defaultKubeConfigDir = ".kube"
	appConfigDir         = ".kubewall"
	appKubeConfigDir     = "kubeconfigs"
)

type Env struct {
	KubeConfigs []KubeConfig `json:"kubeconfigs"`
}

type KubeConfig struct {
	Name         string `json:"name"`
	AbsolutePath string `json:"absolutePath"`
}

type AppConfig struct {
	Version    string                     `json:"version"`
	KubeConfig map[string]*KubeConfigInfo `json:"kubeConfigs"`
}

func NewEnv() *Env {
	env := Env{
		KubeConfigs: make([]KubeConfig, 0),
	}
	createEnvDirAndFile()

	return &env
}

func NewAppConfig(version string, k8sqps, k8sburst int) *AppConfig {
	K8SQPS = k8sqps
	K8SBURST = k8sburst
	return &AppConfig{
		Version:    version,
		KubeConfig: make(map[string]*KubeConfigInfo),
	}
}

func (c *AppConfig) LoadAppConfig() {
	c.buildKubeConfigs(filepath.Join(homeDir(), defaultKubeConfigDir))
	c.buildKubeConfigs(filepath.Join(homeDir(), appConfigDir, appKubeConfigDir))

	i, err := LoadInClusterConfig()
	if err == nil {
		c.KubeConfig["incluster"] = &i
	}
}

func (c *AppConfig) buildKubeConfigs(dirPath string) {
	for _, filePath := range readAllFilesInDir(dirPath) {
		if clusters, err := LoadK8ConfigFromFile(filePath); err == nil {
			if len(clusters) > 0 {
				c.KubeConfig[filepath.Base(filePath)] = &KubeConfigInfo{
					Name:         filePath,
					AbsolutePath: filePath,
					FileExists:   true,
					Clusters:     clusters,
				}
			}
		}
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE")
}

func readAllFilesInDir(dirPath string) []string {
	var files []string
	dirFiles, _ := os.ReadDir(dirPath)
	for _, file := range dirFiles {
		files = append(files, filepath.Join(dirPath, file.Name()))
		if file.IsDir() {
			continue
		}
	}
	return files
}

func (c *AppConfig) RemoveKubeConfig(uuid string) error {
	delete(c.KubeConfig, uuid)
	return os.Remove(filepath.Join(homeDir(), appConfigDir, appKubeConfigDir, uuid))
}

func (c *AppConfig) SaveKubeConfig(uuid string) {
	filePath := filepath.Join(homeDir(), appConfigDir, appKubeConfigDir, uuid)
	if clusters, err := LoadK8ConfigFromFile(filePath); err == nil {
		if len(clusters) > 0 {
			c.KubeConfig[filepath.Base(filePath)] = &KubeConfigInfo{
				Name:         filePath,
				AbsolutePath: filePath,
				FileExists:   true,
				Clusters:     clusters,
			}
		}
	}
}

func createEnvDirAndFile() {
	ensureDirExists(filepath.Join(homeDir(), appConfigDir))
	ensureDirExists(filepath.Join(homeDir(), appConfigDir, appKubeConfigDir))
}

func ensureDirExists(dirPath string) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		os.MkdirAll(dirPath, 0755)
	}
}
