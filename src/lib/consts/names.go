package consts

// Names are predefined names used internally by desugarers and compilers.
var Names = struct {
	DictionaryFunction string
	EmptyDictionary    string
	EmptyList          string
	ListFunction       string
}{
	DictionaryFunction: "$dict",
	EmptyDictionary:    "$emptyDict",
	EmptyList:          "$emptyList",
	ListFunction:       "$list",
}

// FileExtension is a file extension of the language.
const FileExtension = ".cloe"
