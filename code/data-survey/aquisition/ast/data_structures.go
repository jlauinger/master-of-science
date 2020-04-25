package ast

type FindingData struct {
	MatchType            string `csv:"match_type"`
	LineNumber           int    `csv:"line_number"`
	Column               int    `csv:"column"`
	Text                 string `csv:"text"`

	FileName             string `csv:"file_name"`
	PackageImportPath    string `csv:"package_import_path"`
	ModulePath           string `csv:"module_path"`
	ModuleVersion        string `csv:"module_version"`
	ProjectName          string `csv:"project_name"`
}

type FunctionData struct {
	LineNumber           int    `csv:"line_number"`
	Column               int    `csv:"column"`
	Text                 string `csv:"text"`
	NumberUnsafePointer  int    `csv:"number_unsafe_pointer"`
	NumberUnsafeSizeof   int    `csv:"number_unsafe_sizeof"`
	NumberUnsafeAlignof  int    `csv:"number_unsafe_alignof"`
	NumberUnsafeOffsetof int    `csv:"number_unsafe_offsetof"`
	NumberUintptr        int    `csv:"number_uintptr"`
	NumberSliceHeader    int    `csv:"number_slice_header"`
	NumberStringHeader   int    `csv:"number_string_header"`

	FileName             string `csv:"file_name"`
	PackageImportPath    string `csv:"package_import_path"`
	ModulePath           string `csv:"module_path"`
	ModuleVersion        string `csv:"module_version"`
	ProjectName          string `csv:"project_name"`
}

type StatementData struct {
	LineNumber           int    `csv:"line_number"`
	Column               int    `csv:"column"`
	Text                 string `csv:"text"`
	NumberUnsafePointer  int    `csv:"number_unsafe_pointer"`
	NumberUnsafeSizeof   int    `csv:"number_unsafe_sizeof"`
	NumberUnsafeAlignof  int    `csv:"number_unsafe_alignof"`
	NumberUnsafeOffsetof int    `csv:"number_unsafe_offsetof"`
	NumberUintptr        int    `csv:"number_uintptr"`
	NumberSliceHeader    int    `csv:"number_slice_header"`
	NumberStringHeader   int    `csv:"number_string_header"`

	FileName             string `csv:"file_name"`
	PackageImportPath    string `csv:"package_import_path"`
	ModulePath           string `csv:"module_path"`
	ModuleVersion        string `csv:"module_version"`
	ProjectName          string `csv:"project_name"`
}

type ErrorConditionData struct {
	Stage             string `csv:"stage"`
	ProjectName       string `csv:"project_name"`
	PackageImportPath string `csv:"module_import_path"`
	FileName          string `csv:"file_name"`
	Message           string `csv:"message"`
}
