package constmodels

var (
	// STATUS_ACTIVE means account is not deleted (0)
	STATUS_ACTIVE int32 = 0
	// STATUS_DELETED_BY_SELF means account is deleted by self (1)
	STATUS_DELETED_BY_SELF int32 = 1
	// STATUS_DELETED_BY_MOD means account is deleted by mod (2)
	STATUS_DELETED_BY_MOD int32 = 2
)
