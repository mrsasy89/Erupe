package api

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"errors"
	"erupe-ce/server/channelserver"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const (
	NotificationDefault = iota
	NotificationNew
)

// LauncherBanner is a single banner image + link for the launcher home screen.
type LauncherBanner struct {
	Src  string `json:"src"`
	Link string `json:"link"`
}

// LauncherMessage is a text notice for the launcher home screen.
type LauncherMessage struct {
	Message string `json:"message"`
	Date    int64  `json:"date"`
	Link    string `json:"link"`
	Kind    int    `json:"kind"`
}

// LauncherLink is a simple named link on the launcher home screen.
type LauncherLink struct {
	Name string `json:"name"`
	Link string `json:"link"`
	Icon string `json:"icon"`
}

// LauncherResponse is the JSON sent to GET /launcher.
type LauncherResponse struct {
	Banners        []LauncherBanner  `json:"banners"`
	Messages       []LauncherMessage `json:"messages"`
	Links          []LauncherLink    `json:"links"`
	Background     string            `json:"background,omitempty"`
	Cog            string            `json:"cog,omitempty"`
	Capcom         string            `json:"capcom,omitempty"`
	LauncherHeader string            `json:"launcher_header,omitempty"`
}

type User struct {
	TokenID uint32 `json:"tokenId"`
	Token   string `json:"token"`
	Rights  uint32 `json:"rights"`
}

type Character struct {
	ID        uint32 `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	IsFemale  bool   `json:"isFemale" db:"is_female"`
	Weapon    uint32 `json:"weapon" db:"weapon_type"`
	HR        uint32 `json:"hr" db:"hrp"`
	GR        uint32 `json:"gr" db:"gr"`
	LastLogin int32  `json:"lastLogin" db:"last_login"`
}

type MezFes struct {
	ID           uint32   `json:"id"`
	Start        uint32   `json:"start"`
	End          uint32   `json:"end"`
	SoloTickets  uint32   `json:"soloTickets"`
	GroupTickets uint32   `json:"groupTickets"`
	Stalls       []uint32 `json:"stalls"`
}

type Member struct {
	CID  uint32 `json:"cid"`
	ID   uint32 `json:"id"`
	Name string `json:"name"`
}

type AuthData struct {
	CurrentTS     uint32      `json:"currentTs"`
	ExpiryTS      uint32      `json:"expiryTs"`
	EntranceCount uint32      `json:"entranceCount"`
	Notices       []string    `json:"notices"`
	User          User        `json:"user"`
	Characters    []Character `json:"characters"`
	MezFes        *MezFes     `json:"mezFes"`
	Friends       []Member    `json:"friends"`
	PatchServer   string      `json:"patchServer"`
}

type ExportData struct {
	Character map[string]interface{} `json:"character"`
}

func (s *APIServer) newAuthData(userID uint32, userRights uint32, userTokenID uint32, userToken string, characters []Character) AuthData {
	resp := AuthData{
		CurrentTS:     uint32(channelserver.TimeAdjusted().Unix()),
		ExpiryTS:      uint32(s.getReturnExpiry(userID).Unix()),
		EntranceCount: 1,
		User: User{
			Rights:  userRights,
			TokenID: userTokenID,
			Token:   userToken,
		},
		Characters:  characters,
		PatchServer: s.erupeConfig.API.PatchServer,
		Notices:     []string{},
	}

	// 9.2: DevModeOptions instead of DebugOptions
	if s.erupeConfig.DevModeOptions.MaxLauncherHR {
		for i := range resp.Characters {
			resp.Characters[i].HR = 7
		}
	}

	// MezFes event copied from sign server logic
	if s.erupeConfig.DevModeOptions.MezFesEvent {
		stalls := []uint32{10, 3, 6, 9, 4, 8, 5, 7}
		if s.erupeConfig.DevModeOptions.MezFesAlt {
			stalls[4] = 2 // Tokotoko
		}
		resp.MezFes = &MezFes{
			ID:           uint32(channelserver.TimeWeekStart().Unix()),
			Start:        uint32(channelserver.TimeWeekStart().Unix()),
			End:          uint32(channelserver.TimeWeekNext().Unix()),
			SoloTickets:  20,
			GroupTickets: 20,
			Stalls:       stalls,
		}
	} else {
		resp.MezFes = &MezFes{ID: 0, Start: 0, End: 0}
	}

	if !s.erupeConfig.HideLoginNotice {
		resp.Notices = append(resp.Notices, strings.Join(s.erupeConfig.LoginNotices[:], "<PAGE>"))
	}

	friends, err := s.getFriendsForCharacters(characters)
	if err != nil {
		s.logger.Warn("Error fetching friends", zap.Error(err))
	}
	if len(friends) > 255 {
		friends = friends[:255]
	}
	resp.Friends = friends

	return resp
}

func (s *APIServer) Launcher(w http.ResponseWriter, r *http.Request) {
	cfgPath := filepath.Join("PATCHSERVER", "launcher.json")

	raw, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		http.Error(w, "failed to load launcher config", http.StatusInternalServerError)
		return
	}

	var respData LauncherResponse
	if err := json.Unmarshal(raw, &respData); err != nil {
		http.Error(w, "invalid launcher config", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respData)
}

func (s *APIServer) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var reqData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		s.logger.Error("JSON decode error", zap.Error(err))
		w.WriteHeader(400)
		return
	}
	var (
		userID     uint32
		userRights uint32
		password   string
	)
	err := s.db.QueryRow("SELECT id, password, rights FROM users WHERE username = $1", reqData.Username).Scan(&userID, &password, &userRights)
	if err == sql.ErrNoRows {
		w.WriteHeader(400)
		w.Write([]byte("username-error"))
		return
	} else if err != nil {
		s.logger.Warn("SQL query error", zap.Error(err))
		w.WriteHeader(500)
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(password), []byte(reqData.Password)) != nil {
		w.WriteHeader(400)
		w.Write([]byte("password-error"))
		return
	}

	userTokenID, userToken, err := s.createLoginToken(ctx, userID)
	if err != nil {
		s.logger.Warn("Error registering login token", zap.Error(err))
		w.WriteHeader(500)
		return
	}
	characters, err := s.getCharactersForUser(ctx, userID)
	if err != nil {
		s.logger.Warn("Error getting characters from DB", zap.Error(err))
		w.WriteHeader(500)
		return
	}
	if characters == nil {
		characters = []Character{}
	}
	respData := s.newAuthData(userID, userRights, userTokenID, userToken, characters)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respData)
}

func (s *APIServer) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var reqData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		s.logger.Error("JSON decode error", zap.Error(err))
		w.WriteHeader(400)
		return
	}
	if reqData.Username == "" || reqData.Password == "" {
		w.WriteHeader(400)
		return
	}
	s.logger.Info("Creating account", zap.String("username", reqData.Username))
	userID, userRights, err := s.createNewUser(ctx, reqData.Username, reqData.Password)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Constraint == "users_username_key" {
			w.WriteHeader(400)
			w.Write([]byte("username-exists-error"))
			return
		}
		s.logger.Error("Error checking user", zap.Error(err), zap.String("username", reqData.Username))
		w.WriteHeader(500)
		return
	}

	userTokenID, userToken, err := s.createLoginToken(ctx, userID)
	if err != nil {
		s.logger.Error("Error registering login token", zap.Error(err))
		w.WriteHeader(500)
		return
	}
	respData := s.newAuthData(userID, userRights, userTokenID, userToken, []Character{})
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respData)
}

func (s *APIServer) CreateCharacter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var reqData struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		s.logger.Error("JSON decode error", zap.Error(err))
		w.WriteHeader(400)
		return
	}

	userID, err := s.userIDFromToken(ctx, reqData.Token)
	if err != nil {
		w.WriteHeader(401)
		return
	}
	character, err := s.createCharacter(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to create character", zap.Error(err), zap.String("token", reqData.Token))
		w.WriteHeader(500)
		return
	}
	if s.erupeConfig.DevModeOptions.MaxLauncherHR {
		character.HR = 7
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(character)
}

func (s *APIServer) DeleteCharacter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var reqData struct {
		Token  string `json:"token"`
		CharID uint32 `json:"charId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		s.logger.Error("JSON decode error", zap.Error(err))
		w.WriteHeader(400)
		return
	}
	userID, err := s.userIDFromToken(ctx, reqData.Token)
	if err != nil {
		w.WriteHeader(401)
		return
	}
	if err := s.deleteCharacter(ctx, userID, reqData.CharID); err != nil {
		s.logger.Error("Failed to delete character", zap.Error(err), zap.String("token", reqData.Token), zap.Uint32("charID", reqData.CharID))
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct{}{})
}

func (s *APIServer) ExportSave(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var reqData struct {
		Token  string `json:"token"`
		CharID uint32 `json:"charId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		s.logger.Error("JSON decode error", zap.Error(err))
		w.WriteHeader(400)
		return
	}
	userID, err := s.userIDFromToken(ctx, reqData.Token)
	if err != nil {
		w.WriteHeader(401)
		return
	}
	character, err := s.exportSave(ctx, userID, reqData.CharID)
	if err != nil {
		s.logger.Error("Failed to export save", zap.Error(err), zap.String("token", reqData.Token), zap.Uint32("charID", reqData.CharID))
		w.WriteHeader(500)
		return
	}
	save := ExportData{
		Character: character,
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(save)
}

func (s *APIServer) ScreenShotGet(w http.ResponseWriter, r *http.Request) {
	token := mux.Vars(r)["id"]
	var tokenPattern = regexp.MustCompile(`[A-Za-z0-9]+`)
	if !tokenPattern.MatchString(token) || token == "" {
		http.Error(w, "Not Valid Token", http.StatusBadRequest)
		return
	}

	safePath := "./screenshots"
	path := filepath.Join(safePath, fmt.Sprintf("%s.jpg", token))
	result, err := verifyPath(path, safePath)
	if err != nil {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	file, err := os.Open(result)
	if err != nil {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "image/jpeg")
	if _, err := io.Copy(w, file); err != nil {
		http.Error(w, "Unable to send image", http.StatusInternalServerError)
		return
	}
}

func (s *APIServer) ScreenShot(w http.ResponseWriter, r *http.Request) {
	type Result struct {
		XMLName xml.Name `xml:"result"`
		Code    string   `xml:"code"`
	}
	w.Header().Set("Content-Type", "text/xml")
	result := Result{Code: "200"}

	if r.Method != http.MethodPost {
		result.Code = "405"
		writeXML(w, result)
		return
	}

	file, _, err := r.FormFile("img")
	if err != nil {
		result.Code = "400"
		writeXML(w, result)
		return
	}
	tokenPattern := regexp.MustCompile(`[A-Za-z0-9]+`)
	token := r.FormValue("token")
	if !tokenPattern.MatchString(token) || token == "" {
		result.Code = "401"
		writeXML(w, result)
		return
	}

	img, _, err := image.Decode(file)
	if err != nil {
		result.Code = "400"
		writeXML(w, result)
		return
	}

	safePath := "./screenshots"
	path := filepath.Join(safePath, fmt.Sprintf("%s.jpg", token))
	verified, err := verifyPath(path, safePath)
	if err != nil {
		result.Code = "500"
		writeXML(w, result)
		return
	}

	if _, err = os.Stat(safePath); err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(safePath, os.ModePerm); err != nil {
				s.logger.Error("Error writing screenshot, could not create folder")
				result.Code = "500"
				writeXML(w, result)
				return
			}
		} else {
			s.logger.Error("Error writing screenshot")
			result.Code = "500"
			writeXML(w, result)
			return
		}
	}

	outputFile, err := os.Create(verified)
	if err != nil {
		result.Code = "500"
		writeXML(w, result)
		return
	}
	defer outputFile.Close()

	if err = jpeg.Encode(outputFile, img, &jpeg.Options{Quality: 80}); err != nil {
		s.logger.Error("Error writing screenshot, could not write file", zap.Error(err))
		result.Code = "500"
		writeXML(w, result)
		return
	}

	writeXML(w, result)
}

func writeXML(w http.ResponseWriter, v interface{}) {
	xmlData, err := xml.Marshal(v)
	if err != nil {
		http.Error(w, "Unable to marshal XML", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(xmlData)
}
