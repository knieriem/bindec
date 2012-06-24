package bin

import (
	"fmt"
)

// decode binary values according to bit-, or value-wise descriptions

type Decoder interface {
	Decode(w []string, val int) []string
}

type signal struct {
	pos  uint
	mask int
	name string
}

func Sig(pos uint, name string) Decoder {
	return &signal{pos, 1 << pos, name}
}

func (s *signal) Decode(w []string, val int) (list []string) {
	list = w
	if val&s.mask != 0 {
		if s.name == "<reserved>" {
			list = append(list, fmt.Sprintf("bit %d: %s", s.pos, s.name))
		} else {
			list = append(list, s.name)
		}
	}
	return
}

type value struct {
	pos   uint
	mask  int
	desc  string
	names []string
	dflt  string
}

func Val(startBit, endBit uint, desc string, names []string, dflt string) Decoder {
	mask := ((1 << (endBit + 1)) - 1) - ((1 << startBit) - 1)
	return &value{startBit, mask, desc, names, dflt}
}

func (v *value) Decode(w []string, b int) (list []string) {
	b = b & v.mask >> v.pos

	list = w
	s := ""
	switch {
	case b < len(v.names):
		s = v.names[b]
	case v.dflt != "":
		s = v.dflt
	}
	switch s {
	default:
		list = append(list, v.desc+": "+s)
	case "<reserved>":
		list = append(list, fmt.Sprintf("%s: %d: %s", v.desc, b, s))
	case "":
	}
	return
}

type DecoderList []Decoder

func NewDecoderList(decoders ...[]Decoder) Decoder {
	var list []Decoder
	for _, d := range decoders {
		list = append(list, d...)
	}
	return DecoderList(list)
}

func (list DecoderList) Decode(w []string, val int) []string {
	for _, d := range list {
		w = d.Decode(w, val)
	}
	return w
}

type shift struct {
	pos uint
	d   Decoder
}

func Shift(pos uint, d Decoder) Decoder {
	return &shift{pos, d}
}

func (s shift) Decode(w []string, val int) []string {
	return s.d.Decode(w, val>>s.pos)
}

type group struct {
	name string
	d    Decoder
}

func Group(name string, d Decoder) Decoder {
	return &group{name, d}
}

func (g group) Decode(w []string, val int) []string {
	sub := g.d.Decode(nil, val)
	if sub == nil {
		return w
	} else {
		w = append(w, g.name)
		for _, s := range sub {
			w = append(w, "  "+s)
		}
	}
	return w
}
