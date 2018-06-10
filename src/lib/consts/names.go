package consts

// Names are predefined names used internally by desugarers and compilers.
var Names = struct {
	DictionaryFunction string
	EmptyDictionary    string
	EmptyList          string
	ListFunction       string
}{
	DictionaryFunction: "$dictionary",
	EmptyDictionary:    "$emptyDictionary",
	EmptyList:          "$emptyList",
	ListFunction:       "$list",
}

// FileExtension is a file extension of the language.
const FileExtension = ".cloe"

// PathName is the name of the language path where modules are stored.
const PathName = "CLOE_PATH"

// ModuleFilename is the name of the top level scripts in module directories
// with no file extension.
const ModuleFilename = "module"
