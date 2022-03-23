package main

import (
	"fmt"

	clamav "github.com/ca110us/go-clamav"
)

func main() {
	// new clamav instance
	c := new(clamav.Clamav)
	err := c.Init(clamav.SCAN_OPTIONS{
		General:   0,
		Parse:     clamav.CL_SCAN_PARSE_ARCHIVE | clamav.CL_SCAN_PARSE_ELF,
		Heuristic: 0,
		Mail:      0,
		Dev:       0,
	})

	if err != nil {
		panic(err)
	}

	// free clamav memory
	defer c.Free()

	// load db
	signo, err := c.LoadDB("./db", uint(clamav.CL_DB_DIRECTORY))
	if err != nil {
		panic(err)
	}
	fmt.Println("db load succeed:", signo)

	// compile engine
	err = c.CompileEngine()
	if err != nil {
		panic(err)
	}

	c.EngineSetNum(clamav.CL_ENGINE_MAX_SCANSIZE, 1024*1024*40)
	c.EngineSetNum(clamav.CL_ENGINE_PCRE_MAX_FILESIZE, 1024*1024*20)
	c.EngineSetNum(clamav.CL_ENGINE_MAX_SCANTIME, 9000)
	c.EngineSetNum(clamav.CL_ENGINE_PCRE_MATCH_LIMIT, 1000)
	c.EngineSetNum(clamav.CL_ENGINE_PCRE_RECMATCH_LIMIT, 500)

	// fmt.Println(c.EngineGetNum(clamav.CL_ENGINE_PCRE_RECMATCH_LIMIT))

	// scan
	scanned, msg, err := c.ScanFile("./test_file/nmap")
	fmt.Println(scanned, msg, err)
}
