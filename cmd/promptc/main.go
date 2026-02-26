package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/andesdevroot/promptc/pkg/core"
	"github.com/andesdevroot/promptc/pkg/sdk"
	"github.com/gorilla/websocket"
)

// --- INFRAESTRUCTURA DINÁMICA ---
var (
	configPath  string
	metricsPath string
	auditPath   string
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[FATAL] Error obteniendo directorio home: %v\n", err)
		home = "."
	}
	baseDir := filepath.Join(home, ".promptc")
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "[FATAL] Error creando directorio %s: %v\n", baseDir, err)
	}
	configPath = filepath.Join(baseDir, "templates.json")
	metricsPath = filepath.Join(baseDir, "metrics.json")
	auditPath = filepath.Join(baseDir, "audit.log")
}

// --- TIPOS JSON-RPC ---
type JSONRPCMessage struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

type Template struct {
	Description string `json:"description"`
	Content     string `json:"content"`
}

// --- SISTEMA DE AUDITORÍA ---
type AuditEvent struct {
	Timestamp string `json:"ts"`
	Type      string `json:"type"`
	Action    string `json:"action"`
	Actor     string `json:"actor"`
	Resource  string `json:"resource,omitempty"`
	Result    string `json:"result"`
	LatencyMs int64  `json:"latency_ms,omitempty"`
	Detail    string `json:"detail,omitempty"`
}

func auditLog(evt AuditEvent) {
	evt.Timestamp = time.Now().Format("2006-01-02T15:04:05.000")
	line := fmt.Sprintf("[%s] %-10s %-18s actor=%-16s",
		time.Now().Format("15:04:05.000"),
		evt.Type,
		evt.Action,
		evt.Actor,
	)
	if evt.Resource != "" {
		line += fmt.Sprintf(" resource=%-28s", evt.Resource)
	}
	line += fmt.Sprintf(" result=%s", evt.Result)
	if evt.LatencyMs > 0 {
		line += fmt.Sprintf(" latency=%dms", evt.LatencyMs)
	}
	if evt.Detail != "" {
		line += fmt.Sprintf(" | %s", evt.Detail)
	}

	hub.Lock()
	hub.Logs = append(hub.Logs, line)
	for c := range hub.Clients {
		_ = c.WriteMessage(websocket.TextMessage, []byte(line))
	}
	hub.Unlock()
	fmt.Fprintf(os.Stderr, "%s\n", line)

	f, err := os.OpenFile(auditPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		jsonLine, _ := json.Marshal(evt)
		fmt.Fprintf(f, "%s\n", string(jsonLine))
		f.Close()
	}
}

func addLog(msg string) {
	entry := fmt.Sprintf("[%s] SYSTEM     %-18s actor=promptc-engine    result=INFO | %s",
		time.Now().Format("15:04:05.000"),
		"INTERNAL",
		msg,
	)
	hub.Lock()
	hub.Logs = append(hub.Logs, entry)
	for c := range hub.Clients {
		_ = c.WriteMessage(websocket.TextMessage, []byte(entry))
	}
	hub.Unlock()
	fmt.Fprintf(os.Stderr, "%s\n", entry)
}

func sendResponse(id interface{}, result interface{}) {
	resp := JSONRPCResponse{JSONRPC: "2.0", ID: id, Result: result}
	out, _ := json.Marshal(resp)
	fmt.Fprintf(os.Stdout, "%s\n", string(out))
}

// --- MÉTRICAS ---
type MetricsSnapshot struct {
	InferenceCount   int64            `json:"inference_count"`
	InferenceSuccess int64            `json:"inference_success"`
	InferenceFail    int64            `json:"inference_fail"`
	TotalLatencyMs   int64            `json:"total_latency_ms"`
	TotalTokens      int64            `json:"total_tokens"`
	GeminiCallCount  int64            `json:"gemini_call_count"`
	TemplateCalls    map[string]int64 `json:"template_calls"`
	SavedAt          time.Time        `json:"saved_at"`
}

type Metrics struct {
	sync.Mutex
	NodeOnline       bool
	LastHeartbeat    time.Time
	InferenceCount   int64
	InferenceSuccess int64
	InferenceFail    int64
	TotalLatencyMs   int64
	TotalTokens      int64
	GeminiCallCount  int64
	TemplateCalls    map[string]int64
	CurrentMode      string
}

var metrics = &Metrics{
	TemplateCalls: make(map[string]int64),
}

func loadMetrics() {
	data, err := os.ReadFile(metricsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[METRICS] No hay estado previo en %s, arrancando limpio\n", metricsPath)
		return
	}
	var snap MetricsSnapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		fmt.Fprintf(os.Stderr, "[METRICS] metrics.json corrupto, arrancando limpio: %v\n", err)
		return
	}
	atomic.StoreInt64(&metrics.InferenceCount, snap.InferenceCount)
	atomic.StoreInt64(&metrics.InferenceSuccess, snap.InferenceSuccess)
	atomic.StoreInt64(&metrics.InferenceFail, snap.InferenceFail)
	atomic.StoreInt64(&metrics.TotalLatencyMs, snap.TotalLatencyMs)
	atomic.StoreInt64(&metrics.TotalTokens, snap.TotalTokens)
	atomic.StoreInt64(&metrics.GeminiCallCount, snap.GeminiCallCount)
	metrics.Lock()
	if snap.TemplateCalls != nil {
		metrics.TemplateCalls = snap.TemplateCalls
	}
	metrics.Unlock()
}

func saveMetrics() {
	snap := MetricsSnapshot{
		InferenceCount:   atomic.LoadInt64(&metrics.InferenceCount),
		InferenceSuccess: atomic.LoadInt64(&metrics.InferenceSuccess),
		InferenceFail:    atomic.LoadInt64(&metrics.InferenceFail),
		TotalLatencyMs:   atomic.LoadInt64(&metrics.TotalLatencyMs),
		TotalTokens:      atomic.LoadInt64(&metrics.TotalTokens),
		GeminiCallCount:  atomic.LoadInt64(&metrics.GeminiCallCount),
		SavedAt:          time.Now(),
	}
	metrics.Lock()
	snap.TemplateCalls = make(map[string]int64, len(metrics.TemplateCalls))
	for k, v := range metrics.TemplateCalls {
		snap.TemplateCalls[k] = v
	}
	metrics.Unlock()

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return
	}
	tmpPath := metricsPath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return
	}
	os.Rename(tmpPath, metricsPath)
}

func startMetricsPersistence() {
	go func() {
		lastCount := atomic.LoadInt64(&metrics.InferenceCount)
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			current := atomic.LoadInt64(&metrics.InferenceCount)
			if current != lastCount {
				saveMetrics()
				lastCount = current
			}
		}
	}()
}

func recordInference(success bool, latencyMs int64, tokens int64, usedGemini bool) {
	atomic.AddInt64(&metrics.InferenceCount, 1)
	atomic.AddInt64(&metrics.TotalLatencyMs, latencyMs)
	atomic.AddInt64(&metrics.TotalTokens, tokens)
	if success {
		atomic.AddInt64(&metrics.InferenceSuccess, 1)
	} else {
		atomic.AddInt64(&metrics.InferenceFail, 1)
	}
	if usedGemini {
		atomic.AddInt64(&metrics.GeminiCallCount, 1)
	}
	if atomic.LoadInt64(&metrics.InferenceCount)%10 == 0 {
		go saveMetrics()
	}
}

func recordTemplatCall(name string) {
	metrics.Lock()
	defer metrics.Unlock()
	metrics.TemplateCalls[name]++
}

func metricsSnapshotAPI() map[string]interface{} {
	count := atomic.LoadInt64(&metrics.InferenceCount)
	success := atomic.LoadInt64(&metrics.InferenceSuccess)
	fail := atomic.LoadInt64(&metrics.InferenceFail)
	latency := atomic.LoadInt64(&metrics.TotalLatencyMs)
	tokens := atomic.LoadInt64(&metrics.TotalTokens)
	gemini := atomic.LoadInt64(&metrics.GeminiCallCount)

	var avgLatency float64
	if count > 0 {
		avgLatency = float64(latency) / float64(count)
	}
	var tps float64
	if latency > 0 {
		tps = float64(tokens) / (float64(latency) / 1000.0)
	}
	var successRatio float64
	if count > 0 {
		successRatio = float64(success) / float64(count) * 100
	}

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	metrics.Lock()
	templateRanking := make([]map[string]interface{}, 0, len(metrics.TemplateCalls))
	for name, calls := range metrics.TemplateCalls {
		templateRanking = append(templateRanking, map[string]interface{}{
			"name":  name,
			"calls": calls,
		})
	}
	nodeOnline := metrics.NodeOnline
	lastHeartbeat := metrics.LastHeartbeat.Format("15:04:05")
	currentMode := metrics.CurrentMode
	metrics.Unlock()

	return map[string]interface{}{
		"mode":             currentMode,
		"node_online":      nodeOnline,
		"last_heartbeat":   lastHeartbeat,
		"inference_count":  count,
		"success_count":    success,
		"fail_count":       fail,
		"success_ratio":    successRatio,
		"avg_latency_ms":   avgLatency,
		"token_throughput": tps,
		"total_tokens":     tokens,
		"gemini_calls":     gemini,
		"mem_alloc_mb":     float64(memStats.Alloc) / 1024 / 1024,
		"mem_sys_mb":       float64(memStats.Sys) / 1024 / 1024,
		"goroutines":       runtime.NumGoroutine(),
		"template_ranking": templateRanking,
		"uptime_since":     startTime.Format(time.RFC3339),
	}
}

// --- HUB DASHBOARD ---
var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

type DashboardHub struct {
	sync.Mutex
	Clients   map[*websocket.Conn]bool
	Logs      []string
	Templates map[string]Template
}

var hub = &DashboardHub{
	Clients:   make(map[*websocket.Conn]bool),
	Templates: make(map[string]Template),
}

// --- HEARTBEAT ---
func startHeartbeat(remoteIP string) {
	go func() {
		client := &http.Client{Timeout: 3 * time.Second}
		for {
			resp, err := client.Get(fmt.Sprintf("http://%s:11434/api/tags", remoteIP))
			metrics.Lock()
			if err == nil && resp.StatusCode == 200 {
				wasOffline := !metrics.NodeOnline
				metrics.NodeOnline = true
				metrics.LastHeartbeat = time.Now()
				resp.Body.Close()
				metrics.Unlock()
				if wasOffline {
					auditLog(AuditEvent{
						Type:     "KERNEL",
						Action:   "NODE_ONLINE",
						Actor:    "mac-mini",
						Resource: remoteIP + ":11434",
						Result:   "OK",
						Detail:   "Nodo Ollama respondió heartbeat",
					})
				}
			} else {
				wasOnline := metrics.NodeOnline
				metrics.NodeOnline = false
				metrics.Unlock()
				if wasOnline {
					auditLog(AuditEvent{
						Type:     "KERNEL",
						Action:   "NODE_OFFLINE",
						Actor:    "mac-mini",
						Resource: remoteIP + ":11434",
						Result:   "WARN",
						Detail:   "Nodo no responde — activando fallback Gemini",
					})
				}
			}
			time.Sleep(10 * time.Second)
		}
	}()
}

// --- DASHBOARD HTML (NUEVO DISEÑO v0.3.1) ---
const dashboardHTML = `<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <title>PROMPTC // CORE_DASHBOARD</title>
    <style>
        :root { --bg: #0a0a0b; --green: #00ff41; --text: #e0e0e0; --orange: #ff9f1c; --blue: #00d4ff; }
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body { background: var(--bg); color: var(--text); font-family: 'Inter', 'Segoe UI', sans-serif; height: 100vh; display: flex; flex-direction: column; overflow: hidden; }
        .top-bar { padding: 15px 25px; border-bottom: 1px solid #222; display: flex; align-items: center; justify-content: space-between; background: #0f0f11; }
        .logo { font-family: monospace; font-weight: bold; font-size: 1.2rem; color: var(--green); letter-spacing: 2px; }
        .mode-badge { padding: 4px 12px; border-radius: 4px; font-size: 0.7rem; font-weight: bold; text-transform: uppercase; letter-spacing: 1px; }
        .community { background: var(--blue); color: #000; }
        .enterprise { background: var(--green); color: #000; }
        .metrics-grid { display: grid; grid-template-columns: repeat(6, 1fr); gap: 1px; background: #222; border-bottom: 1px solid #222; flex-shrink: 0; }
        .metric-item { background: #0f0f11; padding: 15px 20px; }
        .metric-label { font-size: 0.6rem; color: #666; text-transform: uppercase; letter-spacing: 1px; margin-bottom: 5px; }
        .metric-val { font-size: 1.1rem; font-family: monospace; font-weight: 600; color: var(--text); }
        .main { display: grid; grid-template-columns: 1fr 450px; flex-grow: 1; overflow: hidden; gap: 1px; background: #222; }
        .log-section { background: var(--bg); display: flex; flex-direction: column; overflow: hidden; }
        .section-header { padding: 10px 20px; background: #0f0f11; font-size: 0.7rem; color: #555; border-bottom: 1px solid #222; display: flex; justify-content: space-between; flex-shrink: 0; }
        #logs { flex-grow: 1; overflow-y: auto; padding: 20px; font-family: 'JetBrains Mono', 'Fira Code', monospace; font-size: 0.75rem; line-height: 1.6; color: #888; }
        .entry-kernel { color: var(--green); }
        .entry-mcp { color: var(--blue); }
        .entry-inference { color: var(--orange); }
        .side-panel { background: #0f0f11; display: flex; flex-direction: column; gap: 1px; overflow: hidden; }
        .editor-container { flex-grow: 1; display: flex; flex-direction: column; background: var(--bg); overflow: hidden; }
        textarea { flex-grow: 1; background: #0a0a0b; color: var(--green); border: none; padding: 20px; font-family: monospace; outline: none; font-size: 0.8rem; resize: none; border-bottom: 1px solid #222; }
        .btn-save { padding: 15px; background: var(--green); color: #000; border: none; font-weight: bold; cursor: pointer; text-transform: uppercase; letter-spacing: 1px; transition: 0.2s; flex-shrink: 0; }
        .btn-save:hover { background: #00cc33; }
        .install-box { padding: 15px; background: #151518; color: #888; font-size: 0.65rem; border-top: 1px solid #222; flex-shrink: 0; }
        .install-code { background: #000; padding: 8px; border-radius: 4px; margin-top: 5px; color: var(--blue); cursor: pointer; display: block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
        ::-webkit-scrollbar { width: 6px; }
        ::-webkit-scrollbar-thumb { background: #333; border-radius: 3px; }
    </style>
</head>
<body>
    <div class="top-bar">
        <div class="logo">PROMPTC <span id="v">v0.3.1</span></div>
        <div id="mode-badge" class="mode-badge">CARGANDO...</div>
    </div>
    <div class="metrics-grid">
        <div class="metric-item"><div class="metric-label">Inferencias</div><div class="metric-val" id="m-inf">0</div></div>
        <div class="metric-item"><div class="metric-label">Latencia Avg</div><div class="metric-val" id="m-lat">0 ms</div></div>
        <div class="metric-item"><div class="metric-label">Tokens Total</div><div class="metric-val" id="m-tok">0</div></div>
        <div class="metric-item"><div class="metric-label">Gemini Quota</div><div class="metric-val" id="m-gem">0</div></div>
        <div class="metric-item"><div class="metric-label">Nodo Estado</div><div class="metric-val" id="m-node">--</div></div>
        <div class="metric-item"><div class="metric-label">Goroutines</div><div class="metric-val" id="m-go">0</div></div>
    </div>
    <div class="main">
        <div class="log-section">
            <div class="section-header"><span>AUDIT_LOG_STREAM</span><span id="uptime">UPTIME: --</span></div>
            <div id="logs"></div>
        </div>
        <div class="side-panel">
            <div class="editor-container">
                <div class="section-header">TEMPLATES_CONFIG (HOT-RELOAD)</div>
                <textarea id="ed" placeholder="Cargando plantillas..."></textarea>
                <button class="btn-save" onclick="save()">Sincronizar Cambios</button>
            </div>
            <div class="install-box">
                COMPARTIR INSTALADOR COMUNIDAD:
                <code class="install-code" onclick="copyInstall()">curl -sSL https://raw.githubusercontent.com/andesdevroot/promptc/main/install.sh | bash</code>
            </div>
        </div>
    </div>
    <script>
        const ws = new WebSocket('ws://' + location.host + '/ws');
        const logs = document.getElementById('logs');
        ws.onmessage = (e) => {
            const div = document.createElement('div');
            const txt = e.data;
            if(txt.includes('KERNEL')) div.className = 'entry-kernel';
            else if(txt.includes('MCP')) div.className = 'entry-mcp';
            else if(txt.includes('INFERENCE')) div.className = 'entry-inference';
            div.textContent = txt;
            logs.appendChild(div);
            logs.scrollTop = logs.scrollHeight;
        };

        function update() {
            fetch('/api/metrics').then(r => r.json()).then(d => {
                document.getElementById('m-inf').textContent = d.inference_count || 0;
                document.getElementById('m-lat').textContent = Math.round(d.avg_latency_ms || 0) + ' ms';
                document.getElementById('m-tok').textContent = d.total_tokens || 0;
                document.getElementById('m-gem').textContent = d.gemini_calls || 0;
                document.getElementById('m-go').textContent = d.goroutines || 0;
                const node = document.getElementById('m-node');
                node.textContent = d.node_online ? 'LOCAL_ONLINE' : 'CLOUD_ONLY';
                node.style.color = d.node_online ? '#00ff41' : '#00d4ff';

                const badge = document.getElementById('mode-badge');
                if (d.mode) {
                    badge.textContent = d.mode;
                    badge.className = 'mode-badge ' + d.mode.toLowerCase();
                }
                if (d.uptime_since) {
                    const start = new Date(d.uptime_since);
                    const now = new Date();
                    const diffMs = now - start;
                    const hrs = Math.floor(diffMs / 3600000);
                    const mins = Math.floor((diffMs % 3600000) / 60000);
                    document.getElementById('uptime').textContent = 'UPTIME: ' + hrs + 'h ' + mins + 'm';
                }
            }).catch(e => console.error(e));
        }
        setInterval(update, 2000);
        update();

        fetch('/api/config').then(r => r.json()).then(d => document.getElementById('ed').value = JSON.stringify(d, null, 2));
        function save() {
            fetch('/api/config', { method: 'POST', body: document.getElementById('ed').value }).then(r => r.ok && alert('CONFIGURACIÓN ACTUALIZADA'));
        }
        function copyInstall() {
            navigator.clipboard.writeText('curl -sSL https://raw.githubusercontent.com/andesdevroot/promptc/main/install.sh | bash');
            alert('Comando copiado al portapapeles');
        }
    </script>
</body>
</html>`

// --- DASHBOARD SERVER ---
func startDashboard() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, dashboardHTML)
	})

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			addLog(fmt.Sprintf("WS upgrade falló: %v", err))
			return
		}
		hub.Lock()
		hub.Clients[conn] = true
		for _, l := range hub.Logs {
			_ = conn.WriteMessage(websocket.TextMessage, []byte(l))
		}
		hub.Unlock()
		go func() {
			defer func() {
				hub.Lock()
				delete(hub.Clients, conn)
				hub.Unlock()
				conn.Close()
			}()
			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					return
				}
			}
		}()
	})

	mux.HandleFunc("/api/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(metricsSnapshotAPI())
	})

	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		metrics.Lock()
		nodeOnline := metrics.NodeOnline
		lastHeartbeat := metrics.LastHeartbeat.Format(time.RFC3339)
		tmplCount := len(hub.Templates)
		metrics.Unlock()
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":          "ok",
			"version":         "0.3.1",
			"node_online":     nodeOnline,
			"last_heartbeat":  lastHeartbeat,
			"templates_count": tmplCount,
			"inference_count": atomic.LoadInt64(&metrics.InferenceCount),
			"uptime_since":    startTime.Format(time.RFC3339),
		})
	})

	mux.HandleFunc("/api/config", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			hub.Lock()
			json.NewEncoder(w).Encode(hub.Templates)
			hub.Unlock()
		} else {
			var n map[string]Template
			if err := json.NewDecoder(r.Body).Decode(&n); err == nil {
				hub.Lock()
				hub.Templates = n
				hub.Unlock()
				data, _ := json.MarshalIndent(n, "", "  ")
				_ = os.WriteFile(configPath, data, 0644)
				auditLog(AuditEvent{
					Type:   "SYSTEM",
					Action: "HOT_RELOAD",
					Actor:  "dashboard-operator",
					Result: "OK",
					Detail: fmt.Sprintf("%d templates recargados", len(n)),
				})
			}
		}
	})

	if err := http.ListenAndServe(":8080", mux); err != nil {
		addLog(fmt.Sprintf("Port 8080 busy: %v — MCP-ONLY Mode", err))
	}
}

var startTime = time.Now()

// --- TOOL HANDLERS ---
func handleToolCall(req JSONRPCMessage, app *sdk.PromptC) {
	var call struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if err := json.Unmarshal(req.Params, &call); err != nil {
		auditLog(AuditEvent{
			Type:   "MCP",
			Action: "TOOL_PARSE_ERROR",
			Actor:  "claude-desktop",
			Result: "FAIL",
			Detail: err.Error(),
		})
		sendResponse(req.ID, map[string]interface{}{
			"content": []map[string]interface{}{
				{"type": "text", "text": fmt.Sprintf("Error parseando tool call: %v", err)},
			},
		})
		return
	}

	auditLog(AuditEvent{
		Type:     "MCP",
		Action:   "TOOL_INVOKED",
		Actor:    "claude-desktop",
		Resource: call.Name,
		Result:   "OK",
	})

	switch call.Name {
	case "get_template":
		var args struct {
			Name string `json:"template_name"`
		}
		if err := json.Unmarshal(call.Arguments, &args); err != nil {
			sendResponse(req.ID, map[string]interface{}{
				"content": []map[string]interface{}{{"type": "text", "text": "Error: argumentos inválidos"}},
			})
			return
		}

		hub.Lock()
		tmpl, ok := hub.Templates[args.Name]
		hub.Unlock()

		if !ok {
			sendResponse(req.ID, map[string]interface{}{
				"content": []map[string]interface{}{{"type": "text", "text": "Error: template no encontrado"}},
			})
			return
		}

		recordTemplatCall(args.Name)
		sendResponse(req.ID, map[string]interface{}{
			"content": []map[string]interface{}{{"type": "text", "text": tmpl.Content}},
		})

	case "optimize_prompt":
		var args struct {
			Role        string            `json:"role"`
			Context     string            `json:"context"`
			Task        string            `json:"task"`
			Template    string            `json:"template_name"`
			Constraints []string          `json:"constraints"`
			Variables   map[string]string `json:"variables"`
		}
		if err := json.Unmarshal(call.Arguments, &args); err != nil {
			sendResponse(req.ID, map[string]interface{}{
				"content": []map[string]interface{}{{"type": "text", "text": "Error: argumentos inválidos"}},
			})
			return
		}

		task := args.Task
		if args.Template != "" {
			hub.Lock()
			tmpl, ok := hub.Templates[args.Template]
			hub.Unlock()
			if ok {
				task = tmpl.Content
				recordTemplatCall(args.Template)
			}
		}

		metrics.Lock()
		nodeOnline := metrics.NodeOnline
		metrics.Unlock()

		inferenceActor := "mac-mini"
		if !nodeOnline {
			inferenceActor = "gemini-cloud"
		}

		auditLog(AuditEvent{
			Type:     "INFERENCE",
			Action:   "PIPELINE_START",
			Actor:    "promptc-engine",
			Resource: inferenceActor,
			Result:   "OK",
		})

		start := time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		res, err := app.CompileAndOptimize(ctx, core.Prompt{
			Role:        args.Role,
			Context:     args.Context,
			Task:        task,
			Constraints: args.Constraints,
			Variables:   args.Variables,
		})

		latencyMs := time.Since(start).Milliseconds()
		estimatedTokens := int64(len(res) / 4)

		if err != nil {
			recordInference(false, latencyMs, 0, !nodeOnline)
			sendResponse(req.ID, map[string]interface{}{
				"content": []map[string]interface{}{{"type": "text", "text": fmt.Sprintf("Error: %v", err)}},
			})
			return
		}

		recordInference(true, latencyMs, estimatedTokens, !nodeOnline)
		sendResponse(req.ID, map[string]interface{}{
			"content": []map[string]interface{}{{"type": "text", "text": res}},
		})

	default:
		sendResponse(req.ID, map[string]interface{}{
			"content": []map[string]interface{}{{"type": "text", "text": "Error: herramienta no registrada"}},
		})
	}
}

// --- MAIN ---
func main() {
	modeFlag := flag.String("mode", "enterprise", "Modo de ejecución: 'community' (Viral/Cloud) o 'enterprise' (Soberano/Local)")
	flag.Parse()

	metrics.Lock()
	metrics.CurrentMode = *modeFlag
	metrics.Unlock()

	loadMetrics()

	file, err := os.ReadFile(configPath)
	if err == nil {
		json.Unmarshal(file, &hub.Templates)
	}

	go startDashboard()

	remoteIP := os.Getenv("PROMPTC_MACMINI_IP")
	if remoteIP == "" {
		remoteIP = "100.90.6.101"
	}

	if *modeFlag == "enterprise" {
		startHeartbeat(remoteIP)
	} else {
		metrics.Lock()
		metrics.NodeOnline = false
		metrics.LastHeartbeat = time.Now()
		metrics.Unlock()
	}

	startMetricsPersistence()

	app, err := sdk.NewSDK(context.Background(), os.Getenv("GEMINI_API_KEY"), remoteIP)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[SDK_ERROR] %v\n", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sigChan
		saveMetrics()
		os.Exit(0)
	}()

	auditLog(AuditEvent{
		Type:   "KERNEL",
		Action: "BOOT",
		Actor:  "promptc-engine",
		Result: "OK",
		Detail: fmt.Sprintf("PROMPTC v0.3.1 [%s] iniciado", *modeFlag),
	})

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	for scanner.Scan() {
		var req JSONRPCMessage
		if err := json.Unmarshal(scanner.Bytes(), &req); err != nil {
			continue
		}

		switch req.Method {
		case "initialize":
			sendResponse(req.ID, map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"serverInfo":      map[string]string{"name": "PROMPTC", "version": "0.3.1"},
				"capabilities":    map[string]interface{}{"tools": map[string]interface{}{}},
			})

		case "notifications/initialized":
			// Handshake confirmado

		case "tools/list":
			sendResponse(req.ID, map[string]interface{}{
				"tools": []map[string]interface{}{
					{
						"name":        "get_template",
						"description": "Obtiene una plantilla industrial por nombre desde el almacén local.",
						"inputSchema": map[string]interface{}{
							"type":     "object",
							"required": []string{"template_name"},
							"properties": map[string]interface{}{
								"template_name": map[string]string{"type": "string"},
							},
						},
					},
					{
						"name":        "optimize_prompt",
						"description": "Compila y optimiza un prompt. Acepta template_name.",
						"inputSchema": map[string]interface{}{
							"type":     "object",
							"required": []string{"role", "context", "task"},
							"properties": map[string]interface{}{
								"role":          map[string]string{"type": "string"},
								"context":       map[string]string{"type": "string"},
								"task":          map[string]string{"type": "string"},
								"template_name": map[string]string{"type": "string"},
								"constraints":   map[string]interface{}{"type": "array", "items": map[string]string{"type": "string"}},
								"variables":     map[string]interface{}{"type": "object", "additionalProperties": map[string]string{"type": "string"}},
							},
						},
					},
				},
			})

		case "tools/call":
			handleToolCall(req, app)
		}
	}
}
