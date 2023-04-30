package database

type Bungee struct {
	Id      string `gorm:"primaryKey;"`
	Motd    string
	Favicon string `gorm:"size:4096"`
}

// InitBungeeTable - Initialize table (create default entry)
func (s *Mysql) InitBungeeTable() error {
	r := s.client.FirstOrCreate(&Bungee{
		Id: "default",
	}, "id = ?", "default")
	if r.Error != nil {
		return r.Error
	}

	return nil
}

// GetBungeeEntry - Get Bungee Entry
func (s *Mysql) GetBungeeEntry() (Bungee, error) {
	bungee := Bungee{}
	r := s.client.Find(&bungee)
	if r.Error != nil {
		return Bungee{}, r.Error
	}
	return bungee, nil
}

// SetMotd - Set Motd
func (s *Mysql) SetMotd(motd string) error {
	r := s.client.Model(&Bungee{}).Where("id = ?", "default").Update("motd", motd)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

// SetFavicon - Set Favicon
func (s *Mysql) SetFavicon(favicon string) error {
	r := s.client.Model(&Bungee{}).Where("id = ?", "default").Update("favicon", favicon)
	if r.Error != nil {
		return r.Error
	}

	return nil
}
