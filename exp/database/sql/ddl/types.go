package ddlz

type StmtVerb string

const (
	VerbUnknown StmtVerb = "UNKNOWN_VERB"
	VerbCreate  StmtVerb = "CREATE"
	VerbAlter   StmtVerb = "ALTER"
	VerbDrop    StmtVerb = "DROP"
)

type StmtResource string

const (
	StmtTypeUnknown StmtResource = "UNKNOWN_RESOURCE"
	StmtTypeTable   StmtResource = "TABLE"
	StmtTypeIndex   StmtResource = "INDEX"
)
