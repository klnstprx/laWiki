package dto

// VersionDTO represents the data structure received from the Version service.
type VersionDTO struct {
    ID          string `json:"id"`
    Content     string `json:"content"`
    EntryID     string `json:"entryId"`
    CreatedAt   string `json:"createdAt"`
    UpdatedAt   string `json:"updatedAt"`
    // Add other necessary fields if required
}