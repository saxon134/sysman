package sm

type MenuType int8

const (
	MenuTypeNull MenuType = 0 //未知
	Menu         MenuType = 1 //菜单
	Button       MenuType = 2 //按钮
)

func (m MenuType) Int() int {
	return int(m)
}
