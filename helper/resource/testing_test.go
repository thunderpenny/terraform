package resource

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/terraform"
)

func init() {
	testTesting = true

	if err := os.Setenv(TestEnvVar, "1"); err != nil {
		panic(err)
	}
}

func TestTest(t *testing.T) {
	mp := testProvider()
	mp.ApplyReturn = &terraform.ResourceState{
		ID: "foo",
	}

	checkDestroy := false
	checkStep := false

	checkDestroyFn := func(*terraform.State) error {
		checkDestroy = true
		return nil
	}

	checkStepFn := func(s *terraform.State) error {
		checkStep = true

		rs, ok := s.Resources["test_instance.foo"]
		if !ok {
			t.Error("test_instance.foo is not present")
			return nil
		}
		if rs.ID != "foo" {
			t.Errorf("bad check ID: %s", rs.ID)
		}

		return nil
	}

	mt := new(mockT)
	Test(mt, TestCase{
		Providers: map[string]terraform.ResourceProvider{
			"test": mp,
		},
		CheckDestroy: checkDestroyFn,
		Steps: []TestStep{
			TestStep{
				Config: testConfigStr,
				Check:  checkStepFn,
			},
		},
	})

	if mt.failed() {
		t.Fatalf("test failed: %s", mt.failMessage())
	}
	if !checkStep {
		t.Fatal("didn't call check for step")
	}
	if !checkDestroy {
		t.Fatal("didn't call check for destroy")
	}
}

func TestTest_empty(t *testing.T) {
	destroyCalled := false
	checkDestroyFn := func(*terraform.State) error {
		destroyCalled = true
		return nil
	}

	mt := new(mockT)
	Test(mt, TestCase{
		CheckDestroy: checkDestroyFn,
	})

	if mt.failed() {
		t.Fatal("test failed")
	}
	if destroyCalled {
		t.Fatal("should not call check destroy if there is no steps")
	}
}

func TestTest_noEnv(t *testing.T) {
	// Unset the variable
	if err := os.Setenv(TestEnvVar, ""); err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Setenv(TestEnvVar, "1")

	mt := new(mockT)
	Test(mt, TestCase{})

	if !mt.SkipCalled {
		t.Fatal("skip not called")
	}
}

func TestTest_preCheck(t *testing.T) {
	called := false

	mt := new(mockT)
	Test(mt, TestCase{
		PreCheck:     func() { called = true },
	})

	if !called {
		t.Fatal("precheck should be called")
	}
}

func TestTest_stepError(t *testing.T) {
	mp := testProvider()
	mp.ApplyReturn = &terraform.ResourceState{
		ID: "foo",
	}

	checkDestroy := false

	checkDestroyFn := func(*terraform.State) error {
		checkDestroy = true
		return nil
	}

	checkStepFn := func(*terraform.State) error {
		return fmt.Errorf("error")
	}

	mt := new(mockT)
	Test(mt, TestCase{
		Providers: map[string]terraform.ResourceProvider{
			"test": mp,
		},
		CheckDestroy: checkDestroyFn,
		Steps: []TestStep{
			TestStep{
				Config: testConfigStr,
				Check:  checkStepFn,
			},
		},
	})

	if !mt.failed() {
		t.Fatal("test should've failed")
	}
	t.Logf("Fail message: %s", mt.failMessage())

	if !checkDestroy {
		t.Fatal("didn't call check for destroy")
	}
}

// mockT implements TestT for testing
type mockT struct {
	ErrorCalled bool
	ErrorArgs   []interface{}
	FatalCalled bool
	FatalArgs   []interface{}
	SkipCalled  bool
	SkipArgs    []interface{}

	f bool
}

func (t *mockT) Error(args ...interface{}) {
	t.ErrorCalled = true
	t.ErrorArgs = args
	t.f = true
}

func (t *mockT) Fatal(args ...interface{}) {
	t.FatalCalled = true
	t.FatalArgs = args
	t.f = true
}

func (t *mockT) Skip(args ...interface{}) {
	t.SkipCalled = true
	t.SkipArgs = args
	t.f = true
}

func (t *mockT) failed() bool {
	return t.f
}

func (t *mockT) failMessage() string {
	if t.FatalCalled {
		return t.FatalArgs[0].(string)
	} else if t.ErrorCalled {
		return t.ErrorArgs[0].(string)
	} else if t.SkipCalled {
		return t.SkipArgs[0].(string)
	}

	return "unknown"
}

func testProvider() *terraform.MockResourceProvider {
	mp := new(terraform.MockResourceProvider)
	mp.DiffReturn = &terraform.ResourceDiff{
		Attributes: map[string]*terraform.ResourceAttrDiff{
			"foo": &terraform.ResourceAttrDiff{
				New: "bar",
			},
		},
	}
	mp.ResourcesReturn = []terraform.ResourceType{
		terraform.ResourceType{Name: "test_instance"},
	}

	return mp
}

const testConfigStr = `
resource "test_instance" "foo" {}
`