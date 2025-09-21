package handler

import (
	"fmt"
	"github.com/tomiok/toss/pkg/web"
	"net/http"
)

type Handler struct {
}

func New() Handler {
	return Handler{}
}

func (h Handler) UploadView(w http.ResponseWriter, r *http.Request) error {

	return web.TemplateRender(w, "upload.page.tmpl", nil, false)
}

func (h Handler) PreviewOrDownload(w http.ResponseWriter, r *http.Request) error {
	return web.TemplateRender(w, "preview.page.tmpl", nil, false)
}

func (h Handler) Upload(w http.ResponseWriter, r *http.Request) error {
	//b, err := io.ReadAll(r.Body)
	//if err != nil {
	//	return err
	//}

	//fmt.Println(string(b))

	url := fmt.Sprintf("%s/once/%s", "localhost", web.GenerateHash())
	html := fmt.Sprintf(`
    <div style="background: #f0fdf4; border: 1px solid #bbf7d0; border-radius: 8px; padding: 20px; margin-top: 20px;">
        <div style="color: #065f46; font-weight: 600; margin-bottom: 10px;">✅ File uploaded successfully!</div>
        <div class="result-url" id="resultUrl">%s</div>
        <button class="copy-btn" onclick="copyToClipboard(this)">Copy Link</button>
    </div>
`, url)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
	return nil
}
