package consts

// Names are predefined names used internally by desugarers and compilers.
var Names = struct {
	DictionaryFunction string
	EmptyDictionary    string
	EmptyList          string
	IndexFunction      string
	ListFunction       string
}{
	DictionaryFunction: "$dictionary",
	EmptyDictionary:    "$emptyDictionary",
	EmptyList:          "$emptyList",
	IndexFunction:      "$@",
	ListFunction:       "$list",
}

// FileExtension is a file extension of the language.
const FileExtension = ".cloe"

// ModuleFilename is the name of the top level scripts in module directories
// with no file extension.
const ModuleFilename = "module"
