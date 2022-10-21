package bindec

import (
	"fmt"
	"strings"
)

// decode binary values according to bit-, or value-wise descriptions

type Decoder interface {
	Decode(w []string, val int) []string
}

type signal struct {
	pos    uint
	mask   int
	name   string
	isFlag bool
	negate bool
}

func Sig(pos uint, name string) Decoder {
	return &signal{pos: pos, mask: 1 << pos, name: name}
}

func Flag(pos uint, name string) Decoder {
	negate := false
	if strings.HasPrefix(name, "!") {
		negate = true
		name = name[1:]
	}
	return &signal{pos: pos, mask: 1 << pos, name: name, isFlag: true, negate: negate}
}

func (s *signal) Decode(w []string, val int) (list []string) {
	var str string
	list = w

	v := val&s.mask != 0
	if s.negate {
		v = !v
	}
	switch {
	case s.isFlag:
		if v {
			str = s.name
		} else {
			str = "!" + s.name
		}
	case val&s.mask == 0:
		return
	default:
		if s.name == "<reserved>" {
			str = fmt.Sprintf("bit %d: %s", s.pos, s.name)
		} else {
			str = s.name
		}
	}
	list = append(list, str)
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
	desc := v.desc
	if desc != "" {
		desc += ": "
	}
	switch s {
	default:
		list = append(list, desc+s)
	case "<reserved>":
		list = append(list, fmt.Sprintf("%s%d: %s", desc, b, s))
	case "":
	}
	return
}

type intval struct {
	pos    uint
	mask   int
	desc   string
	format string
	f      func(int) string
}

func Int(startBit, endBit uint, desc string, format string) Decoder {
	mask := ((1 << (endBit + 1)) - 1) - ((1 << startBit) - 1)
	return &intval{startBit, mask, desc, format, nil}
}

func Func(startBit, endBit uint, desc string, f func(int) string) Decoder {
	mask := ((1 << (endBit + 1)) - 1) - ((1 << startBit) - 1)
	return &intval{startBit, mask, desc, "", f}
}

func (v *intval) Decode(w []string, b int) (list []string) {
	var s string

	b = b & v.mask >> v.pos

	if v.f == nil {
		s = fmt.Sprintf(v.format, b)
	} else {
		s = v.f(b)
	}
	list = w
	if v.desc == "" {
		return
	}
	list = append(list, v.desc+": "+s)
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
			w = append(w, "\t"+s)
		}
	}
	return w
}
