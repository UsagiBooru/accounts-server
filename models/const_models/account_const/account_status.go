package account_const

var (
	// STATUS_ACTIVE means account is not deleted
	STATUS_ACTIVE int32 = 0
	// STATUS_DELETED_BY_SELF means account is deleted by self
	STATUS_DELETED_BY_SELF int32 = 1
	// STATUS_DELETED_BY_MOD means account is deleted by mod
	STATUS_DELETED_BY_MOD int32 = 2
)
