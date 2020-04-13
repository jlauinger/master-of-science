#!/usr/bin/env python
import json
import os
from itertools import islice
from os import path
import subprocess

import pandas as pd

PROJECT_DATA_FILE = "../data/projects.csv"
MODULES_DATA_FILE = "../data/modules.csv"
MATCHES_DATA_FILE = "../data/unsafe_matches.csv"
VET_RESULTS_DATA_FILE = "../data/vet_results.csv"

MATCH_TYPES = (
    'unsafe.Pointer', 'unsafe.Sizeof', 'unsafe.Alignof', 'unsafe.Offsetof',
    'uintptr', 'reflect.SliceHeader', 'reflect.StringHeader'
)


def extract_registry(import_path):
    path_components = import_path.split('/')
    registry_components = path_components[0:2] if path_components[1] == 'x' else path_components[0:1]
    return "/".join(registry_components)


def get_go_vet_output(file, unsafe_ptr=False):
    return os.popen("bash -c 'cd \"" + path.dirname(file) + "\" && go vet " +
                    ("-unsafeptr " if unsafe_ptr else "") +
                    "\"" + path.basename(file) + "\" 2>&1 | grep -v \"#\"'").read()


def get_grep_output(module_path, file_name, match_type):
    return os.popen("bash -c 'cd \"" + module_path + "\" && rg " + match_type + " --context 5 --json \"" +
                    file_name + "\"'").read()


projects_df = pd.read_csv(PROJECT_DATA_FILE)

modules_df = pd.DataFrame(columns=['project_name', 'module_import_path', 'module_registry', 'module_version',
                                   'module_number_go_files'])
matches_df = pd.DataFrame(columns=['project_name', 'module_import_path', 'module_registry', 'module_version', 'module_number_go_files',
                                   'file_name', 'file_size_bytes', 'file_size_lines', 'file_imports_unsage_pkg',
                                   'file_go_vet_output', 'text', 'context', 'line_number', 'byte_offset', 'match_type'])
vet_df = pd.DataFrame(columns=['project_name', 'module_import_path', 'module_registry', 'module_version', 'module_number_go_files',
                               'file_name', 'file_size_bytes', 'file_size_lines', 'file_imports_unsage_pkg',
                               'file_go_vet_output', 'line_number', 'message'])

for i, project in projects_df.iterrows():
    try:
        if not path.exists(project['project_checkout_path'] + "/go.mod"):
            print(str(i+1) + "/" + str(len(projects_df)) + ": Skipping " + project['project_name'])
            continue

        print(str(i+1) + "/" + str(len(projects_df)) + ": Analyzing " + project['project_name'])

        modules = [module for module
                   in os.popen("bash -c 'cd \"" + project['project_checkout_path'] +
                               "\" && go mod vendor -v 2>&1 | grep -v \"#\" | sort | uniq'").read().split("\n")
                   if len(module) > 0 and "warning" not in module]

        for module in modules:
            try:
                module_path = project['project_checkout_path'] + "/vendor/" + module

                files = [file for file
                         in os.popen("find '" + module_path + "' -name '*.go'") \
                             .read().split("\n")
                         if len(file) > 0]

                module_data = {
                    'module_import_path': module,
                    'module_registry': extract_registry(module),
                    'module_version': 'n/a',
                    'module_number_go_files': len(files),
                }

                modules_df = modules_df.append({'project_name': project['project_name'], **module_data}, ignore_index=True)

                print("  " + module + " (" + str(module_data['module_number_go_files']) + " files)...")

                for file in files:
                    try:
                        go_vet_output = get_go_vet_output(file)

                        file_data = {
                            'file_name': file[len(module_path)+1:],
                            'file_size_bytes': int(os.popen("wc -c '" + file + "'").read().split(" ")[0]),
                            'file_size_lines': int(os.popen("wc -l '" + file + "'").read().split(" ")[0]),
                            'file_imports_unsage_pkg': subprocess.call(["grep", "unsafe", file], stdout=subprocess.DEVNULL) == 0,
                            'file_go_vet_output': go_vet_output,
                        }
                    except:
                        print("IGNORING ERROR!")
                        continue

                    try:
                        vet_findings = [finding for finding in go_vet_output.split("\n") if len(finding) > 0]

                        for vet_finding in vet_findings:
                            try:
                                vet_finding_data = {
                                    'line_number': vet_finding.split(":")[2],
                                    'message': ":".join(vet_finding.split(":")[4:]).strip(),
                                }

                                vet_df = vet_df.append({'project_name': project['project_name'],
                                                        **module_data, **file_data, **vet_finding_data}, ignore_index=True)
                            except:
                                print("IGNORING ERROR!")
                                continue
                    except:
                        print("IGNORING ERROR!")

                    for match_type in MATCH_TYPES:
                        try:
                            grep_messages = [json.loads(message) for message
                                             in get_grep_output(module_path, file_data['file_name'], match_type).split("\n")
                                             if len(message) > 0]
                        except:
                            print("IGNORING ERROR!")
                            continue

                        for j, line in enumerate(grep_messages):
                            try:
                                if line['type'] == 'match':
                                    context_lines = grep_messages[max(0, j-5) : min(len(grep_messages), j+1+5)]
                                    context = "".join([cl['data']['lines']['text'] for cl in context_lines
                                                       if cl['type'] == 'context' or cl['type'] == 'match'])

                                    match_data = {
                                        'text': line['data']['lines']['text'],
                                        'context': context,
                                        'line_number': line['data']['line_number'],
                                        'byte_offset': line['data']['absolute_offset'],
                                        'match_type': match_type,
                                    }

                                    matches_df = matches_df.append({'project_name': project['project_name'],
                                                                    **module_data, **file_data, **match_data}, ignore_index=True)
                            except:
                                print("IGNORING ERROR!")
                                continue

                try:
                    modules_df.to_csv(MODULES_DATA_FILE, index=False)
                    matches_df.to_csv(MATCHES_DATA_FILE, index=False)
                    vet_df.to_csv(VET_RESULTS_DATA_FILE, index=False)
                except:
                    print("IGNORING ERROR!")
                    continue

            except:
                print("IGNORING ERROR!")
                continue
    except:
        print("IGNORING ERROR!")
        continue