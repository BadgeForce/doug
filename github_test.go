package doug

import "testing"
import "os"

func TestMakeTempDir(t *testing.T) {
	name, err := makeTempDir()
	if err != nil {
		t.Errorf("Expected err to be nil but got %s", err.Error())
	}
	os.RemoveAll("./" + name)
}

func TestRemoveTempDir(t *testing.T) {
	name := "testing-directory-1234"
	err := os.Mkdir(name, 0777)
	if err != nil {
		t.Errorf("Expected err to be nil but got %s", err.Error())
	}

	os.RemoveAll("./" + name)
}
