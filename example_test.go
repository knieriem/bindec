package bindec_test

import (
	"fmt"

	"github.com/knieriem/bindec"
)

var tempStatReg = bindec.Group("TEMP_STAT", bindec.DecoderList{
	bindec.Sig(0, "TEMP_READY"),
	bindec.Sig(1, "OVERTEMP"),
	bindec.Func(4, 13, "TEMP", func(v int) string {
		return fmt.Sprintf("%.3g °C", float64(v)/1.213-273.15)
	}),
})

func Example_tempStat() {
	printTempStat(0x1a53)
	printTempStat(0x1759)
	printTempStat(0x1758)

	// Output:
	// TEMP_STAT
	//	TEMP_READY
	//	OVERTEMP
	//	TEMP: 73.9 °C
	//
	// TEMP_STAT
	//	TEMP_READY
	//	TEMP: 34.4 °C
	//
	// TEMP_STAT
	//	TEMP: 34.4 °C
}

func printTempStat(regVal uint16) {
	fields := tempStatReg.Decode(nil, int(regVal))
	for _, f := range fields {
		fmt.Println(f)
	}
	fmt.Println()
}
