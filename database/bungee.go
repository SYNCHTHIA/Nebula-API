package database

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
	r := s.client.Model(&Bungee{}).Where("motd = ?", motd).Update("status", motd)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

// SetFavicon - Set Favicon
func (s *Mysql) SetFavicon(favicon string) error {
	r := s.client.Model(&Bungee{}).Where("favicon = ?", favicon).Update("status", favicon)
	if r.Error != nil {
		return r.Error
	}

	return nil
}
