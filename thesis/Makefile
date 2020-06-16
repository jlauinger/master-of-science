all:
	pdflatex thesis
	bibtex thesis
	makeglossaries thesis
	pdflatex thesis
	pdflatex thesis

open:
	evince thesis.pdf &

clean:
	rm -f *.bcf *.bib.bak *.acn *.acr *.alg *.glg *.glo *.gls *.ist *.lol *.lot *.fdb_latexmk *.aux *.bbl *.aux *.synctex.gz *.log *.out *.blg *.lof *.toc *.len *.nav *.snm *.vrb *.backup *.tex~ *.bib~ *.glsdefs *.run.xml *-blx.bib *.auxlock *.xmpi
	rm -f frontback/*.lol frontback/*.lot frontback/*.fdb_latexmk frontback/*.aux frontback/*.bbl frontback/*.aux frontback/*.synctex.gz frontback/*.log frontback/*.out frontback/*.blg frontback/*.lof frontback/*.toc frontback/*.len frontback/*.nav frontback/*.snm frontback/*.vrb frontback/*.backup frontback/*.tex~ frontback/*.bib~
	rm -f chapters/*.lol chapters/*.lot chapters/*.fdb_latexmk chapters/*.aux chapters/*.bbl chapters/*.aux chapters/*.synctex.gz chapters/*.log chapters/*.out chapters/*.blg chapters/*.lof chapters/*.toc chapters/*.len chapters/*.nav chapters/*.snm chapters/*.vrb chapters/*.backup chapters/*.tex~ chapters/*.bib~
	rm -f config/*.lol config/*.lot config/*.fdb_latexmk config/*.aux config/*.bbl config/*.aux config/*.synctex.gz config/*.log config/*.out config/*.blg config/*.lof config/*.toc config/*.len config/*.nav config/*.snm config/*.vrb config/*.backup config/*.tex~ config/*.bib~
