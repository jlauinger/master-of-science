# How To Use

In this version of the template, all class files are included.
However, fonts still have to be installed, see the first link in the next section to acquire the fonts.

The main file is `thesis.txt` which should be modified by the student to update the structure, metadata, and includes.
All content files should be in some sub folder, the other files in the main folder must not be modified.

# Important Links

The TUDa Corporate Design and fonts: https://www.intern.tu-darmstadt.de/arbeitsmittel/corporate_design_vorlagen/index.de.jsp

The github page of the new TUDa Latex classes: https://github.com/tudace/tuda_latex_templates

The new TUDa Latex classes: https://www.ce.tu-darmstadt.de/ce/latex_tuda/index.de.jsp

(The old TUDa Latex classes: http://exp1.fkp.physik.tu-darmstadt.de/tuddesign/)

Information about theses can be found on the [webpage of the FB 20](https://www.informatik.tu-darmstadt.de/studium_fb20/im_studium/studienbuero/abschlussarbeiten_fb20/index.de.jsp).

Information on how to submit your final thesis in an electronic form can be found on the [FAQ of the university](https://www.tu-darmstadt.de/studieren/tucan_studienorganisation/tucan_faq/details_96256.de.jsp).

Information about plagiarism and scientific ethic are available on the [webpage of the FB 20](https://www.informatik.tu-darmstadt.de/studium_fb20/im_studium/studienbuero/plagiarismus/index.de.jsp). 

# Other Material to Write a Good Thesis

*  “How to write a successful Bachelor’s/Master’s thesis” by Elmar Jürgens from TUM <https://thesisguide.org/>
*  ["Writing academic papers"](https://sarahnadi.org/writing-papers/) by Sarah Nadi from [University of Alberta](https://sarahnadi.org/smr/), former post-doc in our group. 

# For Supervisors: Thesis Template

**Web interface**:
There is a button on github, on the top of the page called "use this template", that creates a copy of this repository.

**Command line**:
To create a new thesis repoistory from the template, create a new repository "thesis-YYYY-StudentLastName" and run:

    git clone --bare git@github.com:stg-tud/thesis-template.git
    cd thesis-template.git
    git push --mirror git@github.com:stg-tud/thesis-YYYY-StudentLastName.git
    cd ..
    rm -rf thesis-template.git

**Then**:
From the GitHub webpage, add “[supervisor @yourgithubname]” to the repository description, to make it easier to identify the one responsible for the repository.
