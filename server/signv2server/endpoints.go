package signv2server

import (
    "database/sql"
    "encoding/json"
	"io/ioutil"
	"path/filepath"
    "errors"
    "erupe-ce/server/channelserver"
    "net/http"
    "strings"

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
    Banners  []LauncherBanner  `json:"banners"`
    Messages []LauncherMessage `json:"messages"`
    Links    []LauncherLink    `json:"links"`
    Background     string `json:"background,omitempty"`
    Cog            string `json:"cog,omitempty"`
    Capcom         string `json:"capcom,omitempty"`
    LauncherHeader string `json:"launcher_header,omitempty"`
}

// User holds the launcher‐side auth token, its numeric ID, and rights.
type User struct {
    TokenId uint32 `json:"tokenId"` // numeric ID of the login_token record
    Token   string `json:"token"`   // opaque token string
    Rights  uint32 `json:"rights"`
}

// Character is a single character entry returned to the launcher.
type Character struct {
    ID        uint32 `json:"id"`
    Name      string `json:"name"`
    IsFemale  bool   `json:"isFemale" db:"is_female"`
    Weapon    uint32 `json:"weapon" db:"weapon_type"`
    HR        uint32 `json:"hr" db:"hrp"`
    GR        uint32 `json:"gr"`
    LastLogin int32  `json:"lastLogin" db:"last_login"`
}

// AuthData is the JSON sent to POST /login.
type AuthData struct {
    CurrentTS     uint32      `json:"currentTs"`
    ExpiryTS      uint32      `json:"expiryTs"`
    EntranceCount uint32      `json:"entranceCount"`
    Notices       []string    `json:"notices"`
    User          User        `json:"user"`
    Characters    []Character `json:"characters"`
    PatchServer   string      `json:"patchServer"`
}

// newAuthData builds the AuthData, including numeric tokenId.
func (s *Server) newAuthData(userID, userRights, tokenID uint32, userToken string, characters []Character) AuthData {
    resp := AuthData{
        CurrentTS:     uint32(channelserver.TimeAdjusted().Unix()),
        ExpiryTS:      uint32(s.getReturnExpiry(userID).Unix()),
        EntranceCount: 1,
        User: User{
            TokenId: tokenID,
            Token:   userToken,
            Rights:  userRights,
        },
        Characters:  characters,
        PatchServer: s.erupeConfig.SignV2.PatchServer,
        Notices:     []string{},
    }

    // Login notices
    if !s.erupeConfig.HideLoginNotice {
        notice := strings.Join(s.erupeConfig.LoginNotices[:], "<PAGE>")
        resp.Notices = append(resp.Notices, notice)
    }

    return resp
}

// Launcher handles GET /launcher
func (s *Server) Launcher(w http.ResponseWriter, r *http.Request) {
    // look in PATCHSERVER/launcher.json under the project root
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

// Login handles POST /login
func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // Decode incoming JSON
    var reqData struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
        s.logger.Error("JSON decode error", zap.Error(err))
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Invalid data received"))
        return
    }

    // Lookup user
    var userID uint32
    var storedPass string
    var rights uint32
    err := s.db.QueryRowContext(ctx,
        "SELECT id, password, rights FROM users WHERE username = $1",
        reqData.Username,
    ).Scan(&userID, &storedPass, &rights)
    if err == sql.ErrNoRows {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Username does not exist"))
        return
    } else if err != nil {
        s.logger.Warn("SQL query error", zap.Error(err))
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    // Check password
    if bcrypt.CompareHashAndPassword([]byte(storedPass), []byte(reqData.Password)) != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Your password is incorrect"))
        return
    }

    // Create opaque token
    userToken, err := s.createLoginToken(ctx, userID)
    if err != nil {
        s.logger.Warn("Error registering login token", zap.Error(err))
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    // Fetch numeric token ID
    tokenID := uint32(0)

    // Fetch characters
    chars, err := s.getCharactersForUser(ctx, userID)
    if err != nil {
        s.logger.Warn("Error getting characters from DB", zap.Error(err))
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    if chars == nil {
        chars = []Character{}
    }

    // Build and send AuthData
    respData := s.newAuthData(userID, rights, tokenID, userToken, chars)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(respData)
}

// Register handles POST /register
func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    var reqData struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
        s.logger.Error("JSON decode error", zap.Error(err))
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Invalid data received"))
        return
    }
    if reqData.Username == "" || reqData.Password == "" {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Username and password must not be empty"))
        return
    }

    s.logger.Info("Creating account", zap.String("username", reqData.Username))
    userID, rights, err := s.createNewUser(ctx, reqData.Username, reqData.Password)
    if err != nil {
        var pqErr *pq.Error
        if errors.As(err, &pqErr) && pqErr.Constraint == "users_username_key" {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte("User already exists"))
            return
        }
        s.logger.Error("Error checking user", zap.Error(err), zap.String("username", reqData.Username))
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    // Create token and fetch its ID
    userToken, err := s.createLoginToken(ctx, userID)
    if err != nil {
        s.logger.Error("Error registering login token", zap.Error(err))
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    var tokenID uint32
    if err := s.db.QueryRowContext(ctx,
        "SELECT id FROM login_tokens WHERE token = $1", userToken,
    ).Scan(&tokenID); err != nil {
        s.logger.Warn("Error fetching token ID", zap.Error(err))
        tokenID = 0
    }

    respData := s.newAuthData(userID, rights, tokenID, userToken, []Character{})
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(respData)
}

// CreateCharacter handles POST /createCharacter
func (s *Server) CreateCharacter(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    var reqData struct {
        Token string `json:"token"`
    }
    if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
        s.logger.Error("JSON decode error", zap.Error(err))
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Invalid data received"))
        return
    }

    userID, err := s.userIDFromToken(ctx, reqData.Token)
    if err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

    character, err := s.createCharacter(ctx, userID)
    if err != nil {
        s.logger.Error("Failed to create character", zap.Error(err), zap.String("token", reqData.Token))
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(character)
}

// DeleteCharacter handles POST /deleteCharacter
func (s *Server) DeleteCharacter(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    var reqData struct {
        Token  string `json:"token"`
        CharID uint32 `json:"charId"`
    }
    if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
        s.logger.Error("JSON decode error", zap.Error(err))
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Invalid data received"))
        return
    }

    userID, err := s.userIDFromToken(ctx, reqData.Token)
    if err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

    if err := s.deleteCharacter(ctx, userID, reqData.CharID); err != nil {
        s.logger.Error("Failed to delete character", zap.Error(err),
            zap.String("token", reqData.Token), zap.Uint32("charID", reqData.CharID))
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(struct{}{})
}

// ExportSave handles POST /exportSave
func (s *Server) ExportSave(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    var reqData struct {
        Token  string `json:"token"`
        CharID uint32 `json:"charId"`
    }
    if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
        s.logger.Error("JSON decode error", zap.Error(err))
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Invalid data received"))
        return
    }

    userID, err := s.userIDFromToken(ctx, reqData.Token)
    if err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

    save, err := s.exportSave(ctx, userID, reqData.CharID)
    if err != nil {
        s.logger.Error("Failed to export save", zap.Error(err),
            zap.String("token", reqData.Token), zap.Uint32("charID", reqData.CharID))
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(save)
}