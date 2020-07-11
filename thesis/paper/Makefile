all:
	pdflatex main
	bibtex main
	pdflatex main
	pdflatex main

open:
	evince main.pdf &

clean:
	rm -f *.lol *.lot *.fdb_latexmk *.aux *.bbl *.aux *.synctex.gz *.log *.out *.blg *.lof *.toc *.len *.nav *.snm *.vrb *.backup *.tex~ *.bib~

