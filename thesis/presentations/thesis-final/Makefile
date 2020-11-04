all:
	lualatex slides
	bibtex slides
	lualatex slides
	lualatex slides

open:
	evince slides.pdf &

release:
	@cp slides.pdf slides-lauinger-$$(date +%Y%m%d)-$$(git rev-parse --short=7 HEAD).pdf
	@echo Releasing slides-lauinger-$$(date +%Y%m%d)-$$(git rev-parse --short=7 HEAD).pdf

clean:
	rm -f *.bcf *.bib.bak *.acn *.acr *.alg *.glg *.glo *.gls *.ist *.lol *.lot *.fdb_latexmk *.aux *.bbl *.aux *.synctex.gz *.log *.out *.blg *.lof *.toc *.len *.nav *.snm *.vrb *.backup *.tex~ *.bib~ *.glsdefs *.run.xml *-blx.bib *.auxlock *.xmpi *.xmpdata *.mw
	rm -f gfx/*.lol gfx/*.lot gfx/*.fdb_latexmk gfx/*.aux gfx/*.bbl gfx/*.aux gfx/*.synctex.gz gfx/*.log gfx/*.out gfx/*.blg gfx/*.lof gfx/*.toc gfx/*.len gfx/*.nav gfx/*.snm gfx/*.vrb gfx/*.backup gfx/*.tex~ gfx/*.bib~
