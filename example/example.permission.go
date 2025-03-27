package example

type PermissionCode struct {
	Create string `value:"create_example"`
	View   string `value:"view_example"`
	Update string `value:"update_example"`
}

var Permissions PermissionCode
