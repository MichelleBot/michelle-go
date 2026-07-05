package core

type User struct {
	JID           string `db:"jid" json:"jid"`
	LID           string `db:"lid" json:"lid"`
	Name          string `db:"name" json:"name"`
	Banned        bool   `db:"banned" json:"banned"`
	BanTemporary  int64  `db:"ban_temporary" json:"ban_temporary"`
	BanTimes      int    `db:"ban_times" json:"ban_times"`
	Point         int    `db:"point" json:"point"`
	Balance       int    `db:"balance" json:"balance"`
	Pocket        int    `db:"pocket" json:"pocket"`
	Deposito      int    `db:"deposito" json:"deposito"`
	Guard         int    `db:"guard" json:"guard"`
	LastClaim     int64  `db:"last_claim" json:"last_claim"`
	LastRob       int64  `db:"last_rob" json:"last_rob"`
	Premium       bool   `db:"premium" json:"premium"`
	Expired       int64  `db:"expired" json:"expired"`
	LastNotified  int64  `db:"last_notified" json:"last_notified"`
	LastSeen      int64  `db:"last_seen" json:"last_seen"`
	LastJadibot   int64  `db:"last_jadibot" json:"last_jadibot"`
	Hit           int    `db:"hit" json:"hit"`
	Warning       int    `db:"warning" json:"warning"`
	Attempt       int    `db:"attempt" json:"attempt"`
	Code          string `db:"code" json:"code"`
	CodeExpire    int64  `db:"code_expire" json:"code_expire"`
	Email         string `db:"email" json:"email"`
	Verified      bool   `db:"verified" json:"verified"`
	Taken         bool   `db:"taken" json:"taken"`
	Partner       string `db:"partner" json:"partner"`
	SavingHistory string `db:"saving_history" json:"saving_history"`
	PlayerData    string `db:"player_data" json:"player_data"`
}

type Group struct {
	JID          string `db:"jid" json:"jid"`
	Name         string `db:"name" json:"name"`
	Activity     int64  `db:"activity" json:"activity"`
	Adzan        bool   `db:"adzan" json:"adzan"`
	Antibot      bool   `db:"antibot" json:"antibot"`
	Antiporn     bool   `db:"antiporn" json:"antiporn"`
	Antidelete   bool   `db:"antidelete" json:"antidelete"`
	Antilink     bool   `db:"antilink" json:"antilink"`
	Antiphishing bool   `db:"antiphishing" json:"antiphishing"`
	Antitagsw    bool   `db:"antitagsw" json:"antitagsw"`
	Antivirtex   bool   `db:"antivirtex" json:"antivirtex"`
	Antiforward  bool   `db:"antiforward" json:"antiforward"`
	Antisticker  bool   `db:"antisticker" json:"antisticker"`
	Adminonly    bool   `db:"adminonly" json:"adminonly"`
	Captcha      bool   `db:"captcha" json:"captcha"`
	Filter       bool   `db:"filter" json:"filter"`
	Game         bool   `db:"game" json:"game"`
	Mysterybox   bool   `db:"mysterybox" json:"mysterybox"`
	Left         bool   `db:"left" json:"left"`
	Localonly    bool   `db:"localonly" json:"localonly"`
	ListData     string `db:"list_data" json:"list_data"`
	Mute         bool   `db:"mute" json:"mute"`
	Autosticker  bool   `db:"autosticker" json:"autosticker"`
	Restrict     bool   `db:"restrict" json:"restrict"`
	MemberData   string `db:"member_data" json:"member_data"`
	TextLeft     string `db:"text_left" json:"text_left"`
	TextWelcome  string `db:"text_welcome" json:"text_welcome"`
	Welcome      bool   `db:"welcome" json:"welcome"`
	Expired      int64  `db:"expired" json:"expired"`
	LastNotified int64  `db:"last_notified" json:"last_notified"`
	Blocked      string `db:"blocked" json:"blocked"`
	Blacklist    string `db:"blacklist" json:"blacklist"`
	Stay         bool   `db:"stay" json:"stay"`
	OpenAt       string `db:"open_at" json:"open_at"`
	CloseAt      string `db:"close_at" json:"close_at"`
}

type Chat struct {
	JID       string `db:"jid" json:"jid"`
	Chat      int    `db:"chat" json:"chat"`
	LastSeen  int64  `db:"last_seen" json:"last_seen"`
	LastReply int64  `db:"last_reply" json:"last_reply"`
}
