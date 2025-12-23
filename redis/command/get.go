package command

/*
func GetCommand(args interfaces.CommandArgs2, db interfaces.DbOps, resp interfaces.Responder2) {
	if args.Len() != 1 {
		resp.WriteError("ERR wrong number of arguments for 'get' command")
		return
	}

	key := args.Arg(0)

	item := db.GetOrExpire(key, true)
	if item == nil {
		resp.WriteNull()
		return
	}

	val, ok := item.AsString()
	if !ok {
		resp.WriteError("WRONGTYPE Operation against a key holding the wrong kind of value")
		return
	}

	resp.WriteBulkString(val)
}

*/

//هذا لي احتاجة من اوامر لي لازم اسوها ايضا لازم اكون كائن
//Key, Value, Category
//الربط: يقوم باستدعاء دوال مثل AsString (التي سنضيفها) أو Frequency و CreatedAt ليعيد للسيرفر الآخر صورة كاملة عن البيانات.
