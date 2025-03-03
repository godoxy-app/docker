package registry // import "github.com/docker/docker/registry"

import (
	"errors"
	"net"
	"testing"

	"github.com/distribution/reference"
	"github.com/docker/docker/api/types/registry"
	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

// overrideLookupIP overrides net.LookupIP for testing.
func overrideLookupIP(t *testing.T) {
	t.Helper()
	restoreLookup := lookupIP

	// override net.LookupIP
	lookupIP = func(host string) ([]net.IP, error) {
		mockHosts := map[string][]net.IP{
			"":            {net.ParseIP("0.0.0.0")},
			"localhost":   {net.ParseIP("127.0.0.1"), net.ParseIP("::1")},
			"example.com": {net.ParseIP("42.42.42.42")},
			"other.com":   {net.ParseIP("43.43.43.43")},
		}
		if addrs, ok := mockHosts[host]; ok {
			return addrs, nil
		}
		return nil, errors.New("lookup: no such host")
	}
	t.Cleanup(func() {
		lookupIP = restoreLookup
	})
}

func TestParseRepositoryInfo(t *testing.T) {
	type staticRepositoryInfo struct {
		Index         *registry.IndexInfo
		RemoteName    string
		CanonicalName string
		LocalName     string
	}

	expectedRepoInfos := map[string]staticRepositoryInfo{
		"fooo/bar": {
			Index: &registry.IndexInfo{
				Name:     IndexName,
				Official: true,
			},
			RemoteName:    "fooo/bar",
			LocalName:     "fooo/bar",
			CanonicalName: "docker.io/fooo/bar",
		},
		"library/ubuntu": {
			Index: &registry.IndexInfo{
				Name:     IndexName,
				Official: true,
			},
			RemoteName:    "library/ubuntu",
			LocalName:     "ubuntu",
			CanonicalName: "docker.io/library/ubuntu",
		},
		"nonlibrary/ubuntu": {
			Index: &registry.IndexInfo{
				Name:     IndexName,
				Official: true,
			},
			RemoteName:    "nonlibrary/ubuntu",
			LocalName:     "nonlibrary/ubuntu",
			CanonicalName: "docker.io/nonlibrary/ubuntu",
		},
		"ubuntu": {
			Index: &registry.IndexInfo{
				Name:     IndexName,
				Official: true,
			},
			RemoteName:    "library/ubuntu",
			LocalName:     "ubuntu",
			CanonicalName: "docker.io/library/ubuntu",
		},
		"other/library": {
			Index: &registry.IndexInfo{
				Name:     IndexName,
				Official: true,
			},
			RemoteName:    "other/library",
			LocalName:     "other/library",
			CanonicalName: "docker.io/other/library",
		},
		"127.0.0.1:8000/private/moonbase": {
			Index: &registry.IndexInfo{
				Name:     "127.0.0.1:8000",
				Official: false,
			},
			RemoteName:    "private/moonbase",
			LocalName:     "127.0.0.1:8000/private/moonbase",
			CanonicalName: "127.0.0.1:8000/private/moonbase",
		},
		"127.0.0.1:8000/privatebase": {
			Index: &registry.IndexInfo{
				Name:     "127.0.0.1:8000",
				Official: false,
			},
			RemoteName:    "privatebase",
			LocalName:     "127.0.0.1:8000/privatebase",
			CanonicalName: "127.0.0.1:8000/privatebase",
		},
		"localhost:8000/private/moonbase": {
			Index: &registry.IndexInfo{
				Name:     "localhost:8000",
				Official: false,
			},
			RemoteName:    "private/moonbase",
			LocalName:     "localhost:8000/private/moonbase",
			CanonicalName: "localhost:8000/private/moonbase",
		},
		"localhost:8000/privatebase": {
			Index: &registry.IndexInfo{
				Name:     "localhost:8000",
				Official: false,
			},
			RemoteName:    "privatebase",
			LocalName:     "localhost:8000/privatebase",
			CanonicalName: "localhost:8000/privatebase",
		},
		"example.com/private/moonbase": {
			Index: &registry.IndexInfo{
				Name:     "example.com",
				Official: false,
			},
			RemoteName:    "private/moonbase",
			LocalName:     "example.com/private/moonbase",
			CanonicalName: "example.com/private/moonbase",
		},
		"example.com/privatebase": {
			Index: &registry.IndexInfo{
				Name:     "example.com",
				Official: false,
			},
			RemoteName:    "privatebase",
			LocalName:     "example.com/privatebase",
			CanonicalName: "example.com/privatebase",
		},
		"example.com:8000/private/moonbase": {
			Index: &registry.IndexInfo{
				Name:     "example.com:8000",
				Official: false,
			},
			RemoteName:    "private/moonbase",
			LocalName:     "example.com:8000/private/moonbase",
			CanonicalName: "example.com:8000/private/moonbase",
		},
		"example.com:8000/privatebase": {
			Index: &registry.IndexInfo{
				Name:     "example.com:8000",
				Official: false,
			},
			RemoteName:    "privatebase",
			LocalName:     "example.com:8000/privatebase",
			CanonicalName: "example.com:8000/privatebase",
		},
		"localhost/private/moonbase": {
			Index: &registry.IndexInfo{
				Name:     "localhost",
				Official: false,
			},
			RemoteName:    "private/moonbase",
			LocalName:     "localhost/private/moonbase",
			CanonicalName: "localhost/private/moonbase",
		},
		"localhost/privatebase": {
			Index: &registry.IndexInfo{
				Name:     "localhost",
				Official: false,
			},
			RemoteName:    "privatebase",
			LocalName:     "localhost/privatebase",
			CanonicalName: "localhost/privatebase",
		},
		IndexName + "/public/moonbase": {
			Index: &registry.IndexInfo{
				Name:     IndexName,
				Official: true,
			},
			RemoteName:    "public/moonbase",
			LocalName:     "public/moonbase",
			CanonicalName: "docker.io/public/moonbase",
		},
		"index." + IndexName + "/public/moonbase": {
			Index: &registry.IndexInfo{
				Name:     IndexName,
				Official: true,
			},
			RemoteName:    "public/moonbase",
			LocalName:     "public/moonbase",
			CanonicalName: "docker.io/public/moonbase",
		},
		"ubuntu-12.04-base": {
			Index: &registry.IndexInfo{
				Name:     IndexName,
				Official: true,
			},
			RemoteName:    "library/ubuntu-12.04-base",
			LocalName:     "ubuntu-12.04-base",
			CanonicalName: "docker.io/library/ubuntu-12.04-base",
		},
		IndexName + "/ubuntu-12.04-base": {
			Index: &registry.IndexInfo{
				Name:     IndexName,
				Official: true,
			},
			RemoteName:    "library/ubuntu-12.04-base",
			LocalName:     "ubuntu-12.04-base",
			CanonicalName: "docker.io/library/ubuntu-12.04-base",
		},
		"index." + IndexName + "/ubuntu-12.04-base": {
			Index: &registry.IndexInfo{
				Name:     IndexName,
				Official: true,
			},
			RemoteName:    "library/ubuntu-12.04-base",
			LocalName:     "ubuntu-12.04-base",
			CanonicalName: "docker.io/library/ubuntu-12.04-base",
		},
	}

	for reposName, expectedRepoInfo := range expectedRepoInfos {
		named, err := reference.ParseNormalizedNamed(reposName)
		if err != nil {
			t.Error(err)
		}

		repoInfo, err := ParseRepositoryInfo(named)
		if err != nil {
			t.Error(err)
		} else {
			assert.Check(t, is.Equal(repoInfo.Index.Name, expectedRepoInfo.Index.Name), reposName)
			assert.Check(t, is.Equal(reference.Path(repoInfo.Name), expectedRepoInfo.RemoteName), reposName)
			assert.Check(t, is.Equal(reference.FamiliarName(repoInfo.Name), expectedRepoInfo.LocalName), reposName)
			assert.Check(t, is.Equal(repoInfo.Name.Name(), expectedRepoInfo.CanonicalName), reposName)
			assert.Check(t, is.Equal(repoInfo.Index.Official, expectedRepoInfo.Index.Official), reposName)
		}
	}
}

func TestNewIndexInfo(t *testing.T) {
	overrideLookupIP(t)
	testIndexInfo := func(config *serviceConfig, expectedIndexInfos map[string]*registry.IndexInfo) {
		for indexName, expectedIndexInfo := range expectedIndexInfos {
			index := newIndexInfo(config, indexName)
			assert.Check(t, is.Equal(index.Name, expectedIndexInfo.Name), indexName+" name")
			assert.Check(t, is.Equal(index.Official, expectedIndexInfo.Official), indexName+" is official")
			assert.Check(t, is.Equal(index.Secure, expectedIndexInfo.Secure), indexName+" is secure")
			assert.Check(t, is.Equal(len(index.Mirrors), len(expectedIndexInfo.Mirrors)), indexName+" mirrors")
		}
	}

	config := emptyServiceConfig
	var noMirrors []string
	expectedIndexInfos := map[string]*registry.IndexInfo{
		IndexName: {
			Name:     IndexName,
			Official: true,
			Secure:   true,
			Mirrors:  noMirrors,
		},
		"index." + IndexName: {
			Name:     IndexName,
			Official: true,
			Secure:   true,
			Mirrors:  noMirrors,
		},
		"example.com": {
			Name:     "example.com",
			Official: false,
			Secure:   true,
			Mirrors:  noMirrors,
		},
		"127.0.0.1:5000": {
			Name:     "127.0.0.1:5000",
			Official: false,
			Secure:   false,
			Mirrors:  noMirrors,
		},
	}
	testIndexInfo(config, expectedIndexInfos)

	publicMirrors := []string{"http://mirror1.local", "http://mirror2.local"}
	var err error
	config, err = newServiceConfig(ServiceOptions{
		Mirrors:            publicMirrors,
		InsecureRegistries: []string{"example.com"},
	})
	assert.NilError(t, err)

	expectedIndexInfos = map[string]*registry.IndexInfo{
		IndexName: {
			Name:     IndexName,
			Official: true,
			Secure:   true,
			Mirrors:  publicMirrors,
		},
		"index." + IndexName: {
			Name:     IndexName,
			Official: true,
			Secure:   true,
			Mirrors:  publicMirrors,
		},
		"example.com": {
			Name:     "example.com",
			Official: false,
			Secure:   false,
			Mirrors:  noMirrors,
		},
		"example.com:5000": {
			Name:     "example.com:5000",
			Official: false,
			Secure:   true,
			Mirrors:  noMirrors,
		},
		"127.0.0.1": {
			Name:     "127.0.0.1",
			Official: false,
			Secure:   false,
			Mirrors:  noMirrors,
		},
		"127.0.0.1:5000": {
			Name:     "127.0.0.1:5000",
			Official: false,
			Secure:   false,
			Mirrors:  noMirrors,
		},
		"other.com": {
			Name:     "other.com",
			Official: false,
			Secure:   true,
			Mirrors:  noMirrors,
		},
	}
	testIndexInfo(config, expectedIndexInfos)

	config, err = newServiceConfig(ServiceOptions{
		InsecureRegistries: []string{"42.42.0.0/16"},
	})
	assert.NilError(t, err)
	expectedIndexInfos = map[string]*registry.IndexInfo{
		"example.com": {
			Name:     "example.com",
			Official: false,
			Secure:   false,
			Mirrors:  noMirrors,
		},
		"example.com:5000": {
			Name:     "example.com:5000",
			Official: false,
			Secure:   false,
			Mirrors:  noMirrors,
		},
		"127.0.0.1": {
			Name:     "127.0.0.1",
			Official: false,
			Secure:   false,
			Mirrors:  noMirrors,
		},
		"127.0.0.1:5000": {
			Name:     "127.0.0.1:5000",
			Official: false,
			Secure:   false,
			Mirrors:  noMirrors,
		},
		"other.com": {
			Name:     "other.com",
			Official: false,
			Secure:   true,
			Mirrors:  noMirrors,
		},
	}
	testIndexInfo(config, expectedIndexInfos)
}

func TestMirrorEndpointLookup(t *testing.T) {
	containsMirror := func(endpoints []APIEndpoint) bool {
		for _, pe := range endpoints {
			if pe.URL.Host == "my.mirror" {
				return true
			}
		}
		return false
	}
	cfg, err := newServiceConfig(ServiceOptions{
		Mirrors: []string{"https://my.mirror"},
	})
	assert.NilError(t, err)
	s := Service{config: cfg}

	imageName, err := reference.WithName(IndexName + "/test/image")
	if err != nil {
		t.Error(err)
	}
	pushAPIEndpoints, err := s.LookupPushEndpoints(reference.Domain(imageName))
	if err != nil {
		t.Fatal(err)
	}
	if containsMirror(pushAPIEndpoints) {
		t.Fatal("Push endpoint should not contain mirror")
	}

	pullAPIEndpoints, err := s.LookupPullEndpoints(reference.Domain(imageName))
	if err != nil {
		t.Fatal(err)
	}
	if !containsMirror(pullAPIEndpoints) {
		t.Fatal("Pull endpoint should contain mirror")
	}
}

func TestIsSecureIndex(t *testing.T) {
	overrideLookupIP(t)
	tests := []struct {
		addr               string
		insecureRegistries []string
		expected           bool
	}{
		{IndexName, nil, true},
		{"example.com", []string{}, true},
		{"example.com", []string{"example.com"}, false},
		{"localhost", []string{"localhost:5000"}, false},
		{"localhost:5000", []string{"localhost:5000"}, false},
		{"localhost", []string{"example.com"}, false},
		{"127.0.0.1:5000", []string{"127.0.0.1:5000"}, false},
		{"localhost", nil, false},
		{"localhost:5000", nil, false},
		{"127.0.0.1", nil, false},
		{"localhost", []string{"example.com"}, false},
		{"127.0.0.1", []string{"example.com"}, false},
		{"example.com", nil, true},
		{"example.com", []string{"example.com"}, false},
		{"127.0.0.1", []string{"example.com"}, false},
		{"127.0.0.1:5000", []string{"example.com"}, false},
		{"example.com:5000", []string{"42.42.0.0/16"}, false},
		{"example.com", []string{"42.42.0.0/16"}, false},
		{"example.com:5000", []string{"42.42.42.42/8"}, false},
		{"127.0.0.1:5000", []string{"127.0.0.0/8"}, false},
		{"42.42.42.42:5000", []string{"42.1.1.1/8"}, false},
		{"invalid.example.com", []string{"42.42.0.0/16"}, true},
		{"invalid.example.com", []string{"invalid.example.com"}, false},
		{"invalid.example.com:5000", []string{"invalid.example.com"}, true},
		{"invalid.example.com:5000", []string{"invalid.example.com:5000"}, false},
	}
	for _, tc := range tests {
		config, err := newServiceConfig(ServiceOptions{
			InsecureRegistries: tc.insecureRegistries,
		})
		assert.NilError(t, err)

		sec := config.isSecureIndex(tc.addr)
		assert.Equal(t, sec, tc.expected, "isSecureIndex failed for %q %v, expected %v got %v", tc.addr, tc.insecureRegistries, tc.expected, sec)
	}
}
