package govid_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/tomdoherty/govid"
)

func TestFilter(t *testing.T) {
	testCases := []struct {
		file, want string
	}{
		{
			file: "covid.json",
			want: "2020-03-21 00:00:00 +0000 UTC\t1\t0\tZW14645473\n",
		},
	}

	var w govid.LogWriter
	o := bytes.Buffer{}

	if err := w.Init(&o); err != nil {
		t.Fatal(err)
	}

	for _, tc := range testCases {
		f, err := os.Open(tc.file)

		if err != nil {
			t.Fatal(err)
		}

		govid.Filter(f, &w)
		if o.String() != tc.want {
			t.Errorf("want %q, got %q", tc.want, o.String())
		}
	}
}
