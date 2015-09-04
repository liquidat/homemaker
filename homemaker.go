/*
 * Copyright (c) 2015 Alex Yatskov <alex@foosoft.net>
 * Author: Alex Yatskov <alex@foosoft.net>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "Software"), to deal in
 * the Software without restriction, including without limitation the rights to
 * use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 * the Software, and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 * FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 * IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
)

const (
	flagClobber = 1 << iota
	flagForce
	flagVerbose
	flagNoCmd
	flagNoLink
	flagNoMacro
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options] conf src\n", path.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "http://foosoft.net/projects/homemaker/\n\n")
	fmt.Fprintf(os.Stderr, "Parameters:\n")
	flag.PrintDefaults()
}

func main() {
	currUsr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	taskName := flag.String("task", "default", "name of task to execute")
	dstDir := flag.String("dest", currUsr.HomeDir, "target directory for tasks")
	force := flag.Bool("force", true, "create parent directories to target")
	clobber := flag.Bool("clobber", false, "delete files and directories at target")
	verbose := flag.Bool("verbose", false, "verbose output")
	nocmd := flag.Bool("nocmd", false, "don't execute commands")
	nolink := flag.Bool("nolink", false, "don't create links")
	variant := flag.String("variant", "", "execution variant")

	flag.Usage = usage
	flag.Parse()

	flags := 0
	if *clobber {
		flags |= flagClobber
	}
	if *force {
		flags |= flagForce
	}
	if *verbose {
		flags |= flagVerbose
	}
	if *nocmd {
		flags |= flagNoCmd
	}
	if *nolink {
		flags |= flagNoLink
	}

	if flag.NArg() == 2 {
		confFile := makeAbsPath(flag.Arg(0))

		conf, err := newConfig(confFile)
		if err != nil {
			log.Fatal(err)
		}

		conf.srcDir = makeAbsPath(flag.Arg(1))
		conf.dstDir = makeAbsPath(*dstDir)
		conf.variant = *variant
		conf.flags = flags

		os.Setenv("HM_CONFIG", confFile)
		os.Setenv("HM_TASK", *taskName)
		os.Setenv("HM_SRC", conf.srcDir)
		os.Setenv("HM_DEST", conf.dstDir)
		os.Setenv("HM_VARIANT", conf.variant)

		if err := processTask(*taskName, conf); err != nil {
			log.Fatal(err)
		}
	} else {
		usage()
		os.Exit(2)
	}
}
