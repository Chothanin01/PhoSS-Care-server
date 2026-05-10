package databases

import (
	"fmt"
	"os"

	"github.com/chothanin01/PhoSS-Care-server/pkg/utils"
	"gorm.io/gorm"
)

func SeedSuperAdmin(db *gorm.DB, passwordSvc utils.PasswordService) error {
	var count int64
	if err := db.Model(&User{}).Where("role = ?", "admin").Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check existing admins: %w", err)
	}

	if count > 0 {
		return nil
	}

	username := os.Getenv("SUPERADMIN_USERNAME")
	password := os.Getenv("SUPERADMIN_PASSWORD")
	firstname := os.Getenv("SUPERADMIN_FIRSTNAME")
	lastname := os.Getenv("SUPERADMIN_LASTNAME")

	if username == "" || password == "" {
		return fmt.Errorf("missing SUPERADMIN_USERNAME or SUPERADMIN_PASSWORD in .env")
	}

	hashed, err := passwordSvc.Hash(password)
	if err != nil {
		return fmt.Errorf("failed to hash superadmin password: %w", err)
	}

	user := User{
		Username: username,
		Password: hashed,
		Role:     "admin",
	}

	if err := db.Create(&user).Error; err != nil {
		return fmt.Errorf("failed to create superadmin user: %w", err)
	}

	admin := Admin{
		Title:     "Mr.",
		FirstName: firstname,
		LastName:  lastname,
		UserID:    user.ID,
	}

	if err := db.Create(&admin).Error; err != nil {
		return fmt.Errorf("failed to create superadmin profile: %w", err)
	}

	fmt.Printf("Superadmin created: %s (env-based)\n", username)
	return nil
}

func SeedDiseases(db *gorm.DB) error {
	var superadmin User
	superadminUsername := os.Getenv("SUPERADMIN_USERNAME")

	if superadminUsername == "" {
		return fmt.Errorf("missing SUPERADMIN_USERNAME env; cannot assign CreatedBy/UpdatedBy for diseases")
	}

	if err := db.Where("username = ?", superadminUsername).First(&superadmin).Error; err != nil {
		return fmt.Errorf("failed to find superadmin (username=%s): %w", superadminUsername, err)
	}

	defaultDiseases := []Disease{
		{Name: "โรคความดันโลหิตสูง", AvailableDays: StringArray{"Tuseday"}},
		{Name: "โรคเบาหวาน", AvailableDays: StringArray{"Wednesday"}},
		{Name: "วัณโรค", AvailableDays: StringArray{"Friday"}},
		{Name: "วัคซีน", AvailableDays: StringArray{"Monday"}},
	}

	for _, d := range defaultDiseases {
		var count int64
		if err := db.Model(&Disease{}).Where("name = ?", d.Name).Count(&count).Error; err != nil {
			return fmt.Errorf("failed to check disease %s: %w", d.Name, err)
		}

		if count == 0 {
			d.CreatedBy = &superadmin.ID
			d.UpdatedBy = &superadmin.ID
			if err := db.Create(&d).Error; err != nil {
				return fmt.Errorf("failed to seed disease %s: %w", d.Name, err)
			}
			fmt.Printf("Seeded disease: %s (by superadmin: %s)\n", d.Name, superadmin.Username)
		}
	}

	return nil
}

func SeedVaccines(db *gorm.DB) error {
	var superadmin User
	superadminUsername := os.Getenv("SUPERADMIN_USERNAME")

	if superadminUsername == "" {
		return fmt.Errorf("missing SUPERADMIN_USERNAME env; cannot assign createby/updateby for vaccines")
	}

	if err := db.Where("username = ?", superadminUsername).First(&superadmin).Error; err != nil {
		return fmt.Errorf("failed to find superadmin (username=%s): %w", superadminUsername, err)
	}

	defaultVaccines := []Vaccine{
		{Name: "BCG (วัคซีนป้องกันวัณโรค)", Age: "แรกเกิด", Type: "Live Attenuated", Effect: "ป้องกันวัณโรค", Note: "ฉีดก่อนออกจากโรงพยาบาล"},
		{Name: "HB (วัคซีนป้องกันโรคตับอักเสบบี)", Age: "แรกเกิด, 1 เดือน", Type: "Inactivated", Effect: "ป้องกันโรคตับอักเสบบี", Note: "เข็มแรกควรให้ภายใน 24 ชั่วโมงหลังคลอด"},
		{Name: "DTP-HB-Hib (วัคซีนรวม 5 โรค)", Age: "2, 4, 6 เดือน", Type: "Combination", Effect: "ป้องกันโรคคอตีบ บาดทะยัก ไอกรน ตับอักเสบบี และฮิบ", Note: ""},
		{Name: "IPV (วัคซีนป้องกันโรคโปลิโอชนิดฉีด)", Age: "2, 4 เดือน", Type: "Inactivated", Effect: "ป้องกันโรคโปลิโอ", Note: ""},
		{Name: "Rota (วัคซีนโรต้า)", Age: "2, 4, 6 เดือน", Type: "Live Attenuated (Oral)", Effect: "ป้องกันโรคอุจจาระร่วงจากไวรัสโรต้า", Note: "ชนิดหยอด ครั้งแรกต้องให้ก่อนอายุ 15 สัปดาห์"},
		{Name: "OPV (วัคซีนป้องกันโรคโปลิโอชนิดรับประทาน)", Age: "6 เดือน, 1 ปี 6 เดือน, 4 ปี", Type: "Live Attenuated (Oral)", Effect: "ป้องกันโรคโปลิโอ", Note: ""},
		{Name: "MMR (วัคซีนรวมป้องกันโรคหัด-คางทูม-หัดเยอรมัน)", Age: "9 เดือน, 1 ปี 6 เดือน", Type: "Live Attenuated", Effect: "ป้องกันโรคหัด คางทูม และหัดเยอรมัน", Note: ""},
		{Name: "LAJE (วัคซีนป้องกันโรคไข้สมองอักเสบเจอี)", Age: "1 ปี, 2 ปี 6 เดือน", Type: "Live Attenuated", Effect: "ป้องกันโรคไข้สมองอักเสบเจอี", Note: "ชนิดเชื้อเป็นอ่อนฤทธิ์"},
		{Name: "DTP (วัคซีนรวมป้องกันโรคคอตีบ-บาดทะยัก-ไอกรน)", Age: "1 ปี 6 เดือน, 4 ปี", Type: "Combination", Effect: "ป้องกันโรคคอตีบ บาดทะยัก และไอกรน", Note: "เข็มกระตุ้น"},
		{Name: "HPV (วัคซีนป้องกันมะเร็งปากมดลูก)", Age: "11 ปี (ป.5)", Type: "Inactivated", Effect: "ป้องกันมะเร็งปากมดลูกจากเชื้อเอชพีวี", Note: "ฉีด 2 เข็ม ห่างกันอย่างน้อย 6 เดือน (สำหรับเด็กหญิง)"},
		{Name: "dT (วัคซีนรวมป้องกันโรคคอตีบ-บาดทะยัก)", Age: "12 ปี (ป.6)", Type: "Toxoid", Effect: "ป้องกันโรคคอตีบและบาดทะยัก", Note: "เข็มกระตุ้น"},
	}

	for _, v := range defaultVaccines {
		var count int64
		if err := db.Model(&Vaccine{}).Where("name = ?", v.Name).Count(&count).Error; err != nil {
			return fmt.Errorf("failed to check vaccine %s: %w", v.Name, err)
		}

		if count == 0 {
			v.CreatedBy = &superadmin.ID
			v.UpdatedBy = &superadmin.ID
			if err := db.Create(&v).Error; err != nil {
				return fmt.Errorf("failed to seed vaccine %s: %w", v.Name, err)
			}
		}
	}

	return nil
}

func SeedDoctors(db *gorm.DB) error {

    defaultDoctors := []Doctor{
        {Title: "นพ.", FirstName: "สมชาย", LastName: "ใจดี", Role: "doctor"},
        {Title: "พญ.", FirstName: "สมหญิง", LastName: "รักเรียน", Role: "doctor"},
        {Title: "นพ.", FirstName: "อาคม", LastName: "ขยันทำดี", Role: "doctor"},
        {Title: "พญ.", FirstName: "วิไล", LastName: "พรสวัสดิ์", Role: "doctor"},
		{Title: "พย.", FirstName: "ปราณี", LastName: "เมตตา", Role: "nurse"},
		{Title: "พย.", FirstName: "น้ำฝน", LastName: "ชื่นใจ", Role: "nurse"},
		{Title: "พย.", FirstName: "ดวงแก้ว", LastName: "ห่วงใย", Role: "nurse"},
		{Title: "พย.", FirstName: "พิมพ์ชนก", LastName: "รักษ์ดี", Role: "nurse"},
    }

    for _, doc := range defaultDoctors {
        var count int64
        err := db.Model(&Doctor{}).
            Where("first_name = ? AND last_name = ?", doc.FirstName, doc.LastName).
            Count(&count).Error
        
        if err != nil {
            return fmt.Errorf("failed to check doctor %s: %w", doc.FirstName, err)
        }

        if count == 0 {
            if err := db.Create(&doc).Error; err != nil {
                return fmt.Errorf("failed to seed doctor %s %s: %w", doc.FirstName, doc.LastName, err)
            }
        }
    }

    return nil
}