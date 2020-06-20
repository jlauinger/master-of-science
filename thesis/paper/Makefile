all:
	pdflatex paper
	#bibtex paper
	pdflatex paper
	pdflatex paper

open:
	evince paper.pdf &

clean:
	rm -f paper/*.lol paper/*.lot paper/*.fdb_latexmk paper/*.aux paper/*.bbl paper/*.aux paper/*.synctex.gz paper/*.log paper/*.out paper/*.blg paper/*.lof paper/*.toc paper/*.len paper/*.nav paper/*.snm paper/*.vrb paper/*.backup paper/*.tex~ paper/*.bib~

