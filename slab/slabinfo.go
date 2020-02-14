package slab

import (
	"math"
	"strconv"
	"strings"
)

// SlabInfo slabtop -o 命令输出结果
type SlabInfo struct {
	Objs      float64
	Active    float64
	UseObj    float64
	ObjSize   float64
	Slabs     float64
	ObjSlab   float64
	CacheSize float64
	Name      string
}

func New(record []string) *SlabInfo {
	s := &SlabInfo{}
	if len(record) != 8 {
		return s
	}
	s.SetObjs(record[0])
	s.SetActive(record[1])
	s.SetUseObj(record[2])
	s.SetObjSize(record[3])
	s.SetSlabs(record[4])
	s.SetObjSlab(record[5])
	s.SetCacheSize(record[6])
	s.SetName(record[7])
	return s
}

func (s SlabInfo) Equal(b SlabInfo) bool {
	d := 0.001
	if math.Abs(s.Objs-b.Objs) > d {
		return false
	}

	if math.Abs(s.Active-b.Active) > d {
		return false
	}

	if math.Abs(s.UseObj-b.UseObj) > d {
		return false
	}

	if math.Abs(s.ObjSize-b.ObjSize) > d {
		return false
	}

	if math.Abs(s.Slabs-b.Slabs) > d {
		return false
	}

	if math.Abs(s.ObjSlab-b.ObjSlab) > d {
		return false
	}

	if math.Abs(s.CacheSize-b.CacheSize) > d {
		return false
	}

	return s.Name == b.Name
}

func (s *SlabInfo) SetObjs(value string) {
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return
	}
	s.Objs = v
}

func (s *SlabInfo) SetActive(value string) {
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return
	}
	s.Active = v
}

func (s *SlabInfo) SetUseObj(value string) {
	value = strings.Replace(value, "%", "", -1)
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return
	}
	s.UseObj = v
}

func (s *SlabInfo) SetObjSize(value string) {
	value = strings.Replace(value, "K", "", -1)
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return
	}
	s.ObjSize = v
}

func (s *SlabInfo) SetSlabs(value string) {
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return
	}
	s.Slabs = v
}

func (s *SlabInfo) SetObjSlab(value string) {
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return
	}
	s.ObjSlab = v
}

func (s *SlabInfo) SetCacheSize(value string) {
	value = strings.Replace(value, "K", "", -1)
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return
	}
	s.CacheSize = v
}
func (s *SlabInfo) SetName(value string) {
	s.Name = value
}
