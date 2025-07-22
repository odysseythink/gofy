package enumtypes

type DatasetPermissionEnum string

const (
	DatasetPermission_ONLY_ME      DatasetPermissionEnum = "only_me"
	DatasetPermission_ALL_TEAM     DatasetPermissionEnum = "all_team_members"
	DatasetPermission_PARTIAL_TEAM DatasetPermissionEnum = "partial_members"
)
