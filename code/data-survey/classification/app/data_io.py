import pandas as pd
import glob

from app import app


def save_data(filename, interesting_snippets):
    interesting_snippets.to_csv(filename)


def load_data(filename):
    interesting_snippets = pd.read_csv(filename)

    # assign 'unclassified' as default for label
    if 'label' not in interesting_snippets:
        interesting_snippets['label'] = 'unclassified'
        save_data(filename, interesting_snippets)

    # assign 'unclassified' as default for label2
    if 'label2' not in interesting_snippets:
        interesting_snippets['label2'] = 'unclassified'
        save_data(filename, interesting_snippets)

    return interesting_snippets


def get_interesting_files():
    # available files are all CSV files in the classification directory
    return sorted(list(glob.glob(app.config['DATA_DIR'] + '/classification/*.csv')))
