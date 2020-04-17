import csv
import sys
from dataclasses import dataclass
from os import path

offset = int(sys.argv[1])
length = int(sys.argv[2])

PROJECT_DATA_FILE = "dev-data/projects.csv"
MODULES_DATA_FILE = "dev-data/modules_{}_{}.csv".format(offset, offset + length - 1)
MATCHES_DATA_FILE = "dev-data/unsafe_matches_{}_{}.csv".format(offset, offset + length - 1)
VET_RESULTS_DATA_FILE = "dev-data/vet_results_{}_{}.csv".format(offset, offset + length - 1)
ERRORS_DATA_FILE = "dev-data/errors_{}_{}.csv".format(offset, offset + length - 1)

MATCH_TYPES = (
    'unsafe.Pointer', 'unsafe.Sizeof', 'unsafe.Alignof', 'unsafe.Offsetof',
    'uintptr', 'reflect.SliceHeader', 'reflect.StringHeader'
)


@dataclass
class _Module:
    project_name: str
    module_import_path: str
    module_registry: str
    module_version: str
    module_number_go_files: int
    module_checkout_folder: str


class Module(_Module):
    file = None
    writer = None

    def write(self):
        if not Module.file:
            Module.open()

        Module.writer.writerow([self.project_name, self.module_import_path, self.module_registry,
                                self.module_version, self.module_number_go_files, self.module_checkout_folder])

        Module.file.flush()

    @staticmethod
    def open():
        write_header = False

        if path.exists(MODULES_DATA_FILE):
            with open(MODULES_DATA_FILE, 'r') as f:
                if not f.readline().startswith('project_name,'):
                    write_header = True
        else:
            write_header = True

        if write_header:
            with open(MODULES_DATA_FILE, 'w') as f:
                f.write("project_name,module_import_path,module_registry,module_version,module_number_go_files," +
                        "module_checkout_folder\n")

        Module.file = open(MODULES_DATA_FILE, 'w')
        Module.writer = csv.writer(Module.file, delimiter=',', quotechar='"')

    @staticmethod
    def close():
        if Module.file:
            Module.file.close()


@dataclass
class _MatchResult:
    project_name: str
    module_import_path: str
    module_registry: str
    module_version: str
    module_number_go_files: int
    module_checkout_folder: str
    file_name: str
    file_size_bytes: int
    file_size_lines: int
    file_imports_unsafe_pkg: bool
    file_go_vet_output: str
    text: str
    context: str
    line_number: int
    byte_offset: int
    match_type: str


class MatchResult(_MatchResult):
    file = None
    writer = None

    def write(self):
        if not MatchResult.file:
            MatchResult.open()

        MatchResult.writer.writerow([self.project_name, self.module_import_path, self.module_registry,
                                     self.module_version, self.module_number_go_files, self.module_checkout_folder,
                                     self.file_name, self.file_size_bytes, self.file_size_lines,
                                     self.file_imports_unsafe_pkg, self.file_go_vet_output, self.text, self.context,
                                     self.line_number, self.byte_offset, self.match_type])

        MatchResult.file.flush()

    @staticmethod
    def open():
        write_header = False

        if path.exists(MATCHES_DATA_FILE):
            with open(MATCHES_DATA_FILE, 'r') as f:
                if not f.readline().startswith('project_name,'):
                    write_header = True
        else:
            write_header = True

        if write_header:
            with open(MATCHES_DATA_FILE, 'w') as f:
                f.write("project_name,module_import_path,module_registry,module_version,module_number_go_files," +
                        "module_checkout_folder,file_name,file_size_bytes,file_size_lines,file_imports_unsafe_pkg," +
                        "file_go_vet_output,text,context,line_number,byte_offset,match_type\n")

        MatchResult.file = open(MATCHES_DATA_FILE, 'w')
        MatchResult.writer = csv.writer(MatchResult.file, delimiter=',', quotechar='"')

    @staticmethod
    def close():
        if MatchResult.file:
            MatchResult.file.close()


@dataclass
class _VetResult:
    project_name: str
    module_import_path: str
    module_registry: str
    module_version: str
    module_number_go_files: int
    module_checkout_folder: str
    file_name: str
    file_size_bytes: int
    file_size_lines: int
    file_imports_unsafe_pkg: bool
    file_go_vet_output: str
    line_number: int
    message: str
    raw_output: str


class VetResult(_VetResult):
    file = None
    writer = None

    def write(self):
        if not VetResult.file:
            VetResult.open()

        VetResult.writer.writerow([self.project_name, self.module_import_path, self.module_registry,
                                   self.module_version, self.module_number_go_files, self.module_checkout_folder,
                                   self.file_name, self.file_size_bytes, self.file_size_lines,
                                   self.file_imports_unsafe_pkg, self.file_go_vet_output, self.line_number,
                                   self.message, self.raw_output])

        VetResult.file.flush()

    @staticmethod
    def open():
        write_header = False

        if path.exists(VET_RESULTS_DATA_FILE):
            with open(VET_RESULTS_DATA_FILE, 'r') as f:
                if not f.readline().startswith('project_name,'):
                    write_header = True
        else:
            write_header = True

        if write_header:
            with open(VET_RESULTS_DATA_FILE, 'w') as f:
                f.write("project_name,module_import_path,module_registry,module_version,module_number_go_files," +
                        "module_checkout_folder,file_name,file_size_bytes,file_size_lines,file_imports_unsafe_pkg," +
                        "file_go_vet_output,line_number,message,raw_output\n")

        VetResult.file = open(VET_RESULTS_DATA_FILE, 'w')
        VetResult.writer = csv.writer(VetResult.file, delimiter=',', quotechar='"')

    @staticmethod
    def close():
        if VetResult.file:
            VetResult.file.close()


@dataclass
class _ErrorCondition:
    stage: str
    project_name: str
    module_import_path: str
    file_name: str
    message: str


class ErrorCondition(_ErrorCondition):
    file = None
    writer = None

    def write(self):
        if not ErrorCondition.file:
            ErrorCondition.open()

        ErrorCondition.writer.writerow([self.stage, self.project_name, self.module_import_path,
                                        self.file_name, self.message])

        ErrorCondition.file.flush()

    @staticmethod
    def open():
        write_header = False

        if path.exists(ERRORS_DATA_FILE):
            with open(ERRORS_DATA_FILE, 'r') as f:
                if not f.readline().startswith('stage,'):
                    write_header = True
        else:
            write_header = True

        if write_header:
            with open(ERRORS_DATA_FILE, 'w') as f:
                f.write("stage,project_name,module_import_path,file_name,message\n")

        ErrorCondition.file = open(ERRORS_DATA_FILE, 'w')
        ErrorCondition.writer = csv.writer(ErrorCondition.file, delimiter=',', quotechar='"')

    @staticmethod
    def close():
        if ErrorCondition.file:
            ErrorCondition.file.close()
