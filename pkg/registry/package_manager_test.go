// Copyright 2018 The ksonnet authors
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

package registry

import (
	"testing"

	"github.com/ksonnet/ksonnet/pkg/app"
	amocks "github.com/ksonnet/ksonnet/pkg/app/mocks"
	"github.com/ksonnet/ksonnet/pkg/parts"
	"github.com/ksonnet/ksonnet/pkg/pkg"
	"github.com/ksonnet/ksonnet/pkg/prototype"
	"github.com/ksonnet/ksonnet/pkg/util/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_packageManager_Find(t *testing.T) {
	test.WithApp(t, "/app", func(a *amocks.App, fs afero.Fs) {

		test.StageDir(t, fs, "incubator/apache", "/work/apache")
		test.StageDir(t, fs, "incubator/apache", "/app/vendor/incubator/apache")

		a.On("VendorPath").Return("/app/vendor")

		registries := app.RegistryConfigs{
			"incubator": &app.RegistryConfig{
				Protocol: "fs",
				URI:      "/work",
			},
		}

		a.On("Registries").Return(registries, nil)

		libraries := app.LibraryConfigs{
			"apache": &app.LibraryConfig{},
		}

		a.On("Libraries").Return(libraries, nil)

		pm := NewPackageManager(a)

		p, err := pm.Find("incubator/apache")
		require.NoError(t, err)

		require.Equal(t, "apache", p.Name())
	})
}

func Test_packageManager_Packages(t *testing.T) {
	test.WithApp(t, "/app", func(a *amocks.App, fs afero.Fs) {
		makePkg := func(name, registry, version string) pkg.Package {
			p, err := pkg.NewLocal(a, name, registry, version, &pkg.DefaultInstallChecker{App: a})
			require.NoErrorf(t, err, "creating package %s/%s@%s", registry, name, version)
			return p
		}

		test.StageDir(t, fs, "incubator/apache", "/work/apache")
		test.StageDir(t, fs, "incubator/apache", "/app/vendor/incubator/apache")
		test.StageDir(t, fs, "incubator/nginx", "/app/vendor/incubator/nginx@2.0.0")
		test.StageDir(t, fs, "incubator/nginx", "/app/vendor/incubator/nginx@1.2.3")

		a.On("VendorPath").Return("/app/vendor")

		registries := app.RegistryConfigs{
			"incubator": &app.RegistryConfig{
				Protocol: "fs",
				URI:      "/work",
			},
		}

		a.On("Registries").Return(registries, nil)

		libraries := app.LibraryConfigs{
			"apache": &app.LibraryConfig{
				Registry: "incubator",
				Name:     "apache",
			},
			"nginx": &app.LibraryConfig{
				Registry: "incubator",
				Name:     "nginx",
				Version:  "2.0.0",
			},
		}

		a.On("Libraries").Return(libraries, nil)

		envLibraries := app.LibraryConfigs{
			"nginx": &app.LibraryConfig{
				Registry: "incubator",
				Name:     "nginx",
				Version:  "1.2.3",
			},
		}
		environments := app.EnvironmentConfigs{
			"default": &app.EnvironmentConfig{
				Name:      "default",
				Libraries: envLibraries,
			},
		}
		a.On("Environments").Return(environments, nil)

		// Expect global libraries + envLibraries
		expected := make([]pkg.Package, 0, len(libraries)+len(envLibraries))
		for _, l := range libraries {
			p := makePkg(l.Name, l.Registry, l.Version)
			expected = append(expected, p)
		}
		for _, l := range envLibraries {
			p := makePkg(l.Name, l.Registry, l.Version)
			expected = append(expected, p)
		}

		pm := NewPackageManager(a)

		packages, err := pm.Packages()
		require.NoError(t, err)

		assert.Len(t, packages, len(libraries)+len(envLibraries))
		assert.Subset(t, expected, packages)
	})
}

func Test_packageManager_Prototypes(t *testing.T) {
	test.WithApp(t, "/app", func(a *amocks.App, fs afero.Fs) {
		makePkg := func(name, registry, version string) pkg.Package {
			p, err := pkg.NewLocal(a, name, registry, version, &pkg.DefaultInstallChecker{App: a})
			require.NoErrorf(t, err, "creating package %s/%s@%s", registry, name, version)
			return p
		}

		test.StageDir(t, fs, "incubator/apache", "/work/apache")
		test.StageDir(t, fs, "incubator/apache", "/app/vendor/incubator/apache")
		test.StageDir(t, fs, "incubator/apache", "/app/vendor/incubator/apache@1.2.3")

		a.On("VendorPath").Return("/app/vendor")

		pkgs := []pkg.Package{
			makePkg("apache", "incubator", ""),
			makePkg("apache", "incubator", "2.0.1"),
			makePkg("apache", "incubator", "1.2.3"),
		}

		pm := packageManager{
			app:            a,
			InstallChecker: &pkg.DefaultInstallChecker{App: a},
			packagesFn: func() ([]pkg.Package, error) {
				return pkgs, nil
			},
		}

		protos, err := pm.Prototypes()
		require.NoError(t, err)

		// We expect the prototype to be retuned by only one of the packages
		require.Len(t, protos, 1)
		assert.Equal(t, "2.0.1", protos[0].Version)
	})
}

func Test_latestPrototype(t *testing.T) {
	protos := prototype.Prototypes{
		&prototype.Prototype{
			Version: "",
		},
		&prototype.Prototype{
			Version: "2.4.5",
		},
		&prototype.Prototype{
			Version: "v2.0.5",
		},
		&prototype.Prototype{
			Version: "1.2.3",
		},
	}

	p := latestPrototype(protos)
	assert.Equal(t, "2.4.5", p.Version)
}

func Test_latestPrototype_Non_Semver(t *testing.T) {
	protos := prototype.Prototypes{
		&prototype.Prototype{
			Version: "",
		},
		&prototype.Prototype{
			Version: "semanticLargest",
		},
		&prototype.Prototype{
			Version: "abcd",
		},
		&prototype.Prototype{
			Version: "notsemver",
		},
	}

	p := latestPrototype(protos)
	assert.Equal(t, "semanticLargest", p.Version)
}

func Test_remotePackage(t *testing.T) {
	rp := &remotePackage{
		registryName: "registry-name",
		partConfig: &parts.Spec{
			Name:        "pkg-name",
			Description: "description",
		},
	}

	i, err := rp.IsInstalled()
	require.NoError(t, err)

	protos, err := rp.Prototypes()
	require.NoError(t, err)

	assert.Equal(t, "pkg-name", rp.Name())
	assert.Equal(t, "registry-name", rp.RegistryName())
	assert.Equal(t, "description", rp.Description())
	assert.False(t, i)
	assert.Empty(t, protos)
}

func Test_IsInstalled(t *testing.T) {
	libs := app.LibraryConfigs{
		"mysql": &app.LibraryConfig{
			Name:     "mysql",
			Registry: "incubator",
			Version:  "1.2.3",
		},
		"consul": &app.LibraryConfig{
			Name:     "consul",
			Registry: "incubator",
			Version:  "0.6.4",
		},
		"unversioned": &app.LibraryConfig{
			Name:     "unversioned",
			Registry: "helm",
			Version:  "",
		},
	}

	envs := app.EnvironmentConfigs{
		"default": &app.EnvironmentConfig{
			Name: "default",
			Libraries: app.LibraryConfigs{
				"mysql": &app.LibraryConfig{
					Name:     "mysql",
					Registry: "incubator",
					Version:  "4.5.6",
				},
				"nginx": &app.LibraryConfig{
					Name:     "nginx",
					Registry: "incubator",
					Version:  "0.0.1",
				},
			},
		},
	}

	tests := []struct {
		name      string
		libraries app.LibraryConfigs
		envs      app.EnvironmentConfigs
		desc      pkg.Descriptor
		expected  bool
	}{
		{
			name:      "fully qualified",
			libraries: libs,
			envs:      envs,
			desc:      pkg.Descriptor{Name: "mysql", Registry: "incubator", Version: "1.2.3"},
			expected:  true,
		},
		{
			name:      "registry/name, any version",
			libraries: libs,
			envs:      envs,
			desc:      pkg.Descriptor{Name: "mysql", Registry: "incubator", Version: ""},
			expected:  true,
		},
		{
			name:      "just name, any registry, any version",
			libraries: libs,
			envs:      envs,
			desc:      pkg.Descriptor{Name: "mysql", Registry: "", Version: ""},
			expected:  true,
		},
		{
			name:      "wrong version, qualified registry",
			libraries: libs,
			envs:      envs,
			desc:      pkg.Descriptor{Name: "mysql", Registry: "incubator", Version: "9.9.9"},
			expected:  false,
		},
		{
			name:      "wrong version, any registry",
			libraries: libs,
			envs:      envs,
			desc:      pkg.Descriptor{Name: "mysql", Registry: "", Version: "9.9.9"},
			expected:  false,
		},
		{
			name:      "only in environment",
			libraries: libs,
			envs:      envs,
			desc:      pkg.Descriptor{Name: "nginx", Registry: "incubator", Version: ""},
			expected:  true,
		},
		{
			name:      "only in globals",
			libraries: libs,
			envs:      envs,
			desc:      pkg.Descriptor{Name: "consul", Registry: "incubator", Version: ""},
			expected:  true,
		},
		{
			name:      "wrong name",
			libraries: libs,
			envs:      envs,
			desc:      pkg.Descriptor{Name: "fake", Registry: "", Version: ""},
			expected:  false,
		},
		{
			name:      "unversioned library",
			libraries: libs,
			envs:      envs,
			desc:      pkg.Descriptor{Name: "unversioned", Registry: "", Version: ""},
			expected:  true,
		},
		{
			name:      "unversioned library, versioned search",
			libraries: libs,
			envs:      envs,
			desc:      pkg.Descriptor{Name: "unversioned", Registry: "", Version: "uh-oh"},
			expected:  false,
		},
	}

	for _, tc := range tests {
		test.WithApp(t, "/app", func(a *amocks.App, fs afero.Fs) {
			a.On("Libraries").Return(tc.libraries, nil)
			a.On("Environments").Return(tc.envs, nil)
			pm := NewPackageManager(a)
			actual, err := pm.IsInstalled(tc.desc)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, actual, tc.name)
		})
	}
}
