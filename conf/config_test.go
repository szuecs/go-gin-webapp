package conf

import "testing"

func TestNew(t *testing.T) {
	cfg, err := New()
	if err != nil {
		t.Fatalf("ERR: conf.TestNew failed, caused by: %s", err)
	}
	if cfg == nil {
		t.Fatal("ERR: conf.TestNew returned and invalid config")
	}
}
