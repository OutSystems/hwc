package hwcconfig

import (
	"os"
	"path/filepath"
	"io/ioutil"
)

type HwcConfig struct {
	Instance                      string
	Port                          int
	TempDirectory                 string
	IISCompressedFilesDirectory   string
	ASPCompiledTemplatesDirectory string

	Applications              []*HwcApplication
	AspnetConfigPath          string
	WebConfigPath             string
	ApplicationHostConfigPath string
}

func New(port int, rootPath, tmpPath, contextPath, uuid string, multipleApps bool) (error, *HwcConfig) {
	config := &HwcConfig{
		Instance:                      uuid,
		Port:                          port,
		TempDirectory:                 tmpPath,
		IISCompressedFilesDirectory:   filepath.Join(tmpPath, "IIS Temporary Compressed Files"),
		ASPCompiledTemplatesDirectory: filepath.Join(tmpPath, "ASP Compiled Templates"),
	}

	defaultRootPath := filepath.Join(config.TempDirectory, "wwwroot")
	err := os.MkdirAll(defaultRootPath, 0700)
	if err != nil {
		return err, nil
	}

	configPath := filepath.Join(config.TempDirectory, "config")
	err = os.MkdirAll(configPath, 0700)
	if err != nil {
		return err, nil
	}

	err = os.MkdirAll(config.IISCompressedFilesDirectory, 0700)
	if err != nil {
		return err, nil
	}

	err = os.MkdirAll(config.ASPCompiledTemplatesDirectory, 0700)
	if err != nil {
		return err, nil
	}

	if multipleApps == false {
		config.Applications = NewHwcApplications(defaultRootPath, rootPath, contextPath)
	} else {
		files, err := ioutil.ReadDir(rootPath)
		if err != nil {
			return err, nil
		}

		for _, f := range files {
			if f.IsDir() {

				config.Applications = AppendSliceIfMissing(config.Applications, NewHwcApplications(defaultRootPath , filepath.Join(rootPath, f.Name()),contextPath + f.Name()))
			}
		}
	}
	config.ApplicationHostConfigPath = filepath.Join(configPath, "ApplicationHost.config")
	config.AspnetConfigPath = filepath.Join(configPath, "Aspnet.config")
	config.WebConfigPath = filepath.Join(configPath, "Web.config")

	err = config.generateApplicationHostConfig()
	if err != nil {
		return err, nil
	}

	err = config.generateAspNetConfig()
	if err != nil {
		return err, nil
	}

	err = config.generateWebConfig()
	if err != nil {
		return err, nil
	}

	return nil, config
}
