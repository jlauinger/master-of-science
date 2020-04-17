import os
import subprocess
from os import path

import pandas as pd

from .data_structures import PROJECT_DATA_FILE


def get_projects_data(offset, length):
    projects_df = pd.read_csv(PROJECT_DATA_FILE)
    return projects_df[offset:offset+length]


def go_mod_exists(project):
    return path.exists(project['project_checkout_path'] + "/go.mod")


def get_project_modules(project):
    return [module for module
            in os.popen("bash -c 'cd \"" + project['project_checkout_path'] +
                        "\" && go mod vendor -v 2>&1 | grep -v \"#\" | sort | uniq'").read().split("\n")
            if len(module) > 0 and "warning" not in module]


def get_files_in_folder(folder):
    # TODO: non-recursive find to avoid double counting?
    return [file for file
            in os.popen("find '" + folder + "' -name '*.go'").read().split("\n")
            if len(file) > 0]


def extract_registry(import_path):
    path_components = import_path.split('/')
    registry_components = path_components[0:2] if path_components[1] == 'x' else path_components[0:1]
    return "/".join(registry_components)


def get_go_vet_output(file):
    return os.popen("bash -c 'cd \"" + path.dirname(file) + "\" && go vet " +
                    "\"" + path.basename(file) + "\" 2>&1 | grep -v \"#\"'").read()


def parse_vet_finding(vet_finding):
    components = vet_finding.split(":")

    # sometimes they look like this "vet: file.go:42:10: error"
    if components[0] == 'vet':
        line_number = components[2]
        message = ":".join(components[4:]).strip()
    # but sometimes the leading vet is missing: "file.go:42:10: error"
    else:
        line_number = components[1]
        message = ":".join(components[3:]).strip()

    return line_number, message


def get_gosec_output(file):
    # TODO: check and adopt gosec output
    return os.popen("bash -c 'cd \"" + path.dirname(file) + "\" && gosec -quiet -fmt=json'").read()


def get_grep_output(module_path, file_name, match_type):
    return os.popen("bash -c 'cd \"" + module_path + "\" && rg " + match_type + " --context 5 --json \"" +
                    file_name + "\"'").read()


def extract_filename_from_module(file, module_path):
    return file[len(module_path)+1:]


def get_file_features(file, module_path):
    return {
        'file_name': extract_filename_from_module(file, module_path),
        'file_size_bytes': int(os.popen("wc -c '" + file + "'").read().split(" ")[0]),
        'file_size_lines': int(os.popen("wc -l '" + file + "'").read().split(" ")[0]),
        'file_imports_unsafe_pkg': subprocess.call(["grep", "unsafe", file], stdout=subprocess.DEVNULL) == 0,
    }
