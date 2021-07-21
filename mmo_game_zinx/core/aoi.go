package core

import "fmt"

// 定义AOI区域边界值的宏
const (
	AOI_MIN_X  int = 85
	AOI_MAX_X  int = 410
	AOI_CNTS_X int = 10
	AOI_MIN_Y  int = 75
	AOI_MAX_Y  int = 400
	AOI_CNTS_Y int = 20
)

/*
	AOI区域管理模块, 用于管理当前区域内每一个Grid
*/
type AOIManager struct {
	// 区域的左边界坐标
	MinX int
	// 区域的右边界坐标
	MaxX int
	// 区域在X方向上格子的数量
	CntsX int
	// 区域的上边界坐标
	MinY int
	// 区域的下边界坐标
	MaxY int
	// 区域在Y方向上格子的数量
	CntsY int
	// 区域中有哪些九宫格, map：key = 格子的ID, value = 格子对象
	grids map[int]*Grid
}

// 初始化AOI管理区域
func NewAOIManager(minX, maxX, cntsX, minY, maxY, cntsY int) *AOIManager {
	aoiMgr := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		CntsX: cntsX,
		MinY:  minY,
		MaxY:  maxY,
		CntsY: cntsY,
		grids: make(map[int]*Grid),
	}

	// 给AOI区域中每个小格子进行编号和初始化
	// 横轴为x轴, 纵轴为y轴, 因此y行x列
	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			// 根据x, y计算格子ID
			gid := y*cntsX + x
			// 给每个格子初始化
			aoiMgr.grids[gid] = NewGrid(gid,
				aoiMgr.MinX+x*aoiMgr.gridWidth(),
				aoiMgr.MinX+(x+1)*aoiMgr.gridWidth(),
				aoiMgr.MinY+y*aoiMgr.gridLength(),
				aoiMgr.MinY+(y+1)*aoiMgr.gridLength())
		}
	}

	return aoiMgr
}

// 得到每个格子在X轴方向上的宽度
func (m *AOIManager) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CntsX
}

// 得到每个格子在Y周方向上的高度
func (m *AOIManager) gridLength() int {
	return (m.MaxY - m.MinY) / m.CntsY
}

// 调试信息, 用于打印调试AOI区域的信息
func (m *AOIManager) String() string {
	s := fmt.Sprintf("AOIManager:\n MinX:%d, MaxX:%d, CntsX:%d, MinY:%d, MaxY:%d, CntsY:%d\nGrids in AOImanager:\n",
		m.MinX, m.MaxX, m.CntsX, m.MinY, m.MaxY, m.CntsY)

	for _, grid := range m.grids {
		s += fmt.Sprintln(grid)
	}
	return s
}

// 根据当前GID得到周边九宫格格子的GID集合
func (m *AOIManager) GetSuroundGridsByGid(gID int) (grids []*Grid) {
	// 判断当前GID是否在AOIManager中
	if _, ok := m.grids[gID]; !ok {
		return
	}

	// 先将当前gID所在的格子本身加入到九宫格切片集合中
	grids = append(grids, m.grids[gID])

	// 通过gID得出当前格子gID在x轴的编号 idx = gID % cntsX
	idx := gID % m.CntsX

	// 判断idx编号左边是否还有格子, 如果有则加入gidsX集合中
	if idx > 0 {
		grids = append(grids, m.grids[gID-1])
	}

	// 判断idx编号右边是否还有格子, 如果有则加入gidsX集合中
	if idx < m.CntsX-1 {
		grids = append(grids, m.grids[gID+1])
	}

	// gdisX用于保存当前gID所在九宫格的横轴上的所有格子ID
	gidsX := make([]int, 0, len(grids))
	for _, v := range grids {
		//fmt.Println("---------id = ", id)
		gidsX = append(gidsX, v.GID)
	}

	// 遍历gidsX集合, 判断集合中每个gid上边和下边是否还有格子
	for _, v := range gidsX {
		idy := gID / m.CntsY
		if idy > 0 {
			grids = append(grids, m.grids[v-m.CntsX])
		}
		if idy < m.CntsY-1 {
			grids = append(grids, m.grids[v+m.CntsX])
		}

	}
	return
}

// 通过横纵坐标得到当前坐标所在gID的编号
func (m *AOIManager) GetGidByPos(x, y float32) int {
	idx := (int(x) - m.MinX) / m.gridWidth()
	idy := (int(y) - m.MinY) / m.gridLength()

	return idy*m.CntsX + idx
}

// 通过横纵坐标得到周边九宫格内全部PlayersID
func (m *AOIManager) GetPidsbyPos(x, y float32) (playerIDs []int) {
	// 通过横纵坐标得到当前玩家所在格子的ID
	gID := m.GetGidByPos(x, y)

	// 根据gID得到周边九宫格格子的gID集合
	gids := m.GetSuroundGridsByGid(gID)

	// 将九宫格里的全部player的id信息加入到playerIDs
	for _, grid := range gids {
		playerIDs = append(playerIDs, grid.GetPlayerIDs()...)
		//fmt.Printf("------->grid ID: %d, pids: %v\n", grid.GID, grid.GetPlayerIDs())
	}
	return
}

// 添加一个PlayerID到一个格子中
func (m *AOIManager) AddPidToGrid(pID, gID int) {
	m.grids[gID].Add(pID)
}

// 移除一个格子中的PlayerID
func (m *AOIManager) RemovePidFromGrid(pID, gID int) {
	m.grids[gID].Remove(pID)
}

// 通过GID获取全部的PlayerID
func (m *AOIManager) GetPidsByGid(gID int) (playerIDs []int) {
	playerIDs = m.grids[gID].GetPlayerIDs()
	return
}

// 通过坐标将Player添加到一个格子中
func (m *AOIManager) AddToGridByPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	m.AddPidToGrid(pID, gID)
}

// 通过坐标把一个Player从一个格子中删除
func (m *AOIManager) RemoveFromGridByPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	m.RemovePidFromGrid(pID, gID)
}
