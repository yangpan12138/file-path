package main

import "testing"

func TestFilePath(t *testing.T) {
	f := newFolder(`D:\yp\my_work_dir\File-path\`)

	if err := f.Get(); err != nil {
		t.Fatal(err)
	}
}
