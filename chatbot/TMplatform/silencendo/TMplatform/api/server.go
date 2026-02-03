package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"silencendo/cli"
)

// ChatRequest is the JSON body for POST /api/chat.
type ChatRequest struct {
	Message            string `json:"message"`
	ActiveProjectID    string `json:"active_project_id"`
	ActiveProjectTitle string `json:"active_project_title"`
}

// ChatResponse is the JSON response from POST /api/chat.
type ChatResponse struct {
	Reply              string `json:"reply"`
	Exit               bool   `json:"exit"`
	ActiveProjectID    string `json:"active_project_id"`
	ActiveProjectTitle string `json:"active_project_title"`
}

// Server holds the dispatcher and state for the chat API.
type Server struct {
	Dispatcher *cli.Dispatcher
	State      *cli.SessionState
}

// NewServer returns a new API server.
func NewServer(dispatcher *cli.Dispatcher, state *cli.SessionState) *Server {
	return &Server{Dispatcher: dispatcher, State: state}
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(`<!DOCTYPE html>
<html>
<head><meta charset="utf-8"><title>Chatbot API</title></head>
<body style="font-family: system-ui; max-width: 42rem; margin: 2rem auto; padding: 0 1rem;">
  <h1>Chatbot API</h1>
  <p>The API is running.</p>
  <p><strong><a href="/chat">Open Chatbot</a></strong> – use the chatbot in your browser.</p>
  <ul>
    <li><a href="/api/health">GET /api/health</a> – health check</li>
    <li><strong>POST /api/chat</strong> – send messages (JSON body: message, active_project_id, active_project_title)</li>
  </ul>
  <h2 style="font-size: 1.1rem; margin-top: 1.5rem;">Commands to run the code (click to copy)</h2>
  <p><strong>1. Run the backend</strong> (from TMplatform workspace root):</p>
  <pre style="background:#f4f4f4; padding: 1rem; overflow-x: auto; margin: 0 0 1rem 0;"><code class="cmd" data-cmd="./run_chatbot_http.sh" title="Click to copy">./run_chatbot_http.sh</code></pre>
  <p>Or manually:</p>
  <pre style="background:#f4f4f4; padding: 1rem; overflow-x: auto; margin: 0 0 1rem 0;"><code class="cmd" data-cmd="cd silencendo/TMplatform" title="Click to copy">cd silencendo/TMplatform</code>
<code class="cmd" data-cmd="export CHATBOT_HTTP=1" title="Click to copy">export CHATBOT_HTTP=1</code>
<code class="cmd" data-cmd="export CHATBOT_HTTP_PORT=8080" title="Click to copy">export CHATBOT_HTTP_PORT=8080</code>
<code class="cmd" data-cmd="go run ." title="Click to copy">go run .</code></pre>
  <p><strong>2. Run the frontend</strong> (in another terminal, from chatbot folder):</p>
  <pre style="background:#f4f4f4; padding: 1rem; overflow-x: auto; margin: 0;"><code class="cmd" data-cmd="cd .." title="Click to copy">cd ..</code>   # from TMplatform to chatbot
<code class="cmd" data-cmd="export CHATBOT_API_URL=http://localhost:8080" title="Click to copy">export CHATBOT_API_URL=http://localhost:8080</code>
<code class="cmd" data-cmd="npm run dev" title="Click to copy">npm run dev</code></pre>
  <p id="copy-msg" style="color: green; font-size: 0.9rem; margin-top: 0.25rem; min-height: 1.25rem;"></p>
  <script>
    document.querySelectorAll('.cmd').forEach(function(el) {
      el.style.cursor = 'pointer';
      el.style.textDecoration = 'underline';
      el.style.textDecorationStyle = 'dotted';
      el.addEventListener('click', function() {
        var cmd = this.getAttribute('data-cmd');
        navigator.clipboard.writeText(cmd).then(function() {
          var msg = document.getElementById('copy-msg');
          msg.textContent = 'Copied: ' + cmd;
          setTimeout(function() { msg.textContent = ''; }, 2000);
        });
      });
    });
  </script>
</body>
</html>`))
}

func (s *Server) handleChatPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/chat" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(`<!DOCTYPE html>
<html>
<head><meta charset="utf-8"><title>Chatbot</title>
<style>
  * { box-sizing: border-box; }
  body { font-family: system-ui, sans-serif; max-width: 36rem; margin: 0 auto; padding: 1rem; background: #1a1a1a; color: #e0e0e0; min-height: 100vh; display: flex; flex-direction: column; }
  h1 { font-size: 1.25rem; margin: 0 0 1rem 0; }
  #log { flex: 1; overflow-y: auto; border: 1px solid #333; border-radius: 8px; padding: 1rem; margin-bottom: 1rem; min-height: 12rem; background: #252525; font-size: 0.9rem; white-space: pre-wrap; word-break: break-word; }
  .user { color: #7dd3fc; margin-bottom: 0.5rem; }
  .bot { color: #86efac; margin-bottom: 0.75rem; }
  form { display: flex; gap: 0.5rem; }
  #msg { flex: 1; padding: 0.6rem 0.75rem; border: 1px solid #444; border-radius: 6px; background: #252525; color: #e0e0e0; font-size: 1rem; }
  #msg:focus { outline: none; border-color: #0ea5e9; }
  button { padding: 0.6rem 1rem; border: none; border-radius: 6px; background: #0ea5e9; color: #fff; font-weight: 600; cursor: pointer; }
  button:hover { background: #0284c7; }
  button:disabled { opacity: 0.6; cursor: not-allowed; }
  a { color: #7dd3fc; }
</style>
</head>
<body>
  <h1>Chatbot</h1>
  <p style="margin: 0 0 0.5rem 0; font-size: 0.9rem;">Projects &amp; tasks – try <code>/project list</code> or &quot;create project My App&quot;</p>
  <div id="log"></div>
  <form id="f">
    <input id="msg" type="text" placeholder="Type a message..." autocomplete="off" />
    <button type="submit">Send</button>
  </form>
  <script>
    var log = document.getElementById('log');
    var form = document.getElementById('f');
    var input = document.getElementById('msg');
    var session = { active_project_id: '', active_project_title: '' };
    function append(className, text) {
      var div = document.createElement('div');
      div.className = className;
      div.textContent = text;
      log.appendChild(div);
      log.scrollTop = log.scrollHeight;
    }
    form.addEventListener('submit', function(e) {
      e.preventDefault();
      var message = input.value.trim();
      if (!message) return;
      input.value = '';
      append('user', '> ' + message);
      var btn = form.querySelector('button');
      btn.disabled = true;
      fetch('/api/chat', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          message: message,
          active_project_id: session.active_project_id || '',
          active_project_title: session.active_project_title || ''
        })
      }).then(function(r) {
        if (!r.ok) throw new Error(r.statusText);
        return r.json();
      }).then(function(data) {
        session.active_project_id = data.active_project_id || '';
        session.active_project_title = data.active_project_title || '';
        if (data.reply) append('bot', data.reply);
        if (data.exit) append('bot', 'Goodbye!');
      }).catch(function(err) {
        append('bot', 'Error: ' + err.message);
      }).finally(function() { btn.disabled = false; });
    });
  </script>
</body>
</html>`))
}

// corsMiddleware adds CORS headers so the frontend (different port) can call the API.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Handler returns an http.Handler that serves the chat API.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleRoot)
	mux.HandleFunc("/chat", s.handleChatPage)
	mux.HandleFunc("/api/chat", s.handleChat)
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	return corsMiddleware(mux)
}

func (s *Server) handleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	msg := strings.TrimSpace(req.Message)
	if msg == "" {
		http.Error(w, "message is required", http.StatusBadRequest)
		return
	}

	// Restore session from request (do not persist to disk for API calls)
	if s.State != nil {
		s.State.ActiveProjectID = strings.TrimSpace(req.ActiveProjectID)
		s.State.ActiveProjectTitle = strings.TrimSpace(req.ActiveProjectTitle)
	}

	reply, exit, err := s.Dispatcher.DispatchAndCapture(context.Background(), req.Message)
	if err != nil {
		log.Printf("chat dispatch error: %v", err)
		reply = "Something went wrong. Please try again."
	}

	activeID := ""
	activeTitle := ""
	if s.State != nil {
		activeID = strings.TrimSpace(s.State.ActiveProjectID)
		activeTitle = strings.TrimSpace(s.State.ActiveProjectTitle)
	}

	resp := ChatResponse{
		Reply:              reply,
		Exit:               exit,
		ActiveProjectID:    activeID,
		ActiveProjectTitle: activeTitle,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("chat response encode error: %v", err)
	}
}
