package main

import (
	"log"
	"os"
	"runtime/pprof"

	"github.com/pboyd04/godata/filter/parser"
)

func testToken(input string) {
	parser, err := parser.NewParser(input)
	if err != nil {
		log.Println(err)
	}
	_, err = parser.GetOperation()
	if err != nil {
		log.Printf("%s => %#v\n", input, err)
	}
}

//nolint:funlen
func main() {
	f, err := os.Create("cpu.pprof")
	if err != nil {
		log.Panic(err)
	}
	err = pprof.StartCPUProfile(f)
	if err != nil {
		log.Panic(err)
	}
	defer pprof.StopCPUProfile()
	for i := 0; i < 100000; i++ {
		testToken("true")
		testToken("false")
		testToken("Name eq 'Milk'")
		testToken("(Name eq 'Milk')")
		testToken("Name ne 'Milk'")
		testToken("Name gt 'Milk'")
		testToken("Name ge 'Milk'")
		testToken("Name lt 'Milk'")
		testToken("Name le 'Milk'")
		testToken("Name eq 'Milk' and Price lt 2.55")
		testToken("Name EQ 'Milk' AND Price LT 2.55")
		testToken("Name eq 'Milk' AND Price LT 2.55")
		testToken("Name eq 'Milk' or Price lt 2.55")
		testToken("Name in ('Milk', 'Cheese')")
		testToken("Name in ['Milk', 'Cheese']")
		testToken("contains(Name,'red')")
		testToken(`Address eq {"Street":"NE 40th","City":"Redmond","State":"WA","ZipCode":"98052"}`)
		testToken("endswith(Name,'ilk')")
		testToken("not endswith(Name,'ilk')")
		testToken("length(CompanyName) eq 19")
		testToken("startswith(CompanyName,'Futterkiste')")
		testToken(`hassubset(Names,["Milk", "Cheese"])`)
		testToken(`Price add 2.45 eq 5.00`)
		testToken(`Price sub 0.55 eq 2.00`)
		testToken(`Price mul 2.0 eq 5.10`)
		testToken(`Price div 2.55 eq 1`)
		testToken(`Price div 2 eq 2`)
		testToken(`Price divby 2 eq 2.5`)
		testToken(`Rating mod 5 eq 0`)
		testToken(`style has Sales.Pattern'Yellow'`)
		testToken(`(4 add 5) mod (4 sub 1) eq 0`)
		testToken(`concat(concat(City,', '),Country) eq 'Berlin, Germany'`)
		testToken(`indexof(CompanyName,'lfreds') eq 1`)
		testToken(`substring(CompanyName,1) eq 'lfreds Futterkiste'`)
		testToken(`substring(CompanyName,1,2) eq 'lf'`)
		testToken(`hassubsequence([4,1,3],[4,1,3])`)
		testToken(`hassubsequence([4,1,3],[4,1])`)
		testToken(`hassubsequence([4,1,3],[4,3])`)
		testToken(`hassubsequence([4,1,3,1],[1,1])`)
		testToken(`hassubsequence([4,1,3,1],[1,1])`)
		testToken(`hassubsequence([4,1,3],[1,3,4])`)
		testToken(`hassubsequence([4,1,3],[3,1])`)
		testToken(`hassubsequence([1,2],[1,1,2])`)
		testToken(`matchesPattern(CompanyName,'%5EA.*e$')`)
		testToken(`tolower(CompanyName) eq 'alfreds futterkiste'`)
		testToken(`toupper(CompanyName) eq 'ALFREDS FUTTERKISTE'`)
		testToken(`trim(CompanyName) eq CompanyName`)
		testToken(`day(BirthDate) eq 8`)
		testToken(`fractionalseconds(BirthDate) lt 0.1`)
		testToken(`hour(BirthDate) eq 4`)
		testToken(`minute(BirthDate) eq 40`)
		testToken(`month(BirthDate) eq 5`)
		testToken(`second(BirthDate) eq 40`)
		testToken(`year(BirthDate) eq 1971`)
		testToken(`ceiling(Freight) eq 32`)
		testToken(`floor(Freight) eq 32`)
		testToken(`round(Freight) eq 32`)
		testToken(`DiscontinuedDate eq null`)
	}
}
