package status

import (
	"github.com/cossim/hipush/consts"
	"github.com/cossim/hipush/store"
)

type StateStorage struct {
	store store.Store
}

func NewStateStorage(store store.Store) *StateStorage {
	return &StateStorage{
		store: store,
	}
}

func (s *StateStorage) Init() error {
	return s.store.Init()
}

func (s *StateStorage) Close() error {
	return s.store.Close()
}

// Reset Client storage.
func (s *StateStorage) Reset() {
	s.store.Set(consts.HiPushTotal, 0)
	s.store.Set(consts.HiPushSuccess, 0)
	s.store.Set(consts.HiPushFailed, 0)

	s.store.Set(consts.HTTPTotal, 0)
	s.store.Set(consts.HTTPSuccess, 0)
	s.store.Set(consts.HTTPFailed, 0)

	s.store.Set(consts.GRPCTotal, 0)
	s.store.Set(consts.GRPCSuccess, 0)
	s.store.Set(consts.GRPCFailed, 0)

	s.store.Set(consts.IosTotal, 0)
	s.store.Set(consts.IosSuccess, 0)
	s.store.Set(consts.IosFailed, 0)

	s.store.Set(consts.HuaweiTotal, 0)
	s.store.Set(consts.HuaweiSuccess, 0)
	s.store.Set(consts.HuaweiFailed, 0)

	s.store.Set(consts.AndroidTotal, 0)
	s.store.Set(consts.AndroidSuccess, 0)
	s.store.Set(consts.AndroidFailed, 0)

	s.store.Set(consts.VivoTotal, 0)
	s.store.Set(consts.VivoSuccess, 0)
	s.store.Set(consts.VivoFailed, 0)

	s.store.Set(consts.OppoTotal, 0)
	s.store.Set(consts.OppoSuccess, 0)
	s.store.Set(consts.OppoFailed, 0)

	s.store.Set(consts.XiaomiTotal, 0)
	s.store.Set(consts.XiaomiSuccess, 0)
	s.store.Set(consts.XiaomiFailed, 0)

	s.store.Set(consts.MeizuTotal, 0)
	s.store.Set(consts.MeizuSuccess, 0)
	s.store.Set(consts.MeizuFailed, 0)
}

func (s *StateStorage) AddTotalCount(count int64) {
	s.store.Add(consts.HiPushTotal, count)
}

func (s *StateStorage) AddIosTotal(count int64) {
	s.store.Add(consts.IosTotal, count)
}

func (s *StateStorage) AddIosSuccess(count int64) {
	s.store.Add(consts.IosSuccess, count)
}

func (s *StateStorage) AddIosFailed(count int64) {
	s.store.Add(consts.IosFailed, count)
}

func (s *StateStorage) AddAndroidTotal(count int64) {
	s.store.Add(consts.IosTotal, count)
}

func (s *StateStorage) AddAndroidSuccess(count int64) {
	s.store.Add(consts.AndroidSuccess, count)
}

func (s *StateStorage) AddAndroidFailed(count int64) {
	s.store.Add(consts.AndroidFailed, count)
}

func (s *StateStorage) AddHuaweiTotal(count int64) {
	s.store.Add(consts.HuaweiTotal, count)
}

func (s *StateStorage) AddHuaweiSuccess(count int64) {
	s.store.Add(consts.HuaweiSuccess, count)
}

func (s *StateStorage) AddHuaweiFailed(count int64) {
	s.store.Add(consts.HuaweiFailed, count)
}

func (s *StateStorage) AddXiaomiTotal(count int64) {
	s.store.Add(consts.XiaomiTotal, count)
}

func (s *StateStorage) AddXiaomiSuccess(count int64) {
	s.store.Add(consts.XiaomiSuccess, count)
}

func (s *StateStorage) AddXiaomiFailed(count int64) {
	s.store.Add(consts.XiaomiFailed, count)
}

func (s *StateStorage) AddOppoTotal(count int64) {
	s.store.Add(consts.OppoTotal, count)
}

func (s *StateStorage) AddOppoSuccess(count int64) {
	s.store.Add(consts.OppoSuccess, count)
}

func (s *StateStorage) AddOppoFailed(count int64) {
	s.store.Add(consts.OppoFailed, count)
}

func (s *StateStorage) AddVivoTotal(count int64) {
	s.store.Add(consts.VivoTotal, count)
}

func (s *StateStorage) AddVivoSuccess(count int64) {
	s.store.Add(consts.VivoSuccess, count)
}

func (s *StateStorage) AddVivoFailed(count int64) {
	s.store.Add(consts.VivoFailed, count)
}

func (s *StateStorage) AddMeizuTotal(count int64) {
	s.store.Add(consts.MeizuTotal, count)
}

func (s *StateStorage) AddMeizuSuccess(count int64) {
	s.store.Add(consts.MeizuSuccess, count)
}

func (s *StateStorage) AddMeizuFailed(count int64) {
	s.store.Add(consts.MeizuFailed, count)
}

func (s *StateStorage) AddHonorTotal(count int64) {
	s.store.Add(consts.HonorTotal, count)
}

func (s *StateStorage) AddHonorSuccess(count int64) {
	s.store.Add(consts.HonorSuccess, count)
}

func (s *StateStorage) AddHonorFailed(count int64) {
	s.store.Add(consts.HonorFailed, count)
}

// GetTotalCount show counts of all notification.
func (s *StateStorage) GetTotalCount() int64 {
	return s.store.Get(consts.HiPushTotal)
}

func (s *StateStorage) GetIosTotal() int64 {
	return s.store.Get(consts.IosTotal)
}

// GetIosSuccess show success counts of iOS notification.
func (s *StateStorage) GetIosSuccess() int64 {
	return s.store.Get(consts.IosSuccess)
}

// GetIosFailed show Failed counts of iOS notification.
func (s *StateStorage) GetIosFailed() int64 {
	return s.store.Get(consts.IosFailed)
}

func (s *StateStorage) GetAndroidTotal() int64 {
	return s.store.Get(consts.AndroidTotal)
}

func (s *StateStorage) GetAndroidSuccess() int64 {
	return s.store.Get(consts.AndroidSuccess)
}

func (s *StateStorage) GetAndroidFailed() int64 {
	return s.store.Get(consts.AndroidFailed)
}

func (s *StateStorage) GetHuaweiTotal() int64 {
	return s.store.Get(consts.AndroidTotal)
}

func (s *StateStorage) GetHuaweiSuccess() int64 {
	return s.store.Get(consts.HuaweiSuccess)
}

func (s *StateStorage) GetHuaweiFailed() int64 {
	return s.store.Get(consts.HuaweiFailed)
}

func (s *StateStorage) GetXiaomiTotal() int64 {
	return s.store.Get(consts.XiaomiTotal)
}

func (s *StateStorage) GetXiaomiSuccess() int64 {
	return s.store.Get(consts.XiaomiSuccess)
}

func (s *StateStorage) GetXiaomiFailed() int64 {
	return s.store.Get(consts.XiaomiFailed)
}

func (s *StateStorage) GetVivoTotal() int64 {
	return s.store.Get(consts.VivoTotal)
}

func (s *StateStorage) GetVivoSuccess() int64 {
	return s.store.Get(consts.VivoSuccess)
}

func (s *StateStorage) GetVivoFailed() int64 {
	return s.store.Get(consts.VivoFailed)
}

func (s *StateStorage) GetOppoTotal() int64 {
	return s.store.Get(consts.OppoTotal)
}

func (s *StateStorage) GetOppoSuccess() int64 {
	return s.store.Get(consts.OppoSuccess)
}

func (s *StateStorage) GetOppoFailed() int64 {
	return s.store.Get(consts.OppoFailed)
}

func (s *StateStorage) GetMeizuTotal() int64 {
	return s.store.Get(consts.MeizuTotal)
}

func (s *StateStorage) GetMeizuSuccess() int64 {
	return s.store.Get(consts.MeizuSuccess)
}

func (s *StateStorage) GetMeizuFailed() int64 {
	return s.store.Get(consts.MeizuFailed)
}
