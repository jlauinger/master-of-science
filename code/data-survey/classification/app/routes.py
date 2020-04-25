from app import app
from app.data_io import save_data, load_data, get_interesting_files
from app.forms import ClassificationForm

from flask import render_template, flash, redirect
import pandas as pd


current_filename = 'n/a'
interesting_snippets = pd.DataFrame()


@app.route('/')
@app.route('/index')
def index():
    return render_template('index.html', snippets=interesting_snippets.iterrows(), snippets2=interesting_snippets.iterrows(),
                           filename=current_filename)


@app.route('/classify/<int:index>', methods=['GET', 'POST'])
def classify(index):
    form = ClassificationForm()
    next_index = index + 1

    if form.validate_on_submit():
        flash('Classifying as {}'.format(form.label.data))
        interesting_snippets.at[index, 'label'] = form.label.data
        return redirect('/classify/{}'.format(next_index))

    snippet = interesting_snippets.loc[index]
    quick_labels = set(interesting_snippets['label']) | set(['uintptr_type', 'function_call', 'cast', 'protocol'])

    return render_template('classify.html', form=form, snippet=snippet, quick_labels=quick_labels,
                           next_index=next_index, filename=current_filename)


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