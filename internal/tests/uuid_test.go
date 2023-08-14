package tests

import (
	"simple_tiktok_rime/pkg/toolx"
	"testing"
)

func TestUUID(t *testing.T) {
	t.Log("uuid => ", toolx.GenUUID())
}
