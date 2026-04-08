package enumtypes

type DatasetPermissionType string

const (
	DatasetPermission_ONLY_ME      DatasetPermissionType = "only_me"
	DatasetPermission_ALL_TEAM     DatasetPermissionType = "all_team_members"
	DatasetPermission_PARTIAL_TEAM DatasetPermissionType = "partial_members"
)
