package core

import "sync"

/*
	当前游戏的世界管理模块
*/
type WorldManager struct {
	// 当前世界地图的AOI管理模块
	AoiMgr *AOIManager
	// 当前全部在线的玩家Player
	OnlinePlayers map[int32]*Player
	// 保护Players集合的锁
	pLock sync.RWMutex
}

// 对外的全局世界管理模块句柄
var WorldMgrObj *WorldManager

// 初始化方法
func init() {
	WorldMgrObj = &WorldManager{
		// 初始化AOI管理模块
		AoiMgr: NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_CNTS_X, AOI_MIN_Y, AOI_MAX_Y, AOI_CNTS_Y),
		// 初始化在线玩家集合
		OnlinePlayers: make(map[int32]*Player),
	}
}

// 添加一个玩家到在线玩家集合
func (wm *WorldManager) AddPlayer(player *Player) {
	// 加锁, 添加玩家到在线玩家集合
	wm.pLock.Lock()
	wm.OnlinePlayers[player.Pid] = player
	wm.pLock.Unlock()

	// 将Player同时也添加到AOI管理模块
	wm.AoiMgr.AddToGridByPos(int(player.Pid), player.X, player.Z)
}

// 从在线玩家集合中删除一个玩家
func (wm *WorldManager) RemovePlayer(pid int32) {
	// 加写锁, 先获得当前玩家
	wm.pLock.RLock()
	player := wm.OnlinePlayers[pid]
	wm.pLock.RUnlock()

	// 将玩家从AOI管理模块中删除
	wm.AoiMgr.RemoveFromGridByPos(int(pid), player.X, player.Z)

	// 加写锁, 将玩家从在线玩家集合中删除
	wm.pLock.Lock()
	delete(wm.OnlinePlayers, pid)
	wm.pLock.Unlock()

}

// 通过玩家id查询在线玩家对象
func (wm *WorldManager) GetOnlinePlayerByPid(pid int32) *Player {
	// 加读锁
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()
	return wm.OnlinePlayers[pid]
}

// 获取全部在线玩家
func (wm *WorldManager) GetAllOnlinePlayers() []*Player {
	// 加读锁
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	players := make([]*Player, 0)

	for _, p := range wm.OnlinePlayers {
		players = append(players, p)
	}
	return players
}
