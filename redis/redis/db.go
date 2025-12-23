package redis

import (
	"errors"
	"math"
	"redis/interfaces"
	"sync"
	"time"
)

const (
	keysMapSize = 32
)

// ====================================================================
// RedisDb
// ====================================================================

// ====================================================================
// Types (الأنواع الرئيسية)
// ====================================================================

// لتبسيط المثال، سأعيد تعريف الأنواع هنا مؤقتاً، ولكن يفضل استخدامها من حزمة interfaces
// type DatabaseId uint
// type Keys map[string]Item
// type ExpiringKeys map[string]time.Time

// A redis database.
type RedisDb struct {
	// Database id
	id interfaces.DatabaseId

	// All keys in this db.
	keys interfaces.Keys
	// Keys with expire timestamp.
	expiringKeys interfaces.ExpiringKeys
	redis        *Redis
}

// ====================================================================
// Constructors and Accessors (المنشئات و الوصول)
// ====================================================================

// NewRedisDb creates a new db.
func NewRedisDb(id interfaces.DatabaseId, r *Redis) *RedisDb {
	return &RedisDb{
		id:           id,
		redis:        r,
		keys:         make(interfaces.Keys, keysMapSize),
		expiringKeys: make(interfaces.ExpiringKeys, keysMapSize),
	}
}

// RedisDb gets the redis database by its id or creates and returns it if not exists.
func (r *Redis) RedisDb(dbId interfaces.DatabaseId) *RedisDb {
	getDb := func() *RedisDb {
		if db, ok := r.redisDbs[dbId]; ok {
			if redisDb, ok := db.(*RedisDb); ok {
				return redisDb
			}
		}
		return nil
	}

	r.Mu().RLock()
	db := getDb()
	r.Mu().RUnlock()
	if db != nil {
		return db
	}

	r.Mu().Lock()
	defer r.Mu().Unlock()

	db = getDb()
	if db != nil {
		return db
	}

	r.redisDbs[dbId] = NewRedisDb(dbId, r)
	if redisDb, ok := r.redisDbs[dbId].(*RedisDb); ok {
		return redisDb
	}
	return nil
}

// RedisDbs gets all redis databases.
func (r *Redis) RedisDbs() interfaces.RedisDbs {
	r.Mu().RLock()
	defer r.Mu().RUnlock()
	return r.redisDbs
}

// Redis gets the redis instance.
func (db *RedisDb) Redis() *Redis {
	return db.redis
}

// Mu gets the mutex.
func (db *RedisDb) Mu() *sync.RWMutex {
	return db.Redis().Mu()
}

// Id gets the db id.
func (db *RedisDb) Id() interfaces.DatabaseId {
	return db.id
}

// ====================================================================
// Key/Value Operations (عمليات المفاتيح والقيم)
// ====================================================================

// Sets a key with an item which can have an expiration time.
func (db *RedisDb) Set(key string, i interfaces.Item, expires bool, expiry time.Time) {
	db.Mu().Lock()
	defer db.Mu().Unlock()
	db.keys[key] = i
	if expires {
		db.expiringKeys[key] = expiry
	}
}

// Returns the item by the key or nil if key does not exists.
func (db *RedisDb) Get(key string) interfaces.Item {
	db.Mu().RLock()
	defer db.Mu().RUnlock()
	return db.get(key)
}

func (db *RedisDb) get(key string) interfaces.Item {
	i, _ := db.keys[key]
	return i
}

// Deletes a key, returns number of deleted keys.
func (db *RedisDb) Delete(keys ...string) int {
	db.Mu().Lock()
	defer db.Mu().Unlock()
	return db.delete(keys...)
}

// If checkExists is false, then return bool is reprehensible.
func (db *RedisDb) delete(keys ...string) int {
	do := func(k string) bool {
		if k == "" {
			return false
		}
		i := db.get(k)
		if i == nil {
			return false
		}
		i.OnDelete(k, db)
		delete(db.keys, k)
		delete(db.expiringKeys, k)
		return true
	}

	var c int
	for _, k := range keys {
		if do(k) {
			c++
		}
	}

	return c
}

func (db *RedisDb) DeleteExpired(keys ...string) int {
	var c int
	for _, k := range keys {
		if db.Expired(k) && db.Delete(k) > 0 {
			c++
		}
	}
	return c
}

// GetOrExpire gets the item or nil if expired or not exists. If 'deleteIfExpired' is true the key will be deleted.
func (db *RedisDb) GetOrExpire(key string, deleteIfExpired bool) interfaces.Item {
	db.Mu().Lock()
	defer db.Mu().Unlock()
	i, ok := db.keys[key]
	if !ok {
		return nil
	}
	if db.expired(key) {
		if deleteIfExpired {
			db.delete(key)
		}
		return nil
	}
	return i
}

// IsEmpty checks if db is empty.
func (db *RedisDb) IsEmpty() bool {
	db.Mu().RLock()
	defer db.Mu().RUnlock()
	return len(db.keys) == 0
}

// HasExpiringKeys checks if db has any expiring keys.
func (db *RedisDb) HasExpiringKeys() bool {
	db.Mu().RLock()
	defer db.Mu().RUnlock()
	return len(db.expiringKeys) != 0
}

// Check if key exists.
func (db *RedisDb) Exists(key string) bool {
	db.Mu().RLock()
	defer db.Mu().RUnlock()
	return db.exists(key)
}
func (db *RedisDb) exists(key string) bool {
	_, ok := db.keys[key]
	return ok
}

// Check if key has an expiry set.
func (db *RedisDb) Expires(key string) bool {
	db.Mu().RLock()
	defer db.Mu().RUnlock()
	return db.expires(key)
}
func (db *RedisDb) expires(key string) bool {
	_, ok := db.expiringKeys[key]
	return ok
}

// Expired only check if a key can and is expired.
func (db *RedisDb) Expired(key string) bool {
	db.Mu().RLock()
	defer db.Mu().RUnlock()
	return db.expired(key)
}
func (db *RedisDb) expired(key string) bool {
	return db.expires(key) && TimeExpired(db.expiry(key))
}

// Expiry gets the expiry of the key has one.
func (db *RedisDb) Expiry(key string) time.Time {
	db.Mu().RLock()
	defer db.Mu().RUnlock()
	return db.expiry(key)
}

func (db *RedisDb) expiry(key string) time.Time {
	return db.expiringKeys[key]
}

// Keys gets all keys in this db.
func (db *RedisDb) Keys() interfaces.Keys {
	db.Mu().RLock()
	defer db.Mu().RUnlock()
	return db.keys
}

// ExpiringKeys gets keys with an expiry set and their timeout.
func (db *RedisDb) ExpiringKeys() interfaces.ExpiringKeys {
	db.Mu().RLock()
	defer db.Mu().RUnlock()
	return db.expiringKeys
}

// TimeExpired check if a timestamp is older than now.
func TimeExpired(expireAt time.Time) bool {
	return time.Now().After(expireAt)
}

// تأكيد أن Expirer يلتزم بواجهة KeyExpirer
var _ interfaces.ExpirationManager = (*Expirer)(nil)

type Expirer struct {
	redis *Redis

	done chan bool
}

func NewKeyExpirer(r *Redis) *Expirer {
	return &Expirer{
		redis: r,
		done:  make(chan bool, math.MaxInt32),
	}
}

// Start starts the Expirer.
//
// tick - How fast is the cleaner triggered.
//
// randomKeys - Amount of random expiring keys to get checked.
//
// againPercentage - If more than x% of keys were expired, start again in same tick.
func (e *Expirer) Start(tick time.Duration, randomKeys int, againPercentage int) {
	ticker := time.NewTicker(tick)
	for {
		select {
		case <-ticker.C:
			e.do(randomKeys, againPercentage)
		case <-e.done:
			ticker.Stop()
			return
		}
	}
}

// Stop stops the expirator
func (e *Expirer) Stop() {
	if e.done != nil {
		e.done <- true
		close(e.done)
	}
}

func (e *Expirer) do(randomKeys, againPercentage int) {
	var deletedKeys int

	// اجمع قواعد البيانات اللي فيها مفاتيح منتهية
	dbs := make(map[*RedisDb]struct{})
	for _, dbOps := range e.Redis().RedisDbs() {
		// تأكد أن dbOps فعلاً *RedisDb
		db, ok := dbOps.(*RedisDb)
		if !ok || db == nil {
			continue
		}
		if !db.HasExpiringKeys() {
			continue
		}
		dbs[db] = struct{}{}
	}

	if len(dbs) == 0 {
		return
	}

	for c := 0; c < randomKeys; c++ {
		// اختار قاعدة بيانات عشوائية
		var db *RedisDb
		for d := range dbs {
			db = d
			break
		}
		if db == nil {
			continue
		}

		// اختار مفتاح عشوائي من المفاتيح المنتهية
		var key string
		for k := range db.ExpiringKeys() {
			key = k
			break
		}
		if key == "" {
			continue
		}

		// احذف المفتاح لو فعلاً منتهي
		if db.DeleteExpired(key) > 0 {
			deletedKeys++
		}
	}

	// أعد المحاولة إذا كانت نسبة الحذف أعلى من النسبة المحددة
	if againPercentage > 0 && deletedKeys*100/randomKeys > againPercentage {
		go e.do(randomKeys, againPercentage)
	}
}

// Redis gets the redis instance.
func (e *Expirer) Redis() *Redis {
	return e.redis
}

// Update updates an existing key using the Item's internal Update logic.
// If the key doesn't exist, it can either return an error or create a new one.
func (db *RedisDb) Update(key string, newValue interface{}) error {
	db.Mu().Lock()
	defer db.Mu().Unlock()

	item, ok := db.keys[key]
	if !ok {
		return errors.New("key not found") // أو يمكنك ابتكار منطق لإنشاء Item جديد هنا
	}

	// استدعاء منطق التحديث الخاص بالأداة التحليلية
	return item.Update(newValue)
}
