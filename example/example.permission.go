package example

type PermissionCode struct {
	Create string `value:"create_teletubby"`
	View   string `value:"view_teletubby"`
	Update string `value:"update_teletubby"`
}

var Permissions PermissionCode
