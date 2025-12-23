package redis

import "redis/interfaces"

// اسوي دالة ApplyUpdate فقط
// افترض أنك ستستخدم الأنواع من حزمة interfaces
// import "your/module/path/interfaces"

// بدلاً من السجل المعقد، نكتفي بهيكل يدير تنفيذ الأوامر برمجياً
type CommandService struct {
	db interfaces.DatabaseOps
}

// دالة واحدة تعالج التحديثات القادمة من gRPC
func (s *CommandService) ExecuteUpdate(key string, value interface{}) error {
	item := s.db.Get(key)
	if item == nil {
		// هنا تضع منطق إنشاء العنصر الجديد إذا لم يكن موجوداً
		return nil
	}
	return item.Update(value)
}
