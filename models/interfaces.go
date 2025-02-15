package models

type DBGenerator interface {
	GetCsvName() string
	GetDBData() []string
	GetDBDataHeader() []string
}

func GetDBData() []string {
	return []string{}
}

func GetDBDataHeader() []string {
	return []string{}
}

func GetCsvName() string {
	return ""
}

type Response interface {
	GetResponse() Response
}
