# Naming

1. ชื่อของ Handler ทั้งตัว func และชื่อไฟล์ต้องเหมือนกัน เช่น ถ้าหาก func ชื่อ UserById ชื่อไฟล์ต้องเป็น user_by_id.go
1. ถ้าหากเป็น struct ที่ใช้สำหรับรับ request จาก frontend ให้ลงท้ายด้วยคำว่า 'Request' เช่น UserRequest, DeedRequest
2. ถ้าหากเป็น struct ที่ใช้สำหรับ return กลับไปให้ frontend ให้ลงท้ายด้วยคำว่า 'Response' เช่น UserResponse, DeedResponse
3. ถ้าหาก struct ที่ใช้ในการรับและส่งค่าจาก frontend มีหน้าตาเหมือนกับ struct ของ __Service__ หรือ __Entity__ สามารถใช้ struct เหล่านั้นแทนได้
4. ถ้าหาก handler มีเงื่อนไขที่ต้องใช้ในการค้นหา เช่น ต้องใช้ id ในการค้นหา ให้เพิ่มเงื่อนไขนั้นเข้าไปในชื่อด้วย เช่น UserById, DeedById เป็นต้น

<p style="text-align:center;font-size:1.5em;color:red;">ถ้ามี case นอกเหนือจากนี้แนะนำให้ปรึกษากับทีมก่อนเริ่มเขียน</p>

# การแยก folder

ดูที่ __feature__ เป็นหลัก ง่ายๆ คือแยกตาม page ของ frontend