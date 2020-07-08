package vss_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/forensicanalysis/fslib/filesystem/zip"
	"github.com/forensicanalysis/go-vss"
)

func TestReader(t *testing.T) {
	reader, err := os.Open("vss.raw.zip")
	if err != nil {
		t.Fatal("Cannot open file descriptor")
	}
	defer reader.Close()

	zipfs, err := zip.New(reader)
	if err != nil {
		t.Fatal(err)
	}
	// defer zipfs.Close()

	item, err := zipfs.Open("/vss.raw")
	if err != nil {
		t.Fatal(err)
	}
	defer item.Close()

	store, err := vss.NewByReader(item)
	if err != nil {
		t.Fatal(err)
	}

	root, err := store[0].Open("/")
	if err != nil {
		t.Fatal(err)
	}
	filenames, err := root.Readdirnames(0)
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{"$AttrDef", "$BadClus", "$BadClus:$Bad", "$Bitmap", "$Boot", "$Extend", "$LogFile", "$MFT", "$MFTMirr", "$Secure", "$Secure:$SDS", "$UpCase", "$Volume", "System Volume Information", "syslog"}
	if !reflect.DeepEqual(filenames, expected) {
		t.Fatal(filenames)
	}
}
