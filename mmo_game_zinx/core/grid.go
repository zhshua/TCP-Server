package core

import (
	"fmt"
	"sync"
)

/*
	AOI区域中每个小格子的数据类型, 九个小格子构成一个九宫格
	下面的例子中, 格子编号为 0 - 24, 每个格子占50个坐标单位
	0                   250
0	----------------------------> x
    | 0	| 1	| 2	| 3	| 4	|
	---------------------
	| 5	| 6	| 7	| 8	| 9	|
	---------------------
	|10	|11	|12	|13	|14	|
	---------------------
	|15	|16	|17	|18	|19	|
	---------------------
	|20	|21	|22	|23	|24	|
250	---------------------
	y
*/
type Grid struct {
	// 格子ID
	GID int
	/* 下面(MinX, MinY)坐标为每个小格子左上角的坐标
	 * 下面(MaxX, MinY)坐标为每个小格子右上角的坐标
	 * 下面(MinX, MaxY)坐标为每个小格子左下角的坐标
	 * 下面(MaxX, MaxY)坐标为每个小格子右下角的坐标
	 */
	// 格子的左边边界坐标
	MinX int
	// 格子的右边边界坐标
	MaxX int
	// 格子的上边边界坐标
	MinY int
	// 格子的下边边界坐标
	MaxY int
	// 当前小格子内玩家或者物体成员的ID集合
	playerIDs map[int]bool
	// 保护当前集合的锁
	pIDLock sync.RWMutex
}

// 初始化当前格子的方法
func NewGrid(gID, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		GID:       gID,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		playerIDs: make(map[int]bool),
	}
}

// 给格子添加一个玩家
func (g *Grid) Add(playerID int) {
	// 加锁
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	g.playerIDs[playerID] = true
}

// 从格子中删除一个玩家
func (g *Grid) Remove(playerID int) {
	// 加锁
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	delete(g.playerIDs, playerID)
}

// 得到当前小格子中所有的玩家
func (g *Grid) GetPlayerIDs() (playerIDs []int) {
	// 加读锁
	g.pIDLock.RLock()
	defer g.pIDLock.RUnlock()

	for k, _ := range g.playerIDs {
		playerIDs = append(playerIDs, k)
	}
	return
}

// 调试使用-打印出格子的基本信息
func (g *Grid) String() string {
	return fmt.Sprintf("Grid id: %d, minX: %d, maxX: %d, minY: %d, maxY: %d, playerIDs: %v",
		g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIDs)
}
