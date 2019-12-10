# Makefile
# Finley McIlwaine
# Dec. 8, 2019
# COSC 4785 Program 5

OUT=-o invert
FILES=main.go index.go
all: main

.PHONY: clean tarball

main: main.go
	go build ${OUT} ${FILES} 

clean:
	/bin/rm -rf invert indexes.log ./tarball

tarball:
	rm -rf tarball
	mkdir tarball
	tar -cvf ./tarball/invert.tar makefile ${FILES} ./test-docs

