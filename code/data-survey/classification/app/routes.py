from app import app
from app.data_io import interesting_snippets, save_data
from app.forms import ClassificationForm

from flask import render_template, flash, redirect

@app.route('/')
@app.route('/index')
def index():
    return render_template('index.html', snippets=interesting_snippets.iterrows(), snippets2=interesting_snippets.iterrows())

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

    return render_template('classify.html', form=form, snippet=snippet, quick_labels=quick_labels, next_index=next_index)

@app.route('/save')
def save():
    save_data()
    return redirect('/index')
