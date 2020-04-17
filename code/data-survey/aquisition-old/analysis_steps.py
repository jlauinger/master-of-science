import json

from .data_structures import *
from .extractors import *


def analyze_modules(project):
    modules = get_project_modules(project)

    # TODO: analyze project main directory, too

    for module_import_path in modules:
        try:
            analyze_module(project, module_import_path)

        except Exception as e:
            ErrorCondition(stage='module',
                           project_name=project['project_name'],
                           module_import_path=module_import_path,
                           message="{}".format(e),
                           file_name='').write()
            print('SAVING AN ERROR!')
            continue


def analyze_module(project, module_import_path):
    module_checkout_folder = project['project_checkout_path'] + "/vendor/" + module_import_path

    files = get_files_in_folder(module_checkout_folder)

    module = Module(project_name=project['project_name'],
                    module_import_path=module_import_path,
                    module_registry=extract_registry(module_import_path),
                    module_version='',
                    module_number_go_files=len(files),
                    module_checkout_folder=module_checkout_folder)

    module.write()

    print("  " + module_import_path + " (" + str(module.module_number_go_files) + " files)...")

    # gosec_output = get_gosec_output(module_checkout_folder)

    for file in files:
        try:
            analyze_file(file, module)

        except Exception as e:
            ErrorCondition(stage='file',
                           project_name=project['project_name'],
                           module_import_path=module_import_path,
                           file_name=extract_filename_from_module(file, module_checkout_folder),
                           message="{}".format(e)).write()
            print('SAVING AN ERROR!')
            continue


def analyze_file(file, module):
    go_vet_output = get_go_vet_output(file)

    file_data = get_file_features(file, module.module_checkout_folder)
    file_data['file_go_vet_output'] = go_vet_output

    # TODO: save LOC feature

    # TODO: save gosec output

    analyze_vet_findings(go_vet_output, file_data, module)

    run_grep(file_data, module)


def analyze_vet_findings(go_vet_output, file_data, module):
    # TODO: remove trailing newline!

    vet_findings = [finding for finding in go_vet_output.split("\n") if len(finding) > 0]

    for vet_finding in vet_findings:
        try:
            analyze_vet_finding(vet_finding, file_data, module)

        except Exception as e:
            ErrorCondition(stage='vet',
                           project_name=module.project_name,
                           module_import_path=module.module_import_path,
                           file_name=file_data['file_name'],
                           message="{}".format(e)).write()
            print('SAVING AN ERROR!')
            continue


def analyze_vet_finding(vet_finding, file_data, module):
    line_number, message = parse_vet_finding(vet_finding)

    result = VetResult(**module.__dict__,
                       **file_data,
                       line_number=line_number,
                       message=message,
                       raw_output=vet_finding)

    result.write()


def run_grep(file_data, module):
    for match_type in MATCH_TYPES:
        try:
            grep_output = get_grep_output(module.module_checkout_folder, file_data['file_name'], match_type)

            grep_messages = [json.loads(message) for message
                             in grep_output.split("\n")
                             if len(message) > 0]

            analyze_grep_findings(grep_messages, file_data, module, match_type)

        except Exception as e:
            ErrorCondition(stage='grep',
                           project_name=module.project_name,
                           module_import_path=module.module_import_path,
                           file_name=file_data['file_name'],
                           message="{}".format(e)).write()
            print('SAVING AN ERROR!')
            continue


def analyze_grep_findings(grep_messages, file_data, module, match_type):
    for i, line in enumerate(grep_messages):
        try:
            if line['type'] == 'match':
                analyze_grep_finding(i, line, grep_messages, file_data, module, match_type)

        except Exception as e:
            ErrorCondition(stage='grep-analyze',
                           project_name=module.project_name,
                           module_import_path=module.module_import_path,
                           file_name=file_data['file_name'],
                           message="{}".format(e)).write()
            print('SAVING AN ERROR!')
            continue


def analyze_grep_finding(i, line, grep_messages, file_data, module, match_type):
    context_lines = grep_messages[max(0, i-5):min(len(grep_messages), i+1+5)]
    context = "".join([cl['data']['lines']['text'] for cl in context_lines
                       if cl['type'] == 'context' or cl['type'] == 'match'])

    result = MatchResult(**module.__dict__,
                         **file_data,
                         text=line['data']['lines']['text'],
                         context=context,
                         line_number=line['data']['line_number'],
                         byte_offset=line['data']['absolute_offset'],
                         match_type=match_type)

    result.write()
