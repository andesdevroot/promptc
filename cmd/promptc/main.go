package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/andesdevroot/promptc/pkg/core"
	"github.com/andesdevroot/promptc/pkg/sdk"
	"github.com/gorilla/websocket"
)

// --- INFRAESTRUCTURA ---
const configPath = "/Users/cesarrivas/Desktop/GO/promptc/templates.json"
const metricsPath = "/Users/cesarrivas/Desktop/GO/promptc/metrics.json"
const auditPath = "/Users/cesarrivas/Desktop/GO/promptc/audit.log"

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
// AuditEvent representa un evento estructurado de auditoría.
// Cada evento tiene tipo semántico, actor, recurso y resultado.
type AuditEvent struct {
	Timestamp string `json:"ts"`
	Type      string `json:"type"` // KERNEL | MCP | TEMPLATE | INFERENCE | POLICY | SYSTEM
	Action    string `json:"action"`
	Actor     string `json:"actor"` // claude-desktop | promptc-engine | mac-mini | gemini
	Resource  string `json:"resource,omitempty"`
	Result    string `json:"result"` // OK | FAIL | WARN
	LatencyMs int64  `json:"latency_ms,omitempty"`
	Detail    string `json:"detail,omitempty"`
}

// auditLog escribe el evento al archivo de auditoría Y al stream del dashboard.
// El archivo es append-only — nunca se trunca, es el registro regulatorio.
func auditLog(evt AuditEvent) {
	evt.Timestamp = time.Now().Format("2006-01-02T15:04:05.000")

	// Formato legible para el dashboard stream
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

	// 1. Al dashboard WebSocket
	hub.Lock()
	hub.Logs = append(hub.Logs, line)
	for c := range hub.Clients {
		_ = c.WriteMessage(websocket.TextMessage, []byte(line))
	}
	hub.Unlock()

	// 2. A stderr (visible en mcp.log de Claude Desktop)
	fmt.Fprintf(os.Stderr, "%s\n", line)

	// 3. Al archivo de auditoría append-only (registro regulatorio)
	f, err := os.OpenFile(auditPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		jsonLine, _ := json.Marshal(evt)
		fmt.Fprintf(f, "%s\n", string(jsonLine))
		f.Close()
	}
}

// addLog mantiene compatibilidad con logs de sistema genéricos
// que no son eventos de auditoría (WS connect, hot-reload, etc.)
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
}

var metrics = &Metrics{
	TemplateCalls: make(map[string]int64),
}

func loadMetrics() {
	data, err := os.ReadFile(metricsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[METRICS] No hay estado previo, arrancando limpio\n")
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
	fmt.Fprintf(os.Stderr, "[METRICS] Estado restaurado desde disco (guardado: %s)\n",
		snap.SavedAt.Format("2006-01-02 15:04:05"))
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
		fmt.Fprintf(os.Stderr, "[METRICS] Error serializando métricas: %v\n", err)
		return
	}
	tmpPath := metricsPath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "[METRICS] Error escribiendo métricas: %v\n", err)
		return
	}
	if err := os.Rename(tmpPath, metricsPath); err != nil {
		fmt.Fprintf(os.Stderr, "[METRICS] Error en rename atómico: %v\n", err)
		return
	}
	fmt.Fprintf(os.Stderr, "[METRICS] Estado persistido en disco\n")
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
	metrics.Unlock()

	return map[string]interface{}{
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
						Detail:   "Nodo Ollama respondió heartbeat — inferencia local disponible",
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

// --- DASHBOARD HTML ---
const dashboardHTML = `<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <title>PROMPTC // CONTROL_PANEL</title>
    <style>
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body {
            background: #000;
            color: #00ff41;
            font-family: 'Courier New', monospace;
            height: 100vh;
            display: flex;
            flex-direction: column;
            overflow: hidden;
        }
        .header {
            padding: 14px 20px;
            border-bottom: 1px solid #00ff41;
            display: flex;
            align-items: center;
            justify-content: space-between;
            flex-shrink: 0;
        }
        .header h1 { font-size: 1.1em; letter-spacing: 3px; text-shadow: 0 0 10px #00ff41; }
        .blink { animation: blinker 1s linear infinite; }
        @keyframes blinker { 50% { opacity: 0; } }
        .metrics-row {
            display: grid;
            grid-template-columns: repeat(7, 1fr);
            gap: 8px;
            padding: 10px;
            flex-shrink: 0;
        }
        .metric-card {
            border: 1px solid #00ff41;
            padding: 10px 12px;
            background: #050505;
            display: flex;
            flex-direction: column;
            gap: 4px;
        }
        .metric-label { font-size: 0.6em; color: #00aa2a; letter-spacing: 2px; text-transform: uppercase; }
        .metric-value { font-size: 1.3em; font-weight: bold; text-shadow: 0 0 8px #00ff41; }
        .metric-sub { font-size: 0.65em; color: #008820; }
        .online  { color: #00ff41; }
        .offline { color: #ff4141; }
        .warn    { color: #ffaa00; }
        .main-grid {
            display: grid;
            grid-template-columns: 1fr 400px;
            gap: 10px;
            padding: 0 10px 10px 10px;
            flex-grow: 1;
            overflow: hidden;
        }
        .left-col {
            display: flex;
            flex-direction: column;
            gap: 10px;
            overflow: hidden;
            min-height: 0;
        }
        .log-panel {
            border: 1px solid #00ff41;
            background: #050505;
            padding: 12px;
            display: flex;
            flex-direction: column;
            box-shadow: inset 0 0 12px #00ff4122;
            flex: 1 1 0;
            min-height: 0;
            overflow: hidden;
        }
        .log-panel-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 8px;
            border-bottom: 1px solid #003311;
            padding-bottom: 6px;
            flex-shrink: 0;
        }
        .log-panel-header h3 {
            font-size: 0.65em;
            letter-spacing: 2px;
            color: #00aa2a;
        }
        .filter-row {
            display: flex;
            gap: 6px;
        }
        .filter-btn {
            font-size: 0.55em;
            padding: 2px 7px;
            border: 1px solid #00ff41;
            background: transparent;
            color: #00ff41;
            cursor: pointer;
            font-family: inherit;
            letter-spacing: 1px;
            transition: 0.15s;
        }
        .filter-btn.active { background: #00ff41; color: #000; }
        .filter-btn:hover { background: #00ff4133; }
        .log-entries { overflow-y: auto; flex-grow: 1; }

        /* Colores por tipo de evento de auditoría */
        .entry { font-size: 0.72em; margin-bottom: 3px; white-space: pre-wrap; line-height: 1.6; }
        .entry.KERNEL    { color: #00ff41; }
        .entry.MCP       { color: #00ccff; }
        .entry.TEMPLATE  { color: #ffcc00; }
        .entry.INFERENCE { color: #ff9900; }
        .entry.POLICY    { color: #ff4444; }
        .entry.SYSTEM    { color: #556655; }
        .entry.hidden    { display: none; }

        .ranking-panel {
            border: 1px solid #00ff41;
            background: #050505;
            padding: 12px;
            overflow-y: auto;
            flex-shrink: 0;
            max-height: 200px;
            box-shadow: inset 0 0 12px #00ff4122;
        }
        .ranking-panel h3 {
            font-size: 0.65em;
            letter-spacing: 2px;
            color: #00aa2a;
            margin-bottom: 10px;
            border-bottom: 1px solid #003311;
            padding-bottom: 6px;
        }
        .rank-item {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 6px 0;
            border-bottom: 1px solid #001a00;
            font-size: 0.75em;
        }
        .rank-item-left { display: flex; flex-direction: column; flex-grow: 1; margin-right: 12px; }
        .rank-bar {
            height: 4px;
            background: #00ff41;
            margin-top: 4px;
            transition: width 0.5s;
            box-shadow: 0 0 6px #00ff41;
        }
        .editor-panel {
            border: 1px solid #00ff41;
            background: #050505;
            padding: 12px;
            display: flex;
            flex-direction: column;
            gap: 8px;
            box-shadow: inset 0 0 12px #00ff4122;
            overflow: hidden;
        }
        .editor-panel h3 {
            font-size: 0.65em;
            letter-spacing: 2px;
            color: #00aa2a;
            border-bottom: 1px solid #003311;
            padding-bottom: 6px;
            flex-shrink: 0;
        }
        textarea {
            background: #000;
            color: #00ff41;
            border: 1px solid #003311;
            flex-grow: 1;
            padding: 10px;
            font-family: inherit;
            resize: none;
            outline: none;
            font-size: 0.75em;
            line-height: 1.5;
            min-height: 0;
        }
        button {
            background: #00ff41;
            color: #000;
            border: none;
            padding: 10px;
            font-weight: bold;
            cursor: pointer;
            text-transform: uppercase;
            font-family: inherit;
            letter-spacing: 2px;
            font-size: 0.8em;
            transition: 0.15s;
            flex-shrink: 0;
        }
        button:hover { background: #00cc33; box-shadow: 0 0 15px #00ff41; }
        ::-webkit-scrollbar { width: 4px; }
        ::-webkit-scrollbar-thumb { background: #00ff41; }
    </style>
</head>
<body>
    <div class="header">
        <h1>PROMPTC // CONTROL_PANEL_V0.3.0 <span class="blink">_</span></h1>
        <span id="clock" style="font-size:0.8em; color:#00aa2a;"></span>
    </div>

    <div class="metrics-row">
        <div class="metric-card">
            <span class="metric-label">Node Heartbeat</span>
            <span class="metric-value" id="m-node">--</span>
            <span class="metric-sub" id="m-heartbeat">última señal: --</span>
        </div>
        <div class="metric-card">
            <span class="metric-label">Memory Alloc</span>
            <span class="metric-value" id="m-mem">-- MB</span>
            <span class="metric-sub" id="m-goroutines">-- goroutines</span>
        </div>
        <div class="metric-card">
            <span class="metric-label">Avg Latency</span>
            <span class="metric-value" id="m-latency">-- ms</span>
            <span class="metric-sub" id="m-inferences">0 inferencias</span>
        </div>
        <div class="metric-card">
            <span class="metric-label">Token Throughput</span>
            <span class="metric-value" id="m-tps">-- TPS</span>
            <span class="metric-sub" id="m-tokens">0 tokens total</span>
        </div>
        <div class="metric-card">
            <span class="metric-label">Success Ratio</span>
            <span class="metric-value" id="m-ratio">--%</span>
            <span class="metric-sub" id="m-succfail">ok:0 / err:0</span>
        </div>
        <div class="metric-card">
            <span class="metric-label">Gemini Quota</span>
            <span class="metric-value" id="m-gemini">0</span>
            <span class="metric-sub">calls hoy</span>
        </div>
        <div class="metric-card">
            <span class="metric-label">Templates Loaded</span>
            <span class="metric-value" id="m-templates">0</span>
            <span class="metric-sub">en templates.json</span>
        </div>
    </div>

    <div class="main-grid">
        <div class="left-col">
            <div class="log-panel">
                <div class="log-panel-header">
                    <h3>// AUDIT_LOG_STREAM</h3>
                    <div class="filter-row">
                        <button class="filter-btn active" onclick="setFilter('ALL')">ALL</button>
                        <button class="filter-btn" onclick="setFilter('KERNEL')">KERNEL</button>
                        <button class="filter-btn" onclick="setFilter('MCP')">MCP</button>
                        <button class="filter-btn" onclick="setFilter('TEMPLATE')">TEMPLATE</button>
                        <button class="filter-btn" onclick="setFilter('INFERENCE')">INFERENCE</button>
                        <button class="filter-btn" onclick="setFilter('POLICY')">POLICY</button>
                    </div>
                </div>
                <div class="log-entries" id="logs"></div>
            </div>
            <div class="ranking-panel" id="ranking-panel">
                <h3>// TEMPLATE_POPULARITY</h3>
                <div id="ranking"></div>
            </div>
        </div>
        <div class="editor-panel">
            <h3>// TEMPLATES_CONFIG_JSON</h3>
            <textarea id="ed"></textarea>
            <button onclick="save()">SAVE &amp; HOT-RELOAD</button>
        </div>
    </div>

    <script>
        // Reloj
        setInterval(() => {
            document.getElementById('clock').textContent = new Date().toLocaleTimeString('es-CL');
        }, 1000);

        // Filtro activo
        let activeFilter = 'ALL';

        function setFilter(f) {
            activeFilter = f;
            document.querySelectorAll('.filter-btn').forEach(b => {
                b.classList.toggle('active', b.textContent === f);
            });
            document.querySelectorAll('.entry').forEach(el => {
                if (f === 'ALL') {
                    el.classList.remove('hidden');
                } else {
                    el.classList.toggle('hidden', !el.classList.contains(f));
                }
            });
        }

        // Detectar tipo de evento del texto del log para colorear
        function detectType(text) {
            const types = ['KERNEL', 'MCP', 'TEMPLATE', 'INFERENCE', 'POLICY'];
            for (const t of types) {
                if (text.includes(t)) return t;
            }
            return 'SYSTEM';
        }

        // WebSocket logs
        const logsEl = document.getElementById('logs');
        const ws = new WebSocket('ws://' + location.host + '/ws');
        ws.onmessage = (e) => {
            const div = document.createElement('div');
            const type = detectType(e.data);
            div.className = 'entry ' + type;
            if (activeFilter !== 'ALL' && type !== activeFilter) {
                div.classList.add('hidden');
            }
            div.textContent = e.data;
            logsEl.appendChild(div);
            logsEl.scrollTop = logsEl.scrollHeight;
        };

        // Editor templates
        fetch('/api/config')
            .then(r => r.json())
            .then(d => {
                document.getElementById('ed').value = JSON.stringify(d, null, 2);
                document.getElementById('m-templates').textContent = Object.keys(d).length;
            });

        function save() {
            fetch('/api/config', { method: 'POST', body: document.getElementById('ed').value })
                .then(r => {
                    if (!r.ok) alert('PERSISTENCE_ERROR');
                    else fetch('/api/config').then(r => r.json()).then(d => {
                        document.getElementById('m-templates').textContent = Object.keys(d).length;
                    });
                });
        }

        // Polling métricas
        function updateMetrics() {
            fetch('/api/metrics')
                .then(r => r.json())
                .then(d => {
                    const nodeEl = document.getElementById('m-node');
                    if (d.node_online) {
                        nodeEl.textContent = 'ONLINE';
                        nodeEl.className = 'metric-value online';
                    } else {
                        nodeEl.textContent = 'OFFLINE';
                        nodeEl.className = 'metric-value offline';
                    }
                    document.getElementById('m-heartbeat').textContent = 'última señal: ' + d.last_heartbeat;
                    document.getElementById('m-mem').textContent = d.mem_alloc_mb.toFixed(1) + ' MB';
                    document.getElementById('m-goroutines').textContent = d.goroutines + ' goroutines';
                    document.getElementById('m-latency').textContent = d.avg_latency_ms.toFixed(0) + ' ms';
                    document.getElementById('m-inferences').textContent = d.inference_count + ' inferencias';
                    document.getElementById('m-tps').textContent = d.token_throughput.toFixed(1) + ' TPS';
                    document.getElementById('m-tokens').textContent = (d.total_tokens || 0) + ' tokens total';

                    const ratioEl = document.getElementById('m-ratio');
                    ratioEl.textContent = d.success_ratio.toFixed(1) + '%';
                    if (d.success_ratio >= 90) ratioEl.className = 'metric-value online';
                    else if (d.success_ratio >= 70) ratioEl.className = 'metric-value warn';
                    else ratioEl.className = 'metric-value offline';
                    document.getElementById('m-succfail').textContent = 'ok:' + d.success_count + ' / err:' + d.fail_count;

                    const geminiEl = document.getElementById('m-gemini');
                    geminiEl.textContent = d.gemini_calls;
                    if (d.gemini_calls > 1400) geminiEl.className = 'metric-value offline';
                    else if (d.gemini_calls > 1000) geminiEl.className = 'metric-value warn';
                    else geminiEl.className = 'metric-value online';

                    if (d.template_ranking && d.template_ranking.length > 0) {
                        const rankingEl = document.getElementById('ranking');
                        const sorted = d.template_ranking.sort((a, b) => b.calls - a.calls);
                        const maxCalls = sorted[0].calls || 1;
                        rankingEl.innerHTML = sorted.map(t => {
                            const pct = (t.calls / maxCalls * 100).toFixed(0);
                            return '<div class="rank-item">' +
                                '<div class="rank-item-left">' +
                                    '<div>' + t.name + '</div>' +
                                    '<div class="rank-bar" style="width:' + pct + '%"></div>' +
                                '</div>' +
                                '<div style="color:#00aa2a">' + t.calls + '</div>' +
                            '</div>';
                        }).join('');
                    }
                });
        }

        updateMetrics();
        setInterval(updateMetrics, 3000);
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
			"version":         "0.3.0",
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
					Detail: fmt.Sprintf("%d templates recargados desde editor web", len(n)),
				})
			}
		}
	})

	if err := http.ListenAndServe(":8080", mux); err != nil {
		addLog(fmt.Sprintf("Port 8080 busy: %v — MCP-ONLY Mode", err))
	}
}

// startTime para el health endpoint
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

	// Evento MCP: Claude Desktop invocó una herramienta
	auditLog(AuditEvent{
		Type:     "MCP",
		Action:   "TOOL_INVOKED",
		Actor:    "claude-desktop",
		Resource: call.Name,
		Result:   "OK",
		Detail:   "Solicitud recibida vía MCP stdio",
	})

	switch call.Name {
	case "get_template":
		var args struct {
			Name string `json:"template_name"`
		}
		if err := json.Unmarshal(call.Arguments, &args); err != nil {
			auditLog(AuditEvent{
				Type:     "TEMPLATE",
				Action:   "GET_ARGS_INVALID",
				Actor:    "promptc-engine",
				Resource: "get_template",
				Result:   "FAIL",
				Detail:   err.Error(),
			})
			sendResponse(req.ID, map[string]interface{}{
				"content": []map[string]interface{}{
					{"type": "text", "text": fmt.Sprintf("Error: argumentos inválidos: %v", err)},
				},
			})
			return
		}

		hub.Lock()
		tmpl, ok := hub.Templates[args.Name]
		hub.Unlock()

		if !ok {
			auditLog(AuditEvent{
				Type:     "TEMPLATE",
				Action:   "GET_NOT_FOUND",
				Actor:    "promptc-engine",
				Resource: args.Name,
				Result:   "FAIL",
				Detail:   "Template no registrado en templates.json",
			})
			sendResponse(req.ID, map[string]interface{}{
				"content": []map[string]interface{}{
					{"type": "text", "text": fmt.Sprintf("Error: template '%s' no encontrado", args.Name)},
				},
			})
			return
		}

		recordTemplatCall(args.Name)
		auditLog(AuditEvent{
			Type:     "TEMPLATE",
			Action:   "GET_SERVED",
			Actor:    "promptc-engine",
			Resource: args.Name,
			Result:   "OK",
			Detail:   fmt.Sprintf("desc=%q content_len=%d", tmpl.Description, len(tmpl.Content)),
		})
		sendResponse(req.ID, map[string]interface{}{
			"content": []map[string]interface{}{
				{"type": "text", "text": tmpl.Content},
			},
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
			auditLog(AuditEvent{
				Type:     "INFERENCE",
				Action:   "OPTIMIZE_ARGS_INVALID",
				Actor:    "promptc-engine",
				Resource: "optimize_prompt",
				Result:   "FAIL",
				Detail:   err.Error(),
			})
			recordInference(false, 0, 0, false)
			sendResponse(req.ID, map[string]interface{}{
				"content": []map[string]interface{}{
					{"type": "text", "text": fmt.Sprintf("Error: argumentos inválidos: %v", err)},
				},
			})
			return
		}

		// Inyección de template como base del Task
		task := args.Task
		if args.Template != "" {
			hub.Lock()
			tmpl, ok := hub.Templates[args.Template]
			hub.Unlock()
			if ok {
				task = tmpl.Content
				recordTemplatCall(args.Template)
				auditLog(AuditEvent{
					Type:     "TEMPLATE",
					Action:   "INJECT_AS_TASK",
					Actor:    "promptc-engine",
					Resource: args.Template,
					Result:   "OK",
					Detail:   fmt.Sprintf("Template inyectado como base — variables a resolver: %d", len(args.Variables)),
				})
			} else {
				auditLog(AuditEvent{
					Type:     "TEMPLATE",
					Action:   "INJECT_NOT_FOUND",
					Actor:    "promptc-engine",
					Resource: args.Template,
					Result:   "WARN",
					Detail:   "Template no encontrado — usando Task directo del argumento",
				})
			}
		}

		// Enrutamiento al nodo de inferencia
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
			Detail:   fmt.Sprintf("role=%q constraints=%d variables=%d", args.Role, len(args.Constraints), len(args.Variables)),
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
			auditLog(AuditEvent{
				Type:      "INFERENCE",
				Action:    "PIPELINE_FAIL",
				Actor:     inferenceActor,
				Resource:  "optimize_prompt",
				Result:    "FAIL",
				LatencyMs: latencyMs,
				Detail:    err.Error(),
			})
			recordInference(false, latencyMs, 0, !nodeOnline)
			sendResponse(req.ID, map[string]interface{}{
				"content": []map[string]interface{}{
					{"type": "text", "text": fmt.Sprintf("Error en pipeline de optimización: %v", err)},
				},
			})
			return
		}

		auditLog(AuditEvent{
			Type:      "INFERENCE",
			Action:    "PIPELINE_OK",
			Actor:     inferenceActor,
			Resource:  "optimize_prompt",
			Result:    "OK",
			LatencyMs: latencyMs,
			Detail: fmt.Sprintf("tokens~%d soberanía=%s", estimatedTokens, func() string {
				if nodeOnline {
					return "LOCAL"
				}
				return "CLOUD"
			}()),
		})
		recordInference(true, latencyMs, estimatedTokens, !nodeOnline)
		sendResponse(req.ID, map[string]interface{}{
			"content": []map[string]interface{}{
				{"type": "text", "text": res},
			},
		})

	default:
		auditLog(AuditEvent{
			Type:     "MCP",
			Action:   "TOOL_NOT_FOUND",
			Actor:    "claude-desktop",
			Resource: call.Name,
			Result:   "FAIL",
			Detail:   "Herramienta no registrada en el servidor MCP",
		})
		sendResponse(req.ID, map[string]interface{}{
			"content": []map[string]interface{}{
				{"type": "text", "text": fmt.Sprintf("Error: herramienta '%s' no registrada", call.Name)},
			},
		})
	}
}

// --- MAIN ---
func main() {
	// 1. Restaurar métricas
	loadMetrics()

	// 2. Cargar templates
	file, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[WARN] No se pudo leer templates.json: %v\n", err)
	} else {
		if err := json.Unmarshal(file, &hub.Templates); err != nil {
			fmt.Fprintf(os.Stderr, "[WARN] templates.json malformado: %v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "[INFO] %d templates cargados\n", len(hub.Templates))
		}
	}

	// 3. Dashboard
	go startDashboard()

	// 4. Heartbeat
	remoteIP := os.Getenv("PROMPTC_MACMINI_IP")
	if remoteIP == "" {
		remoteIP = "100.90.6.101"
	}
	startHeartbeat(remoteIP)

	// 5. Persistencia periódica
	startMetricsPersistence()

	// 6. SDK
	app, err := sdk.NewSDK(context.Background(), os.Getenv("GEMINI_API_KEY"), remoteIP)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[SDK_ERROR] %v — continuando sin optimizadores\n", err)
	}

	// 7. Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-sigChan
		auditLog(AuditEvent{
			Type:   "KERNEL",
			Action: "SHUTDOWN",
			Actor:  "promptc-engine",
			Result: "OK",
			Detail: fmt.Sprintf("Señal recibida: %v — flush de métricas iniciado", sig),
		})
		saveMetrics()
		os.Exit(0)
	}()

	// 8. Evento de arranque del kernel
	auditLog(AuditEvent{
		Type:   "KERNEL",
		Action: "BOOT",
		Actor:  "promptc-engine",
		Result: "OK",
		Detail: fmt.Sprintf("PROMPTC v0.3.0 iniciado — nodo=%s templates=%d inferencias_previas=%d",
			remoteIP,
			len(hub.Templates),
			atomic.LoadInt64(&metrics.InferenceCount),
		),
	})

	// 9. Scanner MCP
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	for scanner.Scan() {
		var req JSONRPCMessage
		if err := json.Unmarshal(scanner.Bytes(), &req); err != nil {
			fmt.Fprintf(os.Stderr, "[PARSE_ERROR] %v\n", err)
			continue
		}

		switch req.Method {
		case "initialize":
			auditLog(AuditEvent{
				Type:   "MCP",
				Action: "HANDSHAKE_INIT",
				Actor:  "claude-desktop",
				Result: "OK",
				Detail: "Protocolo MCP 2024-11-05 — negociación iniciada",
			})
			sendResponse(req.ID, map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"serverInfo":      map[string]string{"name": "PROMPTC", "version": "0.3.0"},
				"capabilities":    map[string]interface{}{"tools": map[string]interface{}{}},
			})

		case "notifications/initialized":
			auditLog(AuditEvent{
				Type:   "MCP",
				Action: "HANDSHAKE_CONFIRMED",
				Actor:  "claude-desktop",
				Result: "OK",
				Detail: "Canal MCP establecido — herramientas disponibles",
			})

		case "tools/list":
			auditLog(AuditEvent{
				Type:   "MCP",
				Action: "TOOLS_LIST_REQUESTED",
				Actor:  "claude-desktop",
				Result: "OK",
				Detail: "Enviando schema de 2 herramientas: get_template, optimize_prompt",
			})
			sendResponse(req.ID, map[string]interface{}{
				"tools": []map[string]interface{}{
					{
						"name":        "get_template",
						"description": "Obtiene una plantilla industrial por nombre desde el almacén local.",
						"inputSchema": map[string]interface{}{
							"type":     "object",
							"required": []string{"template_name"},
							"properties": map[string]interface{}{
								"template_name": map[string]string{
									"type":        "string",
									"description": "Nombre exacto de la plantilla registrada en templates.json",
								},
							},
						},
					},
					{
						"name":        "optimize_prompt",
						"description": "Compila y optimiza un prompt usando el Mac Mini vía Tailscale con fallback a Gemini. Acepta template_name para usar una plantilla como base con resolución automática de variables.",
						"inputSchema": map[string]interface{}{
							"type":     "object",
							"required": []string{"role", "context", "task"},
							"properties": map[string]interface{}{
								"role": map[string]string{
									"type":        "string",
									"description": "Rol del agente o sistema que ejecutará el prompt",
								},
								"context": map[string]string{
									"type":        "string",
									"description": "Contexto de negocio o técnico relevante para el prompt",
								},
								"task": map[string]string{
									"type":        "string",
									"description": "Tarea concreta. Ignorada si se provee template_name",
								},
								"template_name": map[string]string{
									"type":        "string",
									"description": "Nombre del template en templates.json para usar como base del Task con resolución automática de {{variables}}",
								},
								"constraints": map[string]interface{}{
									"type":        "array",
									"description": "Restricciones opcionales",
									"items":       map[string]string{"type": "string"},
								},
								"variables": map[string]interface{}{
									"type":        "object",
									"description": "Variables de sustitución para resolver {{placeholders}} del template",
									"additionalProperties": map[string]string{
										"type": "string",
									},
								},
							},
						},
					},
				},
			})

		case "tools/call":
			handleToolCall(req, app)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "[SCANNER_ERROR] %v\n", err)
	}
}
