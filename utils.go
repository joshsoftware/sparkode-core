package main

type Language string
type Jobtype string

const (
	LanguageC       Language = "gcc"
	LanguageCPP     Language = "g++"
	LanguageRuby    Language = "rb"
	LanguageGo      Language = "go"
	LanguagePython2 Language = "python2"
	LanguagePython3 Language = "python3"
	languageJava    Language = "javac"
)

const (
	JobCompile Jobtype = "compile"
	JobRun     Jobtype = "run"
)

var languageNameToExtension = map[Language]string{
	LanguageC:       ".c",
	LanguageCPP:     ".cpp",
	LanguageRuby:    ".rb",
	LanguageGo:      ".go",
	LanguagePython2: ".py",
	LanguagePython3: ".py",
	languageJava:    ".java",
}
