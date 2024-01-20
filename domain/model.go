package domain

type NotificationData struct {
	Identity  string                 `json:"identity"`
	Timestamp int64                  `json:"ts"`
	Type      string                 `json:"type"`
	EvtName   string                 `json:"evtName"`
	EvtData   map[string]interface{} `json:"evtData"`
	Caller    string                 `json:"caller"`
	RequestID string                 `json:"request_id"`
}

type PushNotificationRequest struct {
	Data []MetadataPushNotif `json:"d"`
}

type MetadataEmail struct {
	Identity      string                 `json:"identity"`
	Ts            int64                  `json:"ts"`
	Type          string                 `json:"type"`
	TemplateName  string                 `json:"templateName"`
	TemplateParam map[string]interface{} `json:"templateParam"`
}

type MetadataPushNotif struct {
	Identity string                 `json:"identity"`
	Ts       int64                  `json:"ts"`
	Type     string                 `json:"type"`
	EvtName  string                 `json:"evtName"`
	EvtData  map[string]interface{} `json:"evtData"`
}

type ResponseNotif struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message"`
}
