package model

import "MiaoSha/dao"

type Store struct {
	Id   string
	Name string
	Sum  int
}

func GetStoreSumByID(id string) (int, error) {
	var store Store
	err := dao.DB.Where("id = ?", id).First(&store).Error
	if err != nil {
		return 0, err
	}
	return store.Sum, nil
}
func DecreaseStoreByID(vocher_id string) error {
	store := Store{}
	if err := dao.DB.Where("vocher_id = ?", vocher_id).First(&store).Error; err != nil {
		return err
	}
	// 库存减一
	store.Sum -= 1

	// 更新库存记录
	if err := dao.DB.Save(&store).Error; err != nil {
		return err
	}
	return nil
}

func queryOrderCount(id string) (int, error) {
	var store Store
	err := dao.DB.Where("id = ?", id).First(&store).Error
	if err != nil {
		return 0, err
	}
	return store.Sum, nil
}
