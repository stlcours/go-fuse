package test

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/pathfs"
)

var _ = log.Println

type DefaultReadFS struct {
	pathfs.DefaultFileSystem
	size  uint64
	exist bool
}

func (fs *DefaultReadFS) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	if name == "" {
		return &fuse.Attr{Mode: fuse.S_IFDIR | 0755}, fuse.OK
	}
	if name == "file" {
		return &fuse.Attr{Mode: fuse.S_IFREG | 0644, Size: fs.size}, fuse.OK
	}
	return nil, fuse.ENOENT
}

func (fs *DefaultReadFS) Open(name string, f uint32, context *fuse.Context) (fuse.File, fuse.Status) {
	return &fuse.DefaultFile{}, fuse.OK
}

func defaultReadTest(t *testing.T) (root string, cleanup func()) {
	fs := &DefaultReadFS{}
	var err error
	dir, err := ioutil.TempDir("", "go-fuse")
	if err != nil {
		t.Fatalf("TempDir failed: %v", err)
	}
	pathfs := pathfs.NewPathNodeFs(fs, nil)
	state, _, err := fuse.MountNodeFileSystem(dir, pathfs, nil)
	if err != nil {
		t.Fatalf("MountNodeFileSystem failed: %v", err)
	}
	state.Debug = fuse.VerboseTest()
	go state.Loop()

	return dir, func() {
		state.Unmount()
		os.Remove(dir)
	}
}

func TestDefaultRead(t *testing.T) {
	root, clean := defaultReadTest(t)
	defer clean()

	_, err := ioutil.ReadFile(root + "/file")
	if err == nil {
		t.Fatal("should have failed read.")
	}
}