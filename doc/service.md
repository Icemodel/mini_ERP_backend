# Naming

1. ชื่อของ Service ทั้งตัว struct และชื่อไฟล์ต้องเหมือนกัน เช่น ถ้าหาก func ชื่อ UserById ชื่อไฟล์ต้องเป็น user_by_id.go
1. ในการตั้งชื่อตัวแปรและ function ถ้าเป็น public ให้ใช้ PascalCase แต่ถ้าเป็น private ให้ใช้ camelCase
1. สำหรับ Service ที่ใช้ในการค้นหาไม่ต้องมีคำว่า 'Get'
    |    ✅Do |    ❌Don't |
    |:--------:|:-----------:|
    | UserById | GetUserById |  
1. ถ้าหาก Service มีเงื่อนไขที่ต้องใช้ในการค้นหา เช่น ต้องใช้ id ในการค้นหา ให้เพิ่มเงื่อนไขนั้นเข้าไปในชื่อด้วย เช่น UserById, DeedById เป็นต้น
1. ถ้าหากเป็น struct ที่ใช้สำหรับรับ request จาก caller ให้ลงท้ายด้วยคำว่า 'Request' เช่น UserRequest, DeedRequest
1. ถ้าหากเป็น struct ที่ใช้สำหรับ return กลับไปให้ caller ให้ลงท้ายด้วยคำว่า 'Result' เช่น UserResult, DeedResult
1. ถ้าหากเป็น struct ที่ใช้สำหรับเก็บข้อมูลนอกเหนือจาก request และ result ให้เติมคำว่า 'Model' เข้าไปในชื่อด้วย เช่น UserModel, DeedModel โดยแบ่งเป็น case ได้ ดังนี้
    1. ถ้า struct ที่สร้างมาใหม่ชื่อไม่ซ้ำกับที่อื่นใน package เลย สามารถใช้ชื่อที่ตั้งมาแล้วลงท้ายด้วยคำว่า Model ได้เลย
    ```go
    // Service ชื่อ Users
    type CompanyModel struct {
        Roles []UsersRoleModel
        // Fields
    }

    type RoleModel struct {
        // Fields
    }
    ```
    2. ถ้า struct ที่สร้างมาใหม่ชื่อไปซ้ำกับที่อื่นซึ่งอยู่ใน package เดียวกัน ให้เติมชื่อ Service เข้าไปข้างหน้า และลงท้างด้วยคำว่า Model 
    ```go
    // Service ชื่อ Users
    type UsersCompanyModel struct {
        Roles []UsersRoleModel
        // Fields
    }

    type UsersRoleModel struct {
        // Fields
    }

    type UsersDeedModel struct {
        // Fields
    }
    ```
    3. ถ้าใน Service มี struct ที่เอาไว้ใช้ใน result ที่ต้องการ return เป็น array เพียงตัวเดียว ให้เปลี่ยนชื่อ Service ให้อยู่ในรูปเอกพจน์แล้วเติมไปข้างหน้า struct
    ```go
    // Service ชื่อ Companies

    // struct ของ result
    type CompaniesByGroupIdResult struct {
        Companies []CompanyByGroupIdModel
        // Fields
    }

    type CompanyByGroupIdModel struct {
        // Fields
    }
    ```
1. ไม่ใช้ struct ข้าม Service
1. ห้ามใช้ struct จาก Repository เพื่อใช้สำหรับ return ถึงแม้ว่า field ข้างในจะเหมือนกันก็ตาม แต่สามารถใช้ struct จาก Entity ได้

<p style="text-align:center;font-size:1.5em;color:red;">ถ้ามี case นอกเหนือจากนี้แนะนำให้ปรึกษากับทีมก่อนเริ่มเขียน</p>