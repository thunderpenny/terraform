package rpc

import (
	"errors"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/terraform"
)

func TestResourceProvider_configure(t *testing.T) {
	p := new(terraform.MockResourceProvider)
	client, server := testClientServer(t)
	name, err := Register(server, p)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	provider := &ResourceProvider{Client: client, Name: name}

	// Configure
	config := map[string]interface{}{"foo": "bar"}
	e := provider.Configure(config)
	if !p.ConfigureCalled {
		t.Fatal("configure should be called")
	}
	if !reflect.DeepEqual(p.ConfigureConfig, config) {
		t.Fatalf("bad: %#v", p.ConfigureConfig)
	}
	if e != nil {
		t.Fatalf("bad: %#v", e)
	}
}

func TestResourceProvider_configure_errors(t *testing.T) {
	p := new(terraform.MockResourceProvider)
	client, server := testClientServer(t)
	name, err := Register(server, p)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	provider := &ResourceProvider{Client: client, Name: name}

	p.ConfigureReturnError = errors.New("foo")

	// Configure
	config := map[string]interface{}{"foo": "bar"}
	e := provider.Configure(config)
	if !p.ConfigureCalled {
		t.Fatal("configure should be called")
	}
	if !reflect.DeepEqual(p.ConfigureConfig, config) {
		t.Fatalf("bad: %#v", p.ConfigureConfig)
	}
	if e == nil {
		t.Fatal("should have error")
	}
	if e.Error() != "foo" {
		t.Fatalf("bad: %s", e)
	}
}

func TestResourceProvider_configure_warnings(t *testing.T) {
	p := new(terraform.MockResourceProvider)
	client, server := testClientServer(t)
	name, err := Register(server, p)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	provider := &ResourceProvider{Client: client, Name: name}

	// Configure
	config := map[string]interface{}{"foo": "bar"}
	e := provider.Configure(config)
	if !p.ConfigureCalled {
		t.Fatal("configure should be called")
	}
	if !reflect.DeepEqual(p.ConfigureConfig, config) {
		t.Fatalf("bad: %#v", p.ConfigureConfig)
	}
	if e != nil {
		t.Fatalf("bad: %#v", e)
	}
}

func TestResourceProvider_diff(t *testing.T) {
	p := new(terraform.MockResourceProvider)
	client, server := testClientServer(t)
	name, err := Register(server, p)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	provider := &ResourceProvider{Client: client, Name: name}

	p.DiffReturn = &terraform.ResourceDiff{
		Attributes: map[string]*terraform.ResourceAttrDiff{
			"foo": &terraform.ResourceAttrDiff{
				Old: "",
				New: "bar",
			},
		},
	}

	// Diff
	state := &terraform.ResourceState{}
	config := map[string]interface{}{"foo": "bar"}
	diff, err := provider.Diff(state, config)
	if !p.DiffCalled {
		t.Fatal("diff should be called")
	}
	if !reflect.DeepEqual(p.DiffDesired, config) {
		t.Fatalf("bad: %#v", p.DiffDesired)
	}
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}
	if !reflect.DeepEqual(p.DiffReturn, diff) {
		t.Fatalf("bad: %#v", diff)
	}
}

func TestResourceProvider_diff_error(t *testing.T) {
	p := new(terraform.MockResourceProvider)
	client, server := testClientServer(t)
	name, err := Register(server, p)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	provider := &ResourceProvider{Client: client, Name: name}

	p.DiffReturnError = errors.New("foo")

	// Diff
	state := &terraform.ResourceState{}
	config := map[string]interface{}{"foo": "bar"}
	diff, err := provider.Diff(state, config)
	if !p.DiffCalled {
		t.Fatal("diff should be called")
	}
	if !reflect.DeepEqual(p.DiffDesired, config) {
		t.Fatalf("bad: %#v", p.DiffDesired)
	}
	if err == nil {
		t.Fatal("should have error")
	}
	if diff != nil {
		t.Fatal("should not have diff")
	}
}

func TestResourceProvider_resources(t *testing.T) {
	p := new(terraform.MockResourceProvider)
	client, server := testClientServer(t)
	name, err := Register(server, p)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	provider := &ResourceProvider{Client: client, Name: name}

	expected := []terraform.ResourceType{
		{"foo"},
		{"bar"},
	}

	p.ResourcesReturn = expected

	// Resources
	result := provider.Resources()
	if !p.ResourcesCalled {
		t.Fatal("resources should be called")
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("bad: %#v", result)
	}
}
