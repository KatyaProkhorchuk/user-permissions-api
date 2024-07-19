package types

type Response struct {
	Acsess []string `json:"access"`
}

// список записей, к которым у пользователя есть доступ
type ArchiveManager struct {
	Records []string `json:"records"`
}

// список агентов, к которым у пользователя есть доступ
type TaskManager struct {
	Agent []string `json:"agent"`
}

// доступные сервисы для пользователя: ArchiveManager и TaskManager.
type AccessServices struct {
	Archive ArchiveManager `json:"archive"`
	Task    TaskManager    `json:"task"`
}

type User struct {
	Name   string         `json:"name"`
	Access AccessServices `json:"user_access"`
}

// запрос от клиента
type Request struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

func (requestAccess *AccessServices) IsAcess(targetAccess *AccessServices) bool {
	targetAgentMap := make(map[string]struct{}, len(requestAccess.Task.Agent))
	for _, targetAgent := range targetAccess.Task.Agent {
		targetAgentMap[targetAgent] = struct{}{}
	}
	for _, requestAgent := range requestAccess.Task.Agent {
		if _, exists := targetAgentMap[requestAgent]; !exists {
			return false
		}
	}
	return true
}
