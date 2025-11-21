package common

const (
	CaptchaCache = "quebec:core:captcha:cache:%s"
)

type MenuType int8

const (
	MenuTypeDirectory MenuType = 1
	MenuTypeMenu      MenuType = 2
	MenuTypeButton    MenuType = 3
)

type DataRelationshipType int8

const (
	// 角色与菜单
	DataRelationshipTypeRoleToMenu DataRelationshipType = 1
	// 用户与角色
	DataRelationshipTypeUserToRole DataRelationshipType = 2
	// 菜单与按钮
	DataRelationshipTypeManyToMany DataRelationshipType = 3
)
