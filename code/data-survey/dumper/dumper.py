#!/usr/bin/env python

import sys
import os
import hashlib

import pandas as pd

def main():
    # check arguments for validity
    if not os.path.isfile(filename):
        print("{} is not a valid source file.".format(filename))
        return
    if not os.path.isdir(dest_dir):
        print("{} is not a valid destination directory.".format(dest_dir))
        return

    print("[*] dumping data from {} into {}".format(filename, dest_dir))

    data = load_data(filename)

    print("[*] dumping rows...")

    for index, row in data.iterrows():
        dump_row(row, dest_dir)

    print("[+] done!")


def load_data(filename):
    data = pd.read_csv(filename)

    # filter out unclassified snippets
    data = data\
        .where(data['label']!='unclassified')\
        .where(data['label2']!='unclassified')\
        .dropna()

    snippet_count = data['line_number'].count()
    label_count = data['label'].nunique()
    label2_count = data['label2'].nunique()
    combinations_count = data.groupby(['label', 'label2']).ngroups

    print("[*] loaded {} correctly classified snippets".format(snippet_count))
    print("[*] there are {} what-labels and {} purpose-labels, and there are {} combinations".format(
        label_count, label2_count, combinations_count))

    return data


def dump_row(snippet, dest_dir):
    # make sure the directory exists
    dirname = "{}/{}__{}".format(dest_dir, snippet['label2'], snippet['label'])
    if not os.path.exists(dirname):
        os.mkdir(dirname)

    # create the filename
    hash = hashlib.sha256("{}-{}-{}-{}-{}-{}".format(snippet['project_name'], snippet['module_version'],
        snippet['package_import_path'], snippet['file_name'], snippet['line_number'],
        snippet['match_type']).encode('UTF_8')).hexdigest()[:20]
    filename = "{}/{}.txt".format(dirname, hash)

    with open(filename, "w+") as f:
        # dump the data
        f.write(get_file_content(snippet))


def get_file_content(snippet):
    content = ""

    content += "Module: {}\n".format(snippet['module_path'])
    content += "Version: {}\n".format(snippet['module_version'])
    content += "\n"
    content += "Package: {}\n".format(snippet['package_import_path'])
    content += "File: {}\n".format(snippet['file_name'])
    content += "Line: {}\n".format(int(snippet['line_number']))
    content += "\n"
    content += "Imported (possibly among others) by: {}\n".format(snippet['project_name'])
    content += "\n"
    content += "Label 1 (What is happening?): {}\n".format(snippet['label'])
    content += "Label 2 (For what purpose?): {}\n".format(snippet['label2'])
    content += "\n"

    content += "--------------------------------------------------------------\n"
    content += "Snippet line:\n"
    content += "\n"
    content += snippet['text']
    content += "\n"
    content += "--------------------------------------------------------------\n"
    content += "+/- 5 lines context:\n"
    content += "\n"
    content += snippet['context']
    content += "\n"
    content += "--------------------------------------------------------------\n"
    content += "+/- 100 lines context:\n"
    content += "\n"
    content += get_full_context(snippet)
    content += "\n"

    return content


def get_full_context(snippet):
    # check if the snippet is contained within the Go standard library
    if snippet.module_path == "std":
        # if so, the source file will be located within the GOROOT directory
        file_path = "{}/src/{}/{}".format(
            go_lib_path,
            snippet.package_import_path,
            snippet.file_name)
    else:
        # otherwise, it will be in the correct module directory within the GOPATH directory
        file_path = "{}/pkg/mod/{}@{}{}/{}".format(
            go_mod_path,
            snippet.module_path,
            snippet.module_version,
            snippet.package_import_path[len(snippet.module_path):],
            snippet.file_name)

    # provide a fallback content if the file cannot be found for some reason
    if not os.path.exists(file_path):
        return "n/a"

    # get the whole file contents
    with open(file_path, "r") as f:
        content = f.readlines()

    # return only a +/- 100 lines context
    start = max(int(snippet['line_number']) - 100, 0)
    end = min(int(snippet['line_number']) + 100, len(content))

    return "".join(content[start:end])


# entry point: this is run when the file is run on the command line
if __name__ == "__main__":
    # check that the arguments are supplied
    if len(sys.argv) <= 2:
        print("usage: $0 filename dest_dir.")
        sys.exit(1)

    # extract command line arguments
    filename = sys.argv[1]
    dest_dir = sys.argv[2]

    # extract configuration variables from the environment or set defaults
    go_mod_path = os.environ.get('GO_MOD_PATH') or '/home/johannes/.go'
    go_lib_path = os.environ.get('GO_LIB_PATH') or '/usr/lib/go'

    main()
