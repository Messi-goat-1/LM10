package redis

import (
	"redis/interfaces"
	"sync"
)

type Redis struct {
	// مستودع البيانات الأساسي
	redisDbs interfaces.RedisDbs

	// قفل الحماية للعمليات المتزامنة
	mu *sync.RWMutex

	// مدير انتهاء صلاحية المفاتيح
	keyExpirer interfaces.ExpirationManager
}

// دالة الوصول للقفل
func (r *Redis) Mu() *sync.RWMutex {
	return r.mu
}

// الوصول لمدير الصلاحية
func (r *Redis) KeyExpirer() interfaces.ExpirationManager {
	return r.keyExpirer
}

/*

// إنشاء السيرفر مع قناة تتسع لـ 100 دفعة مثلاً
server := &redis.AnalyticsServer{
    DataQueue: make(chan *pb.DataBatch, 100),
}

// تشغيل المحرك في الخلفية قبل تشغيل الـ gRPC
go server.StartWorker()
*/
