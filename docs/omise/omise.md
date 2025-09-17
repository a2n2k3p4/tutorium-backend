# คู่มือการใช้งาน Omise (omise-go) สำหรับโปรเจกต์นี้

เอกสารนี้สรุปวิธีการตั้งค่าและใช้งานระบบรับชำระเงินผ่าน Omise ด้วยไลบรารี `github.com/omise/omise-go` (เวอร์ชันที่ใช้: v1.6.0) ภายในโปรเจกต์นี้ ทั้งฝั่ง API และการทดสอบใช้งาน รวมถึงแนวปฏิบัติที่ถูกต้องปลอดภัย (best practices)

สำคัญ: หลีกเลี่ยงไม่ให้ข้อมูลบัตรเครดิตวิ่งผ่านเซิร์ฟเวอร์ของคุณในสภาพแวดล้อม production ให้ใช้ tokenization จากฝั่ง client เสมอ (Omise.js / Mobile SDK) การทำ tokenization ฝั่ง server ในเอกสารนี้มีไว้เพื่อการทดสอบเท่านั้น เว้นแต่คุณมี PCI‑DSS AoC ที่ถูกต้อง

## ภาพรวมสถาปัตยกรรมในโปรเจกต์
- มี REST endpoints ที่เกี่ยวข้อง:
  - POST `/payments/charge` สร้างรายการชำระเงิน (charge)
  - GET `/payments/transactions` ดูรายการธุรกรรม (พร้อม filter/pagination)
  - GET `/payments/transactions/:id` ดูธุรกรรมตาม `id` ภายในระบบหรือ `charge_id` ของ Omise
  - POST `/payments/transactions/:id/refund` คืนเงินเต็มจำนวนหรือบางส่วน
  - POST `/webhooks/omise` รับ Webhook จาก Omise (สำหรับอัปเดตสถานะธุรกรรมให้ตรงกับความจริง)
  - GET `/health` ตรวจสอบสุขภาพบริการแบบง่าย

- การผนวก Omise client:
  - ระบบจะสร้าง Omise client จาก `OMISE_PUBLIC_KEY` และ `OMISE_SECRET_KEY` แล้ว inject เข้า request context ผ่าน middleware อัตโนมัติ เมื่อ key ถูกตั้งค่า
  - ถ้าไม่ตั้งค่า key จะยังเปิดเซิร์ฟเวอร์ได้ แต่ route ที่ต้องใช้ Omise จะตอบ error

## การตั้งค่า Environment
ตั้งค่าตัวแปรในไฟล์ `.env` (หรือผ่าน environment ของระบบรันจริง)

จำเป็น:
- `OMISE_PUBLIC_KEY` และ `OMISE_SECRET_KEY` (ใช้ test keys ตอนทดสอบ)

ตัวเลือกแนะนำ:
- `PAYMENT_DEFAULT_CURRENCY` ค่าเริ่มต้นสกุลเงิน (เช่น `THB`)
- `PAYMENT_RETURN_URI` สำหรับ flow ที่ต้อง redirect (เช่น Internet Banking/3DS)

ตัวอย่าง `.env` (ทดสอบ):
```
OMISE_PUBLIC_KEY=pkey_test_xxxxxxxxxxxxxxxxxxx
OMISE_SECRET_KEY=skey_test_xxxxxxxxxxxxxxxxxxx
PAYMENT_DEFAULT_CURRENCY=THB
PAYMENT_RETURN_URI=http://localhost:8000/payments/return
```

หมายเหตุ: จำนวนเงินที่ส่งให้ Omise อยู่ในหน่วย satang (เช่น 10000 = 100.00 THB)

## การรันระบบ (Local)
- ใช้ Docker Compose: `docker compose up -d --build`
- เปิด Swagger UI ที่: `http://localhost:8000/swagger/`
- ตรวจสอบ MinIO/DB/pgAdmin ใช้งานได้ตามปกติ (ไม่เกี่ยวโดยตรงกับ Omise แต่เป็นสภาพแวดล้อมของโปรเจกต์)

## การใช้งาน API หลัก

### 1) สร้างรายการชำระเงิน (Charge)
Endpoint: `POST /payments/charge`

Headers แนะนำ:
- `Idempotency-Key`: ใช้ค่าไม่ซ้ำสำหรับแต่ละคำสั่งซื้อ/คำขอ เช่น `ORDER-1234-CHARGE-<random>` เพื่อป้องกัน double charge หาก client ส่งซ้ำ

Payload หลัก (ตัวอย่างกรณีบัตรเครดิตแบบ tokenized):
```json
{
  "amount": 10000,
  "currency": "THB",
  "paymentType": "credit_card",
  "token": "<omise_token_from_frontend>",
  "description": "Order #1234",
  "user_id": 42
}
```

ประเภทที่รองรับ (`paymentType`):
- `credit_card` (แนะนำให้ส่ง `token` ที่สร้างจาก frontend)
- `promptpay` (ระบบจะสร้าง `source` และ `charge` ให้ และคืนข้อมูลสำหรับสร้าง QR)
- `internet_banking` (ต้องระบุ `bank` เช่น `bbl`, `scb`, `bay` และมี `return_uri` สำหรับ redirect กลับ)

ตัวอย่าง PromptPay:
```json
{
  "amount": 5000,
  "currency": "THB",
  "paymentType": "promptpay",
  "description": "Top-up"
}
```
ผลลัพธ์ Charge จะมีข้อมูลที่จำเป็นสำหรับสร้าง QR (เช่น `charge.source.scannable_code`)

ตัวอย่าง Internet Banking:
```json
{
  "amount": 20000,
  "currency": "THB",
  "paymentType": "internet_banking",
  "bank": "bbl",
  "return_uri": "http://localhost:8000/payments/return",
  "description": "Order #5678"
}
```

หมายเหตุสำคัญ:
- Production: ห้ามส่งข้อมูลบัตร (หมายเลข/เดือนปี/CCV) มายัง backend ให้ใช้ token เสมอ
- หน่วยจำนวนเงิน: ใช้ satang
- ใส่ `Idempotency-Key` ทุกครั้ง แก้ปัญหา network retry/double click

### 2) ดูรายการธุรกรรม
- `GET /payments/transactions?user_id=&status=&channel=&limit=&offset=`
  - `status` เช่น `successful`, `failed`, `pending`
  - `channel` เช่น `card`, `promptpay`, `internet_banking_bbl`

### 3) ดูธุรกรรมเฉพาะรายการ
- `GET /payments/transactions/{id}`
  - `{id}` ใส่ได้ทั้ง `id` ภายในระบบ หรือ `charge_id` ของ Omise (เช่น `chrg_test_...`)

### 4) คืนเงิน (Refund)
- `POST /payments/transactions/{id}/refund`
  - body (optional): `{ "amount": 5000 }` เพื่อคืนบางส่วน (satang) ถ้าไม่ส่ง amount = คืนเต็ม
  - แนะนำให้ใส่ `Idempotency-Key` เช่นเดียวกับการ charge

## Webhook (แนะนำให้ตั้งค่าเสมอ)
Endpoint: `POST /webhooks/omise`

การทำงานโดยสรุป:
- ระบบจะ `RetrieveEvent` เพื่อยืนยันความถูกต้อง แล้วดึง `charge_id` จาก event → `RetrieveCharge` เพื่ออ่านสถานะจริง → upsert ลงฐานข้อมูล (idempotent ด้วย `charge_id`)
- ส่งกลับ 200 เมื่อประมวลผลแล้ว หรือ 5xx เพื่อให้ Omise retry ในกรณี error ชั่วคราว

การใช้งานจริง:
1) ใน Omise Dashboard ตั้ง Webhook URL ให้ชี้มาที่ `https://<your-domain>/webhooks/omise` (ทดสอบ local ใช้ `ngrok` เปิด `http://localhost:8000` ออกอินเทอร์เน็ต)
2) รองรับการส่งซ้ำ (idempotent) อยู่แล้ว แต่ควรตั้ง rate-limit ที่ reverse‑proxy/Firewall ด้วย

## แนวปฏิบัติที่สำคัญ (Best Practices)
- Idempotency: ใส่ `Idempotency-Key` ทุกการเรียก charge/refund โดยสร้างจาก business key + random suffix เพื่อลดโอกาสชนกัน
- Tokenization: ใน production ให้รับเฉพาะ `token` ที่สร้างจาก frontend (Omise.js/SDK)
- Status Handling: เช็ก `charge.status` (`pending`/`successful`/`failed`) และรอ webhook เพื่อความแน่นอน
- Amount Unit: ใช้ satang เสมอเมื่อสื่อสารกับ Omise
- Observability: เก็บ log แบบ structured (เช่น `request_id`, `idempotency_key`, `charge_id`, `user_id`, `amount_satang`, `status`)

## ตัวอย่าง cURL
สร้าง charge แบบบัตร (ด้วย token):
```bash
curl -X POST http://localhost:8000/payments/charge \
  -H 'Content-Type: application/json' \
  -H 'Idempotency-Key: ORDER-1234-CHARGE-abc123' \
  -d '{
        "amount": 10000,
        "currency": "THB",
        "paymentType": "credit_card",
        "token": "<omise_token>",
        "description": "Order #1234",
        "user_id": 42
      }'
```

คืนเงินบางส่วน:
```bash
curl -X POST http://localhost:8000/payments/transactions/chrg_test_xxx/refund \
  -H 'Content-Type: application/json' \
  -H 'Idempotency-Key: ORDER-1234-REFUND-abc124' \
  -d '{"amount": 5000}'
```

## Debug & Version (เฉพาะนักพัฒนา)
- เปิด debug ของไลบรารี: สร้าง client แล้วเรียก `client.SetDebug(true)` (โปรดใช้เฉพาะในสภาพแวดล้อมทดสอบ)
- ระบุ API Version: สามารถตั้ง `client.APIVersion = "2015-11-06"` ให้ตรงกับเวอร์ชันที่บัญชีกำหนด (อ้างอิงจากเอกสาร Omise)

## ข้อควรระวังด้านความปลอดภัย
- อย่าบันทึกข้อมูลบัตรเครดิต/ความลับใน log
- จัดการ keys ผ่าน secret manager/CI secret และหมุนเวียนเป็นระยะ
- ป้องกัน DoS ที่ endpoint webhook ด้วย rate‑limit/ACL

## การทดสอบด้วย Test Keys
ตั้งค่า environment:
```bash
export OMISE_PUBLIC_KEY=pkey_test_xxx
export OMISE_SECRET_KEY=skey_test_xxx
```
หรือแก้ที่ไฟล์ `.env` แล้ว `docker compose up -d --build`

สำหรับทดสอบ PromptPay/Internet Banking ให้ดูวิธีการจากเอกสารทางการของ Omise และใช้ test flow ตามที่กำหนด

---
อ้างอิง:
- Omise Go Library: https://pkg.go.dev/github.com/omise/omise-go
- Omise API Docs: https://docs.omise.co/
