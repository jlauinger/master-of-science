#!/usr/bin/env python3

import pandas as pd
import subprocess

projects = pd.read_csv('/root/data/projects.csv')

for i, project in projects.iterrows():
    command = ['git', 'submodule', 'add', project['project_github_clone_url'], 'rawdata/{}'.format(project['project_name'])]
    subprocess.run(command)

