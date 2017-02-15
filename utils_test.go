package glog

import (
	"fmt"
	"testing"
)

type Movie struct {
	Title, Subtitle string
	Year            int
	Color           bool
	Actor           map[string]string
	Oscars          []string
	Sequel          *string
}

func TestDisplay(t *testing.T) {
	//strangelove := Movie{
	//	Title:    "Dr. Strangelove",
	//	Subtitle: "How I Learned to Stop Worrying and Love the Bomb",
	//	Year:     1964,
	//	Color:    false,
	//	Actor: map[string]string{
	//		"Dr. Strangelove":            "Peter Sellers",
	//		"Grp. Capt. Lionel Mandrake": "Peter Sellers",
	//		"Pres. Merkin Muffley":       "Peter Sellers",
	//		"Gen. Buck Turgidson":        "George C. Scott",
	//		"Brig. Gen. 3ack D. Ripper":  "Sterling Hayden",
	//		`Maj.T.3."King" Kong`: "Slim Pickens",
	//	},
	//	Oscars: []string{
	//		"Best Actor (Nomin.)",
	//		"Best Adapted Screenplay (Nomin.)",
	//		"Best Director (Nomin.)",
	//		"Best Picture (Nomin.)",
	//	},
	//}
	//fmt.Println(strangelove)
	//fmt.Println("==========================================")
	//fmt.Println("strangelove", Display(strangelove))                  /
	arr := getStackTrace()
	fmt.Println("getStackTrace", Display(arr))
}

/*
========================================
	BENCHMARKS!!
========================================
	# go test -bench=.
*/

func BenchmarkGetStack1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = getStackTrace()
	}
}

func BenchmarkGetStack2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = CallerInfo2()
	}
}
