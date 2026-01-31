# المرحلة الأولى: بناء التطبيق (Build Stage)
FROM golang:1.24-alpine AS builder

# تثبيت أدوات البناء ومكتبة pcap الضرورية لـ gopacket
RUN apk add --no-cache build-base libpcap-dev

WORKDIR /app

# نسخ ملفات التعريف وتحميل المكتبات
COPY go.mod go.sum ./
RUN go mod download

# نسخ الكود المصدري
COPY . .

# ✅ يجب تفعيل CGO هنا لأن gopacket تعتمد على C
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

# ---------------------------------------------------
# المرحلة الثانية: التشغيل (Runtime Stage)
FROM alpine:latest

# تثبيت مكتبة libpcap فقط لتشغيل البرنامج
RUN apk add --no-cache libpcap

WORKDIR /app

# نسخ الملف التنفيذي من مرحلة البناء
COPY --from=builder /app/server .

# إضافة صلاحية التنفيذ
RUN chmod +x /app/server

# تشغيل السيرفر
CMD ["./server"]