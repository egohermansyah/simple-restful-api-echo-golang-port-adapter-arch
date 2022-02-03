package defined

import "time"

type UserProfile struct {
	UserId     string    `json:"user_id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Otp        int8      `json:"otp"`
	OtpCreated time.Time `json:"otp_created"`
	FcmToken   string    `json:"fcm_token"`
}
