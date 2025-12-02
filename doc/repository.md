# หลักการโดยรวม
1. การตั้งชื่อไฟล์ ให้ตั้งล้อตาม table เลยแต่เปลี่ยนให้อยู่ในรูปเอกพจน์
    |      Table     |       File      |
    |----------------|-----------------|
    | users          | user.go         |
    | deed_histories | deed_history.go |
1. ปกติแล้วแต่ละไฟล์จะเป็นตัวแทนของแต่ละ table แต่ถ้าหากมี function ที่ใช้ join แล้วตอน return ไม่สามารถใช้ struct จาก Entity ได้ ให้แยก function นั้นออกไปเป็นอีกไฟล์เลย
1. แต่ถ้า function ใช้ join แล้วยังสามารถใช้ struct จาก Entity ในการ return ได้ ไม่ต้องแยกไปเป็นไฟล์ใหม่
1. ถ้าหากเป็น struct ที่ใช้สำหรับรับ request จาก caller ให้ลงท้ายด้วยคำว่า 'Request' เช่น UserRequest, DeedRequest
1. ถ้าหากเป็น struct ที่ใช้สำหรับ return กลับไปให้ caller ให้ลงท้ายด้วยคำว่า 'Result' เช่น UserResult, DeedResult
1. ถ้าหากเป็น struct ที่ใช้สำหรับเก็บข้อมูลนอกเหนือจาก request และ result ให้เติมคำว่า 'Model' เข้าไปในชื่อด้วย เช่น UserModel, DeedModel โดยแบ่งเป็น case ได้ ดังนี้
    1. ถ้า struct ที่สร้างมาใหม่ชื่อไม่ซ้ำกับที่อื่นใน package เลย สามารถใช้ชื่อที่ตั้งมาแล้วลงท้ายด้วยคำว่า Model ได้เลย
    ```go
    // Repository ชื่อ Users
    type CompanyModel struct {
        Roles []UsersRoleModel
        // Fields
    }

    type RoleModel struct {
        // Fields
    }
    ```
    2. ถ้า struct ที่สร้างมาใหม่ชื่อไปซ้ำกับที่อื่นซึ่งอยู่ใน package เดียวกัน ให้เติมชื่อ Repository เข้าไปข้างหน้า และลงท้างด้วยคำว่า Model 
    ```go
    // Repository ชื่อ Users
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
    3. ถ้าใน Repository มี struct ที่เอาไว้ใช้ใน result ที่ต้องการ return เป็น array เพียงตัวเดียว ให้เปลี่ยนชื่อ Repository ให้อยู่ในรูปเอกพจน์แล้วเติมไปข้างหน้า struct
    ```go
    // Repository ชื่อ Companies

    // struct ของ result
    type CompaniesByGroupIdResult struct {
        Companies []CompanyByGroupIdModel
        // Fields
    }

    type CompanyByGroupIdModel struct {
        // Fields
    }
    ```
