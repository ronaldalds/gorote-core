package core

type PermissionCode struct {
	SuperUser            string `value:"super_user"`
	CreateUser           string `value:"create_user"`
	ViewUser             string `value:"view_user"`
	UpdateUser           string `value:"update_user"`
	EditePermissionsUser string `value:"edite_permissions_user"`
	CreateRole           string `value:"create_role"`
	ViewRole             string `value:"view_role"`
	UpdateRole           string `value:"update_role"`
}

var Permissions PermissionCode
