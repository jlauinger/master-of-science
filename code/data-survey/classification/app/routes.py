from app import app
from app.data_io import save_data, load_data, get_interesting_files
from app.forms import ClassificationForm1, ClassificationForm2

from flask import render_template, flash, redirect
import pandas as pd

from os import path


# application global state
current_filename = 'n/a'
classifying_label = 'label'
interesting_snippets = pd.DataFrame()


# default page: snippet index
@app.route('/')
@app.route('/index')
def index():
    # render the index / snippet list page by populating its template
    return render_template('index.html', snippets=interesting_snippets.iterrows(), snippets2=interesting_snippets.iterrows(),
                           filename=current_filename, classifying_label=classifying_label)


# snippet details / classification page
@app.route('/classify/<int:index>', methods=['GET', 'POST'])
def classify(index):
    # instantiate classification forms and set the index of the next snippet (used for the Next button)
    form1 = ClassificationForm1()
    form2 = ClassificationForm2()
    next_index = index + 1

    # save label if it was submitted
    if form1.submit1.data and form1.validate():
        flash('Classifying (label 1) as {}'.format(form1.label1.data))
        interesting_snippets.at[index, 'label'] = form1.label1.data
        # then open the next snippet
        return redirect('/classify/{}'.format(next_index))

    # save label2 if it was submitted
    if form2.submit2.data and form2.validate():
            flash('Classifying (label 2) as {}'.format(form2.label2.data))
            interesting_snippets.at[index, 'label2'] = form2.label2.data
            # then open the next snippet
            return redirect('/classify/{}'.format(next_index))

    # get the data to populate the page template
    snippet = interesting_snippets.loc[index]
    quick_labels1 = set(interesting_snippets['label']) | set([])
    quick_labels2 = set(interesting_snippets['label2']) | set([])

    return render_template('classify.html', form1=form1, form2=form2, snippet=snippet, quick_labels1=quick_labels1,
                           quick_labels2=quick_labels2, index=index, next_index=next_index, filename=current_filename)


# code context detail page
@app.route('/file_content/<int:index>', methods=['GET'])
def file_content(index):
    # get the requested snippet
    snippet = interesting_snippets.loc[index]

    # check whether it is contained in the Go standard library
    if snippet.module_path == "std":
        # if so, it is located in the GOROOT directory
        file_path = "{}/src/{}/{}".format(
            app.config['GO_LIB_PATH'],
            snippet.package_import_path,
            snippet.file_name)
    else:
        # otherwise, it is located in the correct module directory within the GOPATH directory
        file_path = "{}/pkg/mod/{}@{}{}/{}".format(
            app.config['GO_MOD_PATH'],
            snippet.module_path,
            snippet.module_version,
            snippet.package_import_path[len(snippet.module_path):],
            snippet.file_name)

    # redirect to the snippet if the file cannot be found and show an error message
    if not path.exists(file_path):
        flash("Path {} not found".format(file_path))
        return redirect("/classify/{}".format(index))

    # read the complete file
    with open(file_path, "r") as f:
        content = f.readlines()

    # add line numbers to the beginning of each line
    content = ["{}: {}".format(str(i+1).rjust(7, " "), line) for i, line in enumerate(content)]

    return render_template('file_content.html', content=content, file_path=file_path)


# action to save the current global state
@app.route('/save')
def save():
    save_data(current_filename, interesting_snippets)
    # there is no actual page to display here, redirect to the list of snippets
    return redirect('/index')


# file selection page
@app.route('/switch-files')
def switch_files_index():
    files = get_interesting_files()
    return render_template('switch_files.html', files=enumerate(files))


# action to select and load a file
@app.route('/switch-files/<int:idx>')
def switch_files_action(idx):
    global current_filename, interesting_snippets

    # the file is referenced by index, therefore I need the list of available files again
    files = get_interesting_files()
    current_filename = files[idx]

    interesting_snippets = load_data(current_filename)

    # there is no actual page to display here, redirect to the list of snippets
    return redirect('/index')


# action to switch the current labeling dimension
@app.route('/switch-label')
def switch_label():
    global classifying_label

    # reverse the label
    if classifying_label == 'label':
        classifying_label = 'label2'
    else:
        classifying_label = 'label'

    # there is no actual page to display here, redirect to the list of snippets
    return redirect('/index')
