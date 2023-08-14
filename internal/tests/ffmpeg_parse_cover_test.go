package tests

import (
	"os/exec"
	"testing"
)

func TestParseCover(t *testing.T) {
	playUrl := "http://ryr42bm4i.hn-bkt.clouddn.com/video/1691375722.mp4"
	time := "10"
	cmd := exec.Command("ffmpeg", "-i", playUrl, "-ss", time, "-frames:v", "1", "./output.jpg")
	_, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
}
