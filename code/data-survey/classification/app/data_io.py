import pandas as pd
from app import app

def save_data():
    interesting_snippets.to_csv(app.config['DATA_DIR'] + '/data/interesting_snippets.csv')

interesting_snippets = \
    pd.read_csv(app.config['DATA_DIR'] + '/data/interesting_snippets.csv')

if not 'label' in interesting_snippets:
    interesting_snippets['label'] = 'unclassified'
    save_data()
