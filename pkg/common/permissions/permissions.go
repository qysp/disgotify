package permissions

// PermissionLevel permission level indicator.
type PermissionLevel uint

// Command permission level.
const (
	PermissionDefault PermissionLevel = iota
	PermissionDeveloper
)
