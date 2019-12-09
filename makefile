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
	/bin/rm -rf invert

tarball:
	rm -rf tarball
	mkdir tarball
	tar -cf ./tarball/program5.tar makefile ${FILES}

