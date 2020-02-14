package slab

import (
	"encoding/csv"
	"io"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseData(t *testing.T) {
	data := `5063448,4803097,94%,0.10K,129832,39,519328K,buffer_head
3305064,2802929,84%,0.19K,78692,42,629536K,dentry
31625,31185,98%,0.31K,1265,25,10120K,nf_conntrack_ffffffffbb3129c0`

	type Test struct {
		Input []string
		T     *SlabInfo
		Want  bool
	}

	tests := make([]Test, 0)

	r := csv.NewReader(strings.NewReader(data))

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		test := Test{
			Input: record,
			Want:  true,
		}
		tests = append(tests, test)
	}

	info0 := SlabInfo{
		Objs:      5063448,
		Active:    4803097,
		UseObj:    94,
		ObjSize:   0.10,
		Slabs:     129832,
		ObjSlab:   39,
		CacheSize: 519328,
		Name:      "buffer_head",
	}

	info1 := SlabInfo{
		Objs:      3305064,
		Active:    2802929,
		UseObj:    84,
		ObjSize:   0.19,
		Slabs:     78692,
		ObjSlab:   42,
		CacheSize: 629536,
		Name:      "dentry",
	}

	info2 := SlabInfo{
		Objs:      31625,
		Active:    31185,
		UseObj:    98,
		ObjSize:   0.31,
		Slabs:     1265,
		ObjSlab:   25,
		CacheSize: 10120,
		Name:      "nf_conntrack_ffffffffbb3129c0",
	}

	test3 := Test{
		Input: []string{},
		T:     &info2,
		Want:  false,
	}

	tests = append(tests, test3)

	tests[0].T = &info0
	tests[1].T = &info1
	tests[2].T = &info2

	for _, test := range tests {
		infot := New(test.Input)
		assert.Equal(t, test.Want, infot.Equal(*test.T))
	}

}
