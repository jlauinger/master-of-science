import pandas as pd
import glob

from app import app


def save_data(filename, interesting_snippets):
    interesting_snippets.to_csv(filename)


def load_data(filename):
    interesting_snippets = pd.read_csv(filename)

    if 'label' not in interesting_snippets:
        interesting_snippets['label'] = 'unclassified'
        save_data(filename, interesting_snippets)

    if 'label2' not in interesting_snippets:
        interesting_snippets['label2'] = 'unclassified'
        save_data(filename, interesting_snippets)

    return interesting_snippets


def get_interesting_files():
    return sorted(list(glob.glob(app.config['DATA_DIR'] + '/classification/*.csv')))
