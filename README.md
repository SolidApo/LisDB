*Lis is still WIP, so this README might contain inaccurate information.*

# LisDB

**LisDB** is a database that stores relations between *nodes* using *tags*.

## Installation / Building

You can find latest binaries & installers for Linux at [https://dl.solidapo.de/LisDB](https://dl.solidapo.de/LisDB/).

If you prefer compiling yourself, you will need [the Go programming language](https://go.dev/).
Simply run `go install codeberg.org/SolidApo/lisdb@latest` in your command line, 
or clone this repo and run `make && sudo make install`, and you're set!
The Makefile also contains targets to build various packages.

Please mind that LisDB requires a Unix-like system (e.g. Linux).

## What data model does LisDB follow?

To be honest, I don't know. I would describe it as being somewhere inbetween of triplestores, relational databases and graph databases.
