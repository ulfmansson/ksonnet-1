// Copyright 2017 The kubecfg authors
//
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
package metadata

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"testing"

	"github.com/ksonnet/ksonnet/metadata/app"
	str "github.com/ksonnet/ksonnet/strings"
	"github.com/spf13/afero"
)

const (
	blankSwagger     = "/blankSwagger.json"
	blankSwaggerData = `{
  "swagger": "2.0",
  "info": {
   "title": "Kubernetes",
   "version": "v1.7.0"
  },
  "paths": {
  },
  "definitions": {
  }
}`
	blankK8sLib = `// AUTOGENERATED from the Kubernetes OpenAPI specification. DO NOT MODIFY.
// Kubernetes version: v1.7.0

{
  local hidden = {
  },
}
`
)

func withFs(fn func(afero.Fs)) {
	ogAppLibUpdater := app.LibUpdater
	app.LibUpdater = app.StubUpdateLibData

	defer func() {
		app.LibUpdater = ogAppLibUpdater
	}()

	fs := afero.NewMemMapFs()
	afero.WriteFile(fs, blankSwagger, []byte(blankSwaggerData), os.ModePerm)
	fn(fs)
}

func TestInitSuccess(t *testing.T) {
	withFs(func(fs afero.Fs) {
		specFlag := fmt.Sprintf("file:%s", blankSwagger)

		appPath := "/fromEmptySwagger"
		reg := newMockRegistryManager("incubator")
		_, err := initManager("fromEmptySwagger", appPath, &specFlag, &mockAPIServer, &mockNamespace, reg, fs)
		if err != nil {
			t.Fatalf("Failed to init cluster spec: %v", err)
		}

		// Verify path locations.
		defaultEnvDir := str.AppendToPath(environmentsDir, defaultEnvName)
		paths := []string{
			ksonnetDir,
			libDir,
			componentsDir,
			environmentsDir,
			vendorDir,
			defaultEnvDir,
		}

		for _, p := range paths {
			path := str.AppendToPath(appPath, p)
			exists, err := afero.DirExists(fs, path)
			if err != nil {
				t.Fatalf("Expected to create directory '%s', but failed:\n%v", p, err)
			} else if !exists {
				t.Fatalf("Expected to create directory '%s', but failed", path)
			}
		}

		paths = []string{
			pkgSrcCacheDir,
		}

		usr, err := user.Current()
		if err != nil {
			t.Fatalf("Could not get user information:\n%v", err)
		}
		userRootPath := str.AppendToPath(usr.HomeDir, userKsonnetRootDir)

		for _, p := range paths {
			path := str.AppendToPath(userRootPath, p)
			exists, err := afero.DirExists(fs, path)
			if err != nil {
				t.Fatalf("Expected to create directory '%s', but failed:\n%v", p, err)
			} else if !exists {
				t.Fatalf("Expected to create directory '%s', but failed", path)
			}
		}

		// Verify contents of metadata.
		envPath := str.AppendToPath(appPath, environmentsDir)

		componentParamsPath := str.AppendToPath(appPath, componentsDir, componentParamsFile)
		componentParamsBytes, err := afero.ReadFile(fs, componentParamsPath)
		if err != nil {
			t.Fatalf("Failed to read params.libsonnet file at '%s':\n%v", componentParamsPath, err)
		} else if len(componentParamsBytes) == 0 {
			t.Fatalf("Expected params.libsonnet at '%s' to be non-empty", componentParamsPath)
		}

		baseLibsonnetPath := str.AppendToPath(envPath, baseLibsonnetFile)
		baseLibsonnetBytes, err := afero.ReadFile(fs, baseLibsonnetPath)
		if err != nil {
			t.Fatalf("Failed to read base.libsonnet file at '%s':\n%v", baseLibsonnetPath, err)
		} else if len(baseLibsonnetBytes) == 0 {
			t.Fatalf("Expected base.libsonnet at '%s' to be non-empty", baseLibsonnetPath)
		}

		appYAMLPath := str.AppendToPath(appPath, appYAMLFile)
		appYAMLBytes, err := afero.ReadFile(fs, appYAMLPath)
		if err != nil {
			t.Fatalf("Failed to read app.yaml file at '%s':\n%v", appYAMLPath, err)
		} else if len(appYAMLBytes) == 0 {
			t.Fatalf("Expected app.yaml at '%s' to be non-empty", appYAMLPath)
		}

		registryYAMLPath := str.AppendToPath(appPath, registriesDir, "incubator", "master.yaml")
		registryYAMLBytes, err := afero.ReadFile(fs, registryYAMLPath)
		if err != nil {
			t.Fatalf("Failed to read registry.yaml file at '%s':\n%v", registryYAMLPath, err)
		} else if len(registryYAMLBytes) == 0 {
			t.Fatalf("Expected registry.yaml at '%s' to be non-empty", registryYAMLPath)
		}
	})
}

func TestFindSuccess(t *testing.T) {
	withFs(func(fs afero.Fs) {
		findSuccess := func(t *testing.T, appDir, currDir string) {
			m, err := findManager(currDir, fs)
			if err != nil {
				t.Fatalf("Failed to find manager at path '%s':\n%v", currDir, err)
			} else if m.rootPath != appDir {
				t.Fatalf("Found manager at incorrect path '%s', expected '%s'", m.rootPath, appDir)
			}
		}

		specFlag := fmt.Sprintf("file:%s", blankSwagger)

		appPath := "/findSuccess"
		reg := newMockRegistryManager("incubator")
		_, err := initManager("findSuccess", appPath, &specFlag, &mockAPIServer, &mockNamespace, reg, fs)
		if err != nil {
			t.Fatalf("Failed to init cluster spec: %v", err)
		}

		findSuccess(t, appPath, appPath)

		components := str.AppendToPath(appPath, componentsDir)
		findSuccess(t, appPath, components)

		// Create empty app file.
		appFile := str.AppendToPath(components, "app.jsonnet")
		f, err := fs.OpenFile(appFile, os.O_RDONLY|os.O_CREATE, 0777)
		if err != nil {
			t.Fatalf("Failed to touch app file '%s'\n%v", appFile, err)
		}
		f.Close()

		findSuccess(t, appPath, appFile)
	})
}

func TestLibPaths(t *testing.T) {
	withFs(func(fs afero.Fs) {
		appName := "test-lib-paths"
		expectedVendorPath := path.Join("/", appName, vendorDir)
		expectedEnvPath := path.Join("/", appName, environmentsDir)
		m := mockEnvironments(t, fs, appName)

		envPath, vendorPath := m.LibPaths()
		if envPath != expectedEnvPath {
			t.Fatalf("Expected env path to be:\n  '%s'\n, got:\n  '%s'", expectedEnvPath, envPath)
		}
		if vendorPath != expectedVendorPath {
			t.Fatalf("Expected vendor lib path to be:\n  '%s'\n, got:\n  '%s'", expectedVendorPath, vendorPath)
		}
	})
}

func TestMakeEnvPaths(t *testing.T) {
	withFs(func(fs afero.Fs) {
		appName := "test-env-paths"
		expectedMainPath := path.Join("/", appName, environmentsDir, mockEnvName, envFileName)
		expectedParamsPath := path.Join("/", appName, environmentsDir, mockEnvName, paramsFileName)
		m := mockEnvironments(t, fs, appName)

		mainPath, paramsPath := m.makeEnvPaths(mockEnvName)

		if mainPath != expectedMainPath {
			t.Fatalf("Expected environment main path to be:\n  '%s'\n, got:\n  '%s'", expectedMainPath, mainPath)
		}
		if paramsPath != expectedParamsPath {
			t.Fatalf("Expected environment params path to be:\n  '%s'\n, got:\n  '%s'", expectedParamsPath, paramsPath)
		}
	})
}

func TestFindFailure(t *testing.T) {
	withFs(func(fs afero.Fs) {
		findFailure := func(t *testing.T, currDir string) {
			_, err := findManager(currDir, fs)
			if err == nil {
				t.Fatalf("Expected to fail to find ksonnet app in '%s', but succeeded", currDir)
			}
		}

		findFailure(t, "/")
		findFailure(t, "/fakePath")
		findFailure(t, "")
	})
}

func TestDoubleNewFailure(t *testing.T) {
	withFs(func(fs afero.Fs) {
		specFlag := fmt.Sprintf("file:%s", blankSwagger)

		appPath := "/doubleNew"
		reg := newMockRegistryManager("incubator")
		_, err := initManager("doubleNew", appPath, &specFlag, &mockAPIServer, &mockNamespace, reg, fs)
		if err != nil {
			t.Fatalf("Failed to init cluster spec: %v", err)
		}

		targetErr := fmt.Sprintf("Could not create app; directory '%s' already exists", appPath)
		_, err = initManager("doubleNew", appPath, &specFlag, &mockAPIServer, &mockNamespace, reg, fs)
		if err == nil || err.Error() != targetErr {
			t.Fatalf("Expected to fail to create app with message '%s', got '%s'", targetErr, err.Error())
		}
	})
}
