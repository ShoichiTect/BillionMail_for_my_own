package batch_mail

import (
	"billionmail-core/internal/model/entity"
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gview"
	"strings"
)

// TemplateEngine
type TemplateEngine struct {
	view *gview.View
}

// NewTemplateEngine Create new template engine instance
func NewTemplateEngine() *TemplateEngine {
	view := gview.New()

	// register template functions
	view.BindFuncMap(g.Map{
		"get": func(m interface{}, key string) interface{} {
			// handle map
			if m, ok := m.(map[string]string); ok {
				return m[key]
			}
			// handle g.Map
			if m, ok := m.(g.Map); ok {
				return m[key]
			}
			return ""
		},
		"UnsubscribeURL": func(data ...interface{}) interface{} {
			if len(data) > 0 {
				if m, ok := data[0].(g.Map); ok {
					if url, ok := m["UnsubscribeURL"].(string); ok {
						return url
					}
				}
			}
			return ""
		},
	})

	return &TemplateEngine{view: view}
}

// RenderTemplate Render template
func (e *TemplateEngine) RenderTemplate(ctx context.Context, content string, data g.Map) (string, error) {
	return e.view.ParseContent(ctx, content, data)
}

func (e *TemplateEngine) RenderEmailTemplate(ctx context.Context, content string, contact *entity.Contact, task *entity.EmailTask, unsubscribeURL string) (string, error) {
	return e.RenderEmailTemplateWithAPI(ctx, content, contact, task, unsubscribeURL, nil)
}

// RenderEmailTemplate Render email template
func (e *TemplateEngine) RenderEmailTemplateWithAPI(ctx context.Context, content string, contact *entity.Contact, task *entity.EmailTask, unsubscribeURL string, apiAttribs map[string]interface{}) (string, error) {
	// handle escape characters
	content = strings.NewReplacer(
		"\\r\\n", "\n",
		"\\n", "\n",
		"\\r", "\n",
		"\\t", "\t",
	).Replace(content)

	// prepare subscriber data
	subscriberData := make(map[string]interface{})
	if contact != nil {
		// Custom properties are added directly to subscriberData
		if contact.Attribs != nil {
			for k, v := range contact.Attribs {
				subscriberData[k] = v
			}
		}
		subscriberData["Email"] = contact.Email
		subscriberData["Active"] = contact.Active
		subscriberData["Status"] = contact.Status

	}

	// prepare api data
	apiData := make(map[string]interface{})
	if contact != nil {
		// Custom properties are added directly to subscriberData
		if apiAttribs != nil {
			for k, v := range apiAttribs {
				apiData[k] = v
			}
		}

	}

	// prepare task data
	taskData := g.Map{
		"Id":             0,
		"TaskName":       "",
		"Addresser":      "",
		"Subject":        "",
		"FullName":       "",
		"RecipientCount": 0,
		"TaskProcess":    0,
		"Pause":          0,
		"TemplateId":     0,
		"IsRecord":       0,
		"Unsubscribe":    0,
		"Threads":        0,
		"Etypes":         "",
		"TrackOpen":      0,
		"TrackClick":     0,
		"StartTime":      0,
		"CreateTime":     0,
		"UpdateTime":     0,
		"Remark":         "",
		"Active":         0,
	}
	if task != nil {
		taskData["Id"] = task.Id
		taskData["TaskName"] = task.TaskName
		taskData["Addresser"] = task.Addresser
		taskData["Subject"] = task.Subject
		taskData["FullName"] = task.FullName
		taskData["RecipientCount"] = task.RecipientCount
		taskData["TaskProcess"] = task.TaskProcess
		taskData["Pause"] = task.Pause
		taskData["TemplateId"] = task.TemplateId
		taskData["IsRecord"] = task.IsRecord
		taskData["Unsubscribe"] = task.Unsubscribe
		taskData["Threads"] = task.Threads
		taskData["Etypes"] = task.Etypes
		taskData["TrackOpen"] = task.TrackOpen
		taskData["TrackClick"] = task.TrackClick
		taskData["StartTime"] = task.StartTime
		taskData["CreateTime"] = task.CreateTime
		taskData["UpdateTime"] = task.UpdateTime
		taskData["Remark"] = task.Remark
		taskData["Active"] = task.Active
	}

	// prepare template data
	templateData := g.Map{
		"Subscriber":     subscriberData,
		"Task":           taskData,
		"UnsubscribeURL": unsubscribeURL,
		"API":            apiAttribs,
	}

	// use template engine to render
	result, err := e.view.ParseContent(ctx, content, templateData)
	if err != nil {
		return "", err
	}

	return result, nil
}

// global template engine instance
var defaultEngine = NewTemplateEngine()

// GetTemplateEngine Get default template engine instance
func GetTemplateEngine() *TemplateEngine {
	return defaultEngine
}
