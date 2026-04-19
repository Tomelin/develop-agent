package agent

import "testing"

func TestChannelRegistryDefaultBuffer(t *testing.T) {
	r := NewChannelRegistry(0)
	rt := r.Create("a1")
	if cap(rt.In) != 10 || cap(rt.Out) != 10 {
		t.Fatalf("expected default buffer 10, got in=%d out=%d", cap(rt.In), cap(rt.Out))
	}

	if rt.CurrentStatus() != StatusIdle {
		t.Fatalf("expected status IDLE, got %s", rt.CurrentStatus())
	}

	r.Destroy("a1")
	if _, ok := r.Get("a1"); ok {
		t.Fatal("runtime should be removed")
	}
}
