package base

import (
	"time"
)

// CSV output formats -----------------------------------------------------------------------------------

type DateTime struct {
	time.Time
}

type ProjectData struct {
	Name           string   `csv:"project_name"`
	Rank           int      `csv:"project_rank"`
	GithubCloneUrl string   `csv:"project_github_clone_url"`
	NumberOfStars  int      `csv:"project_number_of_stars"`
	NumberOfForks  int      `csv:"project_number_of_forks"`
	GithubId       int64    `csv:"project_github_id"`
	Revision       string   `csv:"project_revision"`
	CreatedAt      DateTime `csv:"project_created_at"`
	LastPushedAt   DateTime `csv:"project_last_pushed_at"`
	UpdatedAt      DateTime `csv:"project_updated_at"`
	Size           int      `csv:"project_size"`
	CheckoutPath   string   `csv:"project_checkout_path"`
	UsesModules    bool     `csv:"project_uses_modules"`
	RootModule     string   `csv:"project_root_module"`
}

type PackageData struct {
	Name            string `csv:"name"`
	ImportPath      string `csv:"import_path"`
	Dir             string `csv:"dir"`
	IsStandard      bool   `csv:"is_standard"`
	IsDepOnly       bool   `csv:"is_dep_only"`
	NumberOfGoFiles int    `csv:"number_of_go_files"`
	Loc             int    `csv:"loc"`
	ByteSize        int    `csv:"byte_size"`

	ModulePath       string `csv:"module_path"`
	ModuleVersion    string `csv:"module_version"`
	ModuleRegistry   string `csv:"module_registry"`
	ModuleIsIndirect bool   `csv:"module_is_indirect"`

	ProjectName string `csv:"project_name"`

	GoFiles     []string `csv:"-"`
	Imports     []string `csv:"-"`
	Deps        []string `csv:"-"`
	ImportStack []string `csv:"-"`
	HopCount    int      `csv:"package_hop_count"`

	UnsafeSum                  int  `csv:"package_unsafe_sum"`

	UnsafePointerSum        int  `csv:"package_geiger_unsafe_pointer_sum"`
	UnsafePointerVariable   int  `csv:"package_geiger_unsafe_pointer_variable"`
	UnsafePointerParameter  int  `csv:"package_geiger_unsafe_pointer_parameter"`
	UnsafePointerAssignment int  `csv:"package_geiger_unsafe_pointer_assignment"`
	UnsafePointerCall       int  `csv:"package_geiger_unsafe_pointer_call"`
	UnsafePointerOther      int  `csv:"package_geiger_unsafe_pointer_other"`

	UnsafeSizeofSum        int  `csv:"package_geiger_unsafe_sizeof_sum"`
	UnsafeSizeofVariable   int  `csv:"package_geiger_unsafe_sizeof_variable"`
	UnsafeSizeofParameter  int  `csv:"package_geiger_unsafe_sizeof_parameter"`
	UnsafeSizeofAssignment int  `csv:"package_geiger_unsafe_sizeof_assignment"`
	UnsafeSizeofCall       int  `csv:"package_geiger_unsafe_sizeof_call"`
	UnsafeSizeofOther      int  `csv:"package_geiger_unsafe_sizeof_other"`

	UnsafeOffsetofSum        int  `csv:"package_geiger_unsafe_offsetof_sum"`
	UnsafeOffsetofVariable   int  `csv:"package_geiger_unsafe_offsetof_variable"`
	UnsafeOffsetofParameter  int  `csv:"package_geiger_unsafe_offsetof_parameter"`
	UnsafeOffsetofAssignment int  `csv:"package_geiger_unsafe_offsetof_assignment"`
	UnsafeOffsetofCall       int  `csv:"package_geiger_unsafe_offsetof_call"`
	UnsafeOffsetofOther      int  `csv:"package_geiger_unsafe_offsetof_other"`

	UnsafeAlignofSum        int  `csv:"package_geiger_unsafe_alignof_sum"`
	UnsafeAlignofVariable   int  `csv:"package_geiger_unsafe_alignof_variable"`
	UnsafeAlignofParameter  int  `csv:"package_geiger_unsafe_alignof_parameter"`
	UnsafeAlignofAssignment int  `csv:"package_geiger_unsafe_alignof_assignment"`
	UnsafeAlignofCall       int  `csv:"package_geiger_unsafe_alignof_call"`
	UnsafeAlignofOther      int  `csv:"package_geiger_unsafe_alignof_other"`

	SliceHeaderSum        int  `csv:"package_geiger_slice_header_sum"`
	SliceHeaderVariable   int  `csv:"package_geiger_slice_header_variable"`
	SliceHeaderParameter  int  `csv:"package_geiger_slice_header_parameter"`
	SliceHeaderAssignment int  `csv:"package_geiger_slice_header_assignment"`
	SliceHeaderCall       int  `csv:"package_geiger_slice_header_call"`
	SliceHeaderOther      int  `csv:"package_geiger_slice_header_other"`

	StringHeaderSum        int  `csv:"package_geiger_string_header_sum"`
	StringHeaderVariable   int  `csv:"package_geiger_string_header_variable"`
	StringHeaderParameter  int  `csv:"package_geiger_string_header_parameter"`
	StringHeaderAssignment int  `csv:"package_geiger_string_header_assignment"`
	StringHeaderCall       int  `csv:"package_geiger_string_header_call"`
	StringHeaderOther      int  `csv:"package_geiger_string_header_other"`

	UintptrSum        int  `csv:"package_geiger_uintptr_sum"`
	UintptrVariable   int  `csv:"package_geiger_uintptr_variable"`
	UintptrParameter  int  `csv:"package_geiger_uintptr_parameter"`
	UintptrAssignment int  `csv:"package_geiger_uintptr_assignment"`
	UintptrCall       int  `csv:"package_geiger_uintptr_call"`
	UintptrOther      int  `csv:"package_geiger_uintptr_other"`
}

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
	PackageDir           string `csv:"package_dir"`
	ModulePath           string `csv:"module_path"`
	ModuleVersion        string `csv:"module_version"`
	ProjectName          string `csv:"project_name"`
}

type GrepFindingData struct {
	Text                 string `csv:"text"`
	Context              string `csv:"context"`
	LineNumber           int    `csv:"line_number"`
	Column               int    `csv:"column"`
	AbsoluteOffset       int    `csv:"absolute_offset"`
	MatchType            string `csv:"match_type"`

	FileName             string `csv:"file_name"`
	FileLoc              int    `csv:"file_loc"`
	FileByteSize         int    `csv:"file_byte_size"`
	PackageImportPath    string `csv:"package_import_path"`
	ModulePath           string `csv:"module_path"`
	ModuleVersion        string `csv:"module_version"`
	ProjectName          string `csv:"project_name"`

	FileCopyPath         string `csv:"file_copy_path"`
}

type VetFindingData struct {
	Message              string `csv:"message"`
	Context              string `csv:"context"`
	LineNumber           int    `csv:"line_number"`
	Column               int    `csv:"column"`
	RawOutput            string `csv:"raw_output"`

	FileName             string `csv:"file_name"`
	FileLoc              int    `csv:"file_loc"`
	FileByteSize         int    `csv:"file_byte_size"`
	PackageImportPath    string `csv:"package_import_path"`
	ModulePath           string `csv:"module_path"`
	ModuleVersion        string `csv:"module_version"`
	ProjectName          string `csv:"project_name"`

	FileCopyPath         string `csv:"file_copy_path"`
}

type GosecFindingData struct {
	Message              string `csv:"message"`
	Context              string `csv:"context"`
	Confidence           string `csv:"confidence"`
	Severity             string `csv:"severity"`
	CweId                string `csv:"cwe_id"`
	RuleId               string `csv:"rule_id"`
	LineNumber           int    `csv:"line_number"`
	Column               int    `csv:"column"`

	FileName             string `csv:"file_name"`
	FileLoc              int    `csv:"file_loc"`
	FileByteSize         int    `csv:"file_byte_size"`
	PackageImportPath    string `csv:"package_import_path"`
	ModulePath           string `csv:"module_path"`
	ModuleVersion        string `csv:"module_version"`
	ProjectName          string `csv:"project_name"`

	FileCopyPath         string `csv:"file_copy_path"`
}

type GosaferFindingData struct {
	Message              string `csv:"message"`
	Context              string `csv:"context"`
	LineNumber           int    `csv:"line_number"`
	Column               int    `csv:"column"`
	RawOutput            string `csv:"raw_output"`

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


// Go list parsing -----------------------------------------------------------------------------------------------------

type GoListOutputPackage struct {
	Dir           string   // directory containing package sources
	ImportPath    string   // import path of package in dir
	ImportComment string   // path in import comment on package statement
	Name          string   // package name
	Doc           string   // package documentation string
	Target        string   // install path
	Shlib         string   // the shared library that contains this package (only set when -linkshared)
	Goroot        bool     // is this package in the Go root?
	Standard      bool     // is this package part of the standard Go library?
	Stale         bool     // would 'go install' do anything for this package?
	StaleReason   string   // explanation for Stale==true
	Root          string   // Go root or Go path dir containing this package
	ConflictDir   string   // this directory shadows Dir in $GOPATH
	BinaryOnly    bool     // binary-only package (no longer supported)
	ForTest       string   // package is only for use in named test
	Export        string   // file containing export data (when using -export)
	Module        *GoListOutputModule  // info about package's containing module, if any (can be nil)
	Match         []string // command-line patterns matching this package
	DepOnly       bool     // package is only a dependency, not explicitly listed

	// Source files
	GoFiles         []string // .go source files (excluding CgoFiles, TestGoFiles, XTestGoFiles)
	CgoFiles        []string // .go source files that import "C"
	CompiledGoFiles []string // .go files presented to compiler (when using -compiled)
	IgnoredGoFiles  []string // .go source files ignored due to build constraints
	CFiles          []string // .c source files
	CXXFiles        []string // .cc, .cxx and .cpp source files
	MFiles          []string // .m source files
	HFiles          []string // .h, .hh, .hpp and .hxx source files
	FFiles          []string // .f, .F, .for and .f90 Fortran source files
	SFiles          []string // .s source files
	SwigFiles       []string // .swig files
	SwigCXXFiles    []string // .swigcxx files
	SysoFiles       []string // .syso object files to add to archive
	TestGoFiles     []string // _test.go files in package
	XTestGoFiles    []string // _test.go files outside package

	// Cgo directives
	CgoCFLAGS    []string // cgo: flags for C compiler
	CgoCPPFLAGS  []string // cgo: flags for C preprocessor
	CgoCXXFLAGS  []string // cgo: flags for C++ compiler
	CgoFFLAGS    []string // cgo: flags for Fortran compiler
	CgoLDFLAGS   []string // cgo: flags for linker
	CgoPkgConfig []string // cgo: pkg-config names

	// Dependency information
	Imports      []string          // import paths used by this package
	ImportMap    map[string]string // map from source import to ImportStack (identity entries omitted)
	Deps         []string          // all (recursively) imported dependencies
	TestImports  []string          // imports from TestGoFiles
	XTestImports []string          // imports from XTestGoFiles

	// Error information
	Incomplete bool            // this package or a dependency has an error
	Error      *GoListOutputPackageError   // error loading package
	DepsErrors []*GoListOutputPackageError // errors loading dependencies
}

type GoListOutputModule struct {
	Path      string       // module path
	Version   string       // module version
	Versions  []string     // available module versions (with -versions)
	Replace   *GoListOutputModule      // replaced by this module
	Time      *time.Time   // time version was created
	Update    *GoListOutputModule      // available update, if any (with -u)
	Main      bool         // is this the main module?
	Indirect  bool         // is this module only an indirect dependency of main module?
	Dir       string       // directory holding files for this module, if any
	GoMod     string       // path to go.mod file used when loading this module, if any
	GoVersion string       // go version used in module
	Error     *GoListOutputModuleError // error loading module
}

type GoListOutputModuleError struct {
	Err string // the error itself
}

type GoListOutputPackageError struct {
	ImportStack   []string // shortest path from package named on command line to this one
	Pos           string   // position of error (if present, file:line:col)
	Err           string   // the error itself
}


// Ripgrep parsing -----------------------------------------------------------------------------------------------------

type RipgrepOutputLine struct {
	MessageType string 				`json:"type"` // begin,end,match,context,summary
	Data        *RipgrepMessageData	`json:"data"`
}

type RipgrepMessageData struct {
	Path           *RipgrepTextEncapsulation `json:"path"`
	Lines          *RipgrepTextEncapsulation `json:"lines"`
	LineNumber     int                       `json:"line_number"`
	AbsoluteOffset int                       `json:"absolute_offset"`
	SubMatches     []RipgrepSubmatch         `json:"submatches"`

	BinaryOffset   int                       `json:"binary_offset"`
	Stats   	   *RipgrepEndStats          `json:"stats"`
	ElapsedTotal   *RipgrepTime              `json:"elapsed_total"`
}

type RipgrepTextEncapsulation struct {
	Text string `json:"text"`
}

type RipgrepSubmatch struct {
	Match *RipgrepTextEncapsulation `json:"match"`
	Start int                       `json:"start"`
	End   int                       `json:"end"`
}

type RipgrepEndStats struct {
	Elapsed           *RipgrepTime `json:"elapsed"`
	Searches          int          `json:"searches"`
	SearchesWithMatch int          `json:"searches_with_match"`
	BytesSearched     int          `json:"bytes_searched"`
	BytesPrinted      int          `json:"bytes_printed"`
	MatchedLines      int          `json:"matched_lines"`
	Matches           int          `json:"matches"`
}

type RipgrepTime struct {
	Human string `json:"human"`
	Nanos int    `json:"nanos"`
	Secs  int    `json:"secs"`
}


// Vet parsing ---------------------------------------------------------------------------------------------------------

type VetFindingLine struct {
	Message     string
	ContextLine string
}


// Gosec parsing -------------------------------------------------------------------------------------------------------

type GosecOutput struct {
	GolangErrors *GosecGolangErrorsOutput `json:"Golang errors"`
	Issues       []GosecIssueOutput       `json:"Issues"`
	Stats        *GosecStatsOutput        `json:"Stats"`
}

type GosecIssueOutput struct {
	Severity   string          `json:"severity"`
	Confidence string          `json:"confidence"`
	Cwe        *GosecCweOutput `json:"cwe"`
	RuleId     string          `json:"rule_id"`
	Details    string          `json:"details"`
	File       string          `json:"file"`
	Code       string          `json:"code"`
	Line       string          `json:"line"`
	Column     string          `json:"column"`
}

type GosecCweOutput struct {
	Id  string `json:"ID"`
	Url string `json:"URL"`
}

type GosecGolangErrorsOutput struct {}

type GosecStatsOutput struct {
	Files int `json:"files"`
	Lines int `json:"lines"`
	Nosec int `json:"nosec"`
	Found int `json:"found"`
}


// Go-safer parsing ------------------------------------------------------------------------------------------------------

type GosaferFindingLine struct {
	Message     string
	ContextLine string
}


// AST analysis data ---------------------------------------------------------------------------------------------------

type AstFindingData struct {
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

type AstFunctionData struct {
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

type AstStatementData struct {
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