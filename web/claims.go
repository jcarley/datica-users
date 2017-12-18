package web

import (
	"net/http"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jcarley/datica-users/helper/jsonutil"
	"github.com/jcarley/datica-users/helper/structutil"
)

var (
	mySigningKey = []byte(`sfXrQ7g@_bkTFM_gdbA6u$5C6gNVghuK8J6wHHX!j_EEZh-!$GmcsG_*6P^Pg!_ZG&rdrEx*T9?jQ2%zPF9VJXA7$v=g4ZaJkp-ks!VG??$yPUD%vnUeCWk=F?c*6Yeh?mFbUew5BQyMt5WfMZM7CJ=-eA5h2yZ-y6sZ#^Dyysfqt*rB9EXzPyDy7E%K!BBH84Lq9r7Q7^fSq-V3aWgT*d32g&kq6!XTk^m=L=TCZGGD!QPeqU+Sx!#qL*9Y@sfm_JNLsQ2HBt5DU*sS=^#BLF8EHb=-k7gFK&_Z+4UkfmCaD-36ub#=K*R67C+2*KqmBxc!+7Ec9Q%-kw3L5M^TUsPYeE388%kma!?2H8-gL?tfEJh4v3z6ncE4vFR=EnAt-gKcZy5gyV=qbk6H!wVLghDd8Z3Db5&2PVj6T_pkSHeM_?7QEbuY5mMsgP=MUYS^m#^+9%zpwT@avUU7Z8Pgn2bzZRP?6GSrjzd+?C=7rRtZ&fz3R9Fpu?--a*Sx!E^kr9Yj+rheJVH_Vx2^-5V3WG%Rxbb#w+PqbM_+dcnQy&t97Gr!*xAHGj!tF^C@CWwzYyNCuKH=F9mN?dB%qG#MMYPs6LP&2vY_YWbBC-Fwf73u&kSCzpQrJcr&mPN3cAfNXzDY^wQrYB6m?$?E3%Aee+SR_jd8GVSWFpdVSK3Pwa9LSFcA*AtYc_2zY!8n_#U7ZnFZb?z7S_ZJH*qd$tWD%w+K3R!LyjK6$f83v+-k8MQWFuQ_Uhc%?8uC7+XWzD=aM@VER*jU+DPhP_2xECD8aD^L2v9P?nhLuZhq-H9!q_hZ26wv%t^T3NW2Y2-2CJjN!wNNrjRgjrquS4S%a-q6kFc^%3p4fqQZd&V7FPFBbx955gV+YyPG$qL^qMTfsAfs==C#kFqnz4kaK^WK#-XuJF%UXBc2RFUXbyNJhRB_BaLP-XKX4dDP5kQgbMTL_XNmaraSVpuUC2dMA@UDqF#F_vW6_!=qX&MG#p=ru&WqB82Y6HH$2m-2gdv94T#EtPdkPXb!+_cM+rAs6fL!GA5j?6MjneA9*8zZ9dZ4+QP$%ezUFPAM2fcApEMTv5TxDSX&c4#n^BJcJTJheaa8*?sq+95$3eqaQ2yaUb45ax!QcZt!fK-kfXRndjH9+*_*#GCBd8Wj%TfUnX*zDu-YD!qPkt=Eu+S2j@x9JL7-TZ74R9%^B$&kjQyF&$NjNL9y7@B5*j^Ht%nesP*H8eRJn?Bf&xq9a3ent3!n_=$PU!xHc&VC?P^hk^2BhZzn8FeWTJNw8LJmya^K*QEsFUJ&$y@TTtGngbPz_=GAg4%_UcUd^#2CAf_=?Yxv$Kzsq2_fZ3YM=v8SKvvCu+qX#C4c?c^xy_89*j6cPZtHc3WUNHscuVc7=V?5N67x@va8fyT?w+DZMNW2vm*xdfL94%?euZ8@K^k3=fq^yFVgEMXUETj!kZAeft2hh2TTx-EWKVJr+F%9S&YSL*hSAx+-*BQ8Jqp2S_kK^VQUMfDHSmPUs8jnT+xf%TUt^#nR?Mk$Y7z6+_F@$d8v#_!GgZ^@vsEx&#tAud@EB7-@$P+%qSa_Yv_JRgCLAx97Ug-VWCf?sccfJ6%&bZp9Vt&KV9D%h+tFKDeM5hZ?E^BTvfyjCsP_BGzMPzvc-X6C8$#bM+WEW9UEcfd+7e=w9RQV@8*FLjJm22HT8$5ErBqe638cpQC294ACtc#-qd@ET@rL%y_yp4eEZuxx2xy+K?SaY@nbLtHAM+L-v?5fFhpHv8W$9x&y=6%+VqH#N@TgWDjj?zSvS2TZn=Hg%HQ#D_Hc69Kfs&KN@zdDpW&^wgA?!^^kRR*3kv8MNyeWUR8g#A_6Edp!H5r^CLkSR5J2EvuE^VXbTDwSez6jK!&Tn-%Nt7FwnEF$KrWzEf!K%asXL#A-EhwXG+M5X3mAYHUeEVtffd$vRbBAm4KBX+?AzPre$7w^pH7_^p$sxbvsjHB-V?dK4dJ!MT*ySmQ#2duWnX^eRKv5NJ2@kp6S&H?^R93V$T7vngM^Ac$*++dzjWKH!dYbWyUzp-#mG=xR`)
)

//TODO: Add expires date
type AuthToken struct {
	Header    map[string]interface{} `json:"header"`
	Claim     Claim                  `json:"claim"`
	Signature string                 `json:"signature"`
}

type Header struct {
	Type string `json:"typ"`
	Alg  string `json:"alg"`
}

type Claim struct {
	UserId   string `json:"userid"`
	Username string `json:"username"`
	IAT      string `json:"iat"` // Issued At Time
}

// authToken := web.NewAuthToken(token, sig)
func NewAuthToken(token *jwt.Token, signature string) *AuthToken {
	return &AuthToken{
		Header:    token.Header,
		Claim:     token.Claims.(Claim),
		Signature: signature,
	}
}

func NewClaim(userId, username string) Claim {
	issuedAt := formattedTime()
	return Claim{userId, username, issuedAt}
}

func (this Claim) Valid() error {
	return nil
}

func Signature(claim Claim) (string, *jwt.Token) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	ss, _ := token.SignedString(mySigningKey)
	return ss, token
}

func DecodeFromJSON(values map[string]interface{}, token *AuthToken) error {
	return structutil.Decode(values, token)
}

func FromSignature(r *http.Request) (string, error) {
	header, err := jwtmiddleware.FromAuthHeader(r)

	var values map[string]interface{}
	err = jsonutil.DecodeJSON([]byte(header), &values)
	if err != nil {
		return "", err
	}

	signature := values["signature"].(string)

	return signature, nil
}

func formattedTime() string {
	t := time.Now().UTC()
	return t.Format(time.RFC3339)
}

func fromFormattedTime(value string) (time.Time, error) {
	return time.Parse(time.RFC3339, value)
}

// {
// "header": {
// "typ": "JWT",
// "alg": "HS256"
// },
// "claim": {
// "userid": 1
// },
// "signature": ""
// }
