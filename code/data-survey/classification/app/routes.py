from app import app
from app.data_io import save_data, load_data, get_interesting_files
from app.forms import ClassificationForm1, ClassificationForm2

from flask import render_template, flash, redirect
import pandas as pd

from os import path


current_filename = 'n/a'
classifying_label = 'label'
interesting_snippets = pd.DataFrame()


@app.route('/')
@app.route('/index')
def index():
    return render_template('index.html', snippets=interesting_snippets.iterrows(), snippets2=interesting_snippets.iterrows(),
                           filename=current_filename, classifying_label=classifying_label)


@app.route('/classify/<int:index>', methods=['GET', 'POST'])
def classify(index):
    form1 = ClassificationForm1()
    form2 = ClassificationForm2()
    next_index = index + 1

    if form1.submit1.data and form1.validate():
        flash('Classifying (label 1) as {}'.format(form1.label1.data))
        interesting_snippets.at[index, 'label'] = form1.label1.data
        return redirect('/classify/{}'.format(next_index))

    if form2.submit2.data and form2.validate():
            flash('Classifying (label 2) as {}'.format(form2.label2.data))
            interesting_snippets.at[index, 'label2'] = form2.label2.data
            return redirect('/classify/{}'.format(next_index))

    snippet = interesting_snippets.loc[index]
    quick_labels1 = set(interesting_snippets['label']) | set([])
    quick_labels2 = set(interesting_snippets['label2']) | set([])

    return render_template('classify.html', form1=form1, form2=form2, snippet=snippet, quick_labels1=quick_labels1,
                           quick_labels2=quick_labels2, index=index, next_index=next_index, filename=current_filename)


@app.route('/file_content/<int:index>', methods=['GET'])
def file_content(index):
    snippet = interesting_snippets.loc[index]

    if snippet.module_path == "std":
        file_path = "{}/src/{}/{}".format(
            app.config['GO_LIB_PATH'],
            snippet.package_import_path,
            snippet.file_name)
    else:
        file_path = "{}/pkg/mod/{}@{}{}/{}".format(
            app.config['GO_MOD_PATH'],
            snippet.module_path,
            snippet.module_version,
            snippet.package_import_path[len(snippet.module_path):],
            snippet.file_name)

    if not path.exists(file_path):
        flash("Path {} not found".format(file_path))
        return redirect("/classify/{}".format(index))

    with open(file_path, "r") as f:
        content = f.readlines()

    content = ["{}: {}".format(str(i+1).rjust(7, " "), line) for i, line in enumerate(content)]

    return render_template('file_content.html', content=content, file_path=file_path)


@app.route('/save')
def save():
    save_data(current_filename, interesting_snippets)

    return redirect('/index')


@app.route('/switch-files')
def switch_files_index():
    files = get_interesting_files()

    return render_template('switch_files.html', files=enumerate(files))


@app.route('/switch-files/<int:idx>')
def switch_files_action(idx):
    global current_filename, interesting_snippets

    files = get_interesting_files()
    current_filename = files[idx]

    interesting_snippets = load_data(current_filename)

    return redirect('/index')


@app.route('/switch-label')
def switch_label():
    global classifying_label

    if classifying_label == 'label':
        classifying_label = 'label2'
    else:
        classifying_label = 'label'

    return redirect('/index')