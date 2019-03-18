package duitdraw

import "testing"

func TestStringWidth(t *testing.T) {
	tt := []string{
		"Hello world!",
		"I can eat glass and it doesn't hurt me.",
		"私はガラスを食べられます。それは私を傷つけません。",
		"আমি কাঁচ খেতে পারি, তাতে আমার কোনো ক্ষতি হয় না।",
	}
	for _, tc := range tt {
		sum := 0
		for _, c := range tc {
			sum += defaultFont.StringWidth(string(c))
		}
		dx := defaultFont.StringWidth(tc)
		if dx != sum {
			t.Errorf("StringWidth(%q) is %v; expected %v", tc, dx, sum)
		}
	}
}
