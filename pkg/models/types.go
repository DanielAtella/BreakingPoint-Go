package models

import (
    "fmt"
    "time"
)

// Common request/response types

// LoginRequest represents authentication request
type LoginRequest struct {
    Username  string `json:"username"`
    Password  string `json:"password"`
    SessionID string `json:"sessionId,omitempty"`
}

// LoginResponse represents authentication response
type LoginResponse struct {
    SessionID    string `json:"sessionId"`
    APIKey       string `json:"apiKey"`
    APIServer    string `json:"apiServer,omitempty"`
    ServerTime   string `json:"serverTime,omitempty"`
}

// TestRunRequest represents test execution request
type TestRunRequest struct {
    ModelName    string `json:"modelname"`
    Group        int    `json:"group"`
    AllowMalware bool   `json:"allowMalware,omitempty"`
}

// TestRunResponse represents test execution response
type TestRunResponse struct {
    RunID   int    `json:"runid"`
    Status  string `json:"status"`
    Message string `json:"message,omitempty"`
}

// SearchRequest represents generic search request
type SearchRequest struct {
    SearchString string `json:"searchString,omitempty"`
    Limit        string `json:"limit,omitempty"`
    Sort         string `json:"sort,omitempty"`
    SortOrder    string `json:"sortorder,omitempty"`
    Offset       string `json:"offset,omitempty"`
}

// ExportRequest represents file export request
type ExportRequest struct {
    Name        string `json:"name"`
    Attachments bool   `json:"attachments,omitempty"`
    FilePath    string `json:"filepath"`
    RunID       *int   `json:"runid,omitempty"`
}

// ImportRequest represents file import request
type ImportRequest struct {
    Name     string `json:"name"`
    Filename string `json:"filename"`
    Force    bool   `json:"force,omitempty"`
}

// RealtimeStatsRequest represents real-time statistics request
type RealtimeStatsRequest struct {
    RunID         int      `json:"runid"`
    RTSGroup      string   `json:"rtsgroup"`
    NumSeconds    int      `json:"numSeconds"`
    NumDataPoints int      `json:"numDataPoints,omitempty"`
    Aggregate     string   `json:"aggregate,omitempty"`
    Protocol      []string `json:"protocol,omitempty"`
}

// RealtimeStatsResponse represents real-time statistics response
type RealtimeStatsResponse struct {
    TestStuck bool        `json:"testStuck"`
    Time      int64       `json:"time"`
    Progress  float64     `json:"progress"`
    Values    interface{} `json:"values"`
}

// ComponentInfo represents test component information
type ComponentInfo struct {
    ID                   string      `json:"id"`
    Label                string      `json:"label"`
    Type                 string      `json:"type"`
    Active               bool        `json:"active"`
    OriginalPreset       string      `json:"originalPreset,omitempty"`
    OriginalPresetLabel  string      `json:"originalPresetLabel,omitempty"`
    ReportResults        bool        `json:"reportResults,omitempty"`
    Tags                 []Tag       `json:"tags,omitempty"`
    Timeline             Timeline    `json:"timeline,omitempty"`
}

// Tag represents component tags
type Tag struct {
    ID       string   `json:"id"`
    Type     string   `json:"type"`
    DomainID DomainID `json:"domainId"`
}

// DomainID represents domain identification
type DomainID struct {
    Name     string `json:"name"`
    External bool   `json:"external"`
    Iface    string `json:"iface"`
}

// Timeline represents component timeline
type Timeline struct {
    TimeSegments []TimeSegment `json:"timesegment"`
}

// TimeSegment represents timeline segments
type TimeSegment struct {
    Label string `json:"label"`
    Size  int    `json:"size"`
    Type  string `json:"type"`
}

// NetworkInfo represents network configuration information
type NetworkInfo struct {
    Name           string      `json:"name"`
    Label          string      `json:"label"`
    Description    string      `json:"description"`
    InterfaceCount int         `json:"interfaceCount"`
    CreatedBy      string      `json:"createdBy"`
    CreatedOn      time.Time   `json:"createdOn"`
    Revision       int         `json:"revision"`
    NetworkModel   interface{} `json:"networkModel"`
}

// TopologyInfo represents chassis topology information
type TopologyInfo struct {
    Model        string `json:"model"`
    SerialNumber string `json:"serialNumber"`
    Slots        []Slot `json:"slot"`
}

// Slot represents chassis slot information
type Slot struct {
    ID             int    `json:"id"`
    Model          string `json:"model"`
    SerialNumber   string `json:"serialNumber"`
    State          string `json:"state"`
    InterfaceCount int    `json:"interfaceCount"`
    Ports          []Port `json:"port"`
}

// Port represents port information
type Port struct {
    ID           string `json:"id"`
    Number       int    `json:"number"`
    State        string `json:"state"`
    Link         string `json:"link"`
    Speed        int    `json:"speed"`
    Media        string `json:"media"`
    Group        int    `json:"group,omitempty"`
    ReservedBy   string `json:"reservedBy,omitempty"`
    Capture      bool   `json:"capture,omitempty"`
    Owner        string `json:"owner,omitempty"`
}

// StrikeInfo represents security strike information
type StrikeInfo struct {
    ID            string      `json:"id"`
    Name          string      `json:"name"`
    Path          string      `json:"path"`
    Protocol      string      `json:"protocol"`
    Category      string      `json:"category"`
    Direction     string      `json:"direction"`
    Severity      string      `json:"severity"`
    Year          string      `json:"year"`
    Variants      int         `json:"variants"`
    FileSize      string      `json:"fileSize"`
    FileExtension string      `json:"fileExtension"`
    Keywords      []Keyword   `json:"keyword"`
    References    []Reference `json:"reference"`
}

// Keyword represents strike keywords
type Keyword struct {
    Name string `json:"name"`
}

// Reference represents strike references
type Reference struct {
    Label string `json:"label"`
    Type  string `json:"type"`
    Value string `json:"value"`
}

// Error represents API error response
type Error struct {
    StatusCode int         `json:"status_code"`
    Content    interface{} `json:"content"`
    Message    string      `json:"message,omitempty"`
}

// Error implements the error interface
func (e *Error) Error() string {
    if e.Message != "" {
        return e.Message
    }
    return fmt.Sprintf("API error %d: %v", e.StatusCode, e.Content)
}

