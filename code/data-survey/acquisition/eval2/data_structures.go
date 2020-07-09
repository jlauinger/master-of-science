package eval2

type GeigerFindingData struct {
	Text                 string `csv:"text"`
	Context              string `csv:"context"`
	LineNumber           int    `csv:"line_number"`
	Column               int    `csv:"column"`
	AbsoluteOffset       int    `csv:"absolute_offset"`
	MatchType            string `csv:"match_type"`
	ContextType          string `csv:"context_type"`

	FileName             string `csv:"file_name"`
	FileLoc              int    `csv:"file_loc"`
	FileByteSize         int    `csv:"file_byte_size"`
	PackageImportPath    string `csv:"package_import_path"`
	ModulePath           string `csv:"module_path"`
	ModuleVersion        string `csv:"module_version"`
	ProjectName          string `csv:"project_name"`
}
