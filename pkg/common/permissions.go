package common

// PermissionLevel represents the permission level for a command.
type PermissionLevel uint

// Command permission level.
const (
	PermissionDefault PermissionLevel = iota
	PermissionDeveloper
)
