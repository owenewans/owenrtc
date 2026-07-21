// wails injects bound methods at window.go.panel.API; fall back to HTTP fetch in browser
function wapi() {
    return (window.go && window.go.panel && window.go.panel.API) || null
}

async function getMode() {
    const a = wapi()
    if (a) return a.Mode()
    const r = await fetch("/api/mode")
    return r.json()
}

async function getJitsiHosts() {
    const a = wapi()
    if (a) return a.JitsiHosts()
    const r = await fetch("/api/rooms/jitsi")
    return r.json()
}

async function testRoom(provider, transport, roomID) {
    const a = wapi()
    if (a) return a.TestRoom(provider, transport, roomID)
    const r = await fetch("/api/rooms/test", {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({provider, transport, room_id: roomID}),
    })
    const data = await r.json()
    return data.result || "failed"
}

async function createInstanceAPI(inst) {
    const a = wapi()
    if (a) return a.CreateInstance(inst)
    const r = await fetch("/api/servers", {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify(inst),
    })
    return r.json()
}

async function listServers() {
    const a = wapi()
    if (a) return a.ListInstances()
    const r = await fetch("/api/servers")
    return r.json()
}

async function installAPI(host, port, user, password) {
    const a = wapi()
    if (a) return a.Install(host, port, user, password)
    const r = await fetch("/api/install", {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({host, port, user, password}),
    })
    const res = await r.json()
    return res.status || "failed"
}

const transports = {
    jitsi:    ["datachannel", "vp8channel", "seichannel", "videochannel"],
    telemost: ["vp8channel", "videochannel"],
    wbstream: ["vp8channel", "seichannel", "videochannel"],
}

async function init() {
    let mode
    try {
        mode = await getMode()
    } catch (e) {
        mode = {mode: "client"}
    }
    document.getElementById("mode-badge").textContent = mode.mode + " mode"

    if (mode.mode === "server") {
        showMain()
    } else {
        show("setup")
    }
}

function show(id) {
    document.querySelectorAll(".screen").forEach(s => s.classList.add("hidden"))
    document.getElementById(id).classList.remove("hidden")
}

function showMain() {
    loadServers()
    show("main")
}

function showCreate() {
    loadJitsiHosts()
    updateTransports()
    show("create")
}

async function loadServers() {
    let servers = []
    try { servers = await listServers() } catch (e) {}
    const list = document.getElementById("server-list")
    list.innerHTML = ""
    servers.forEach(srv => {
        const card = document.createElement("div")
        card.className = "server-card"
        card.innerHTML = `
            <div>
                <div class="name">${srv.name}</div>
                <div class="meta">${srv.provider} / ${srv.transport}</div>
            </div>
            <button class="btn-link" onclick="alert('todo')">manage</button>
        `
        list.appendChild(card)
    })
}

async function loadJitsiHosts() {
    let hosts = []
    try { hosts = await getJitsiHosts() } catch (e) {}
    const dl = document.getElementById("jitsi-hosts")
    dl.innerHTML = ""
    hosts.forEach(h => {
        const opt = document.createElement("option")
        opt.value = h
        dl.appendChild(opt)
    })
}

function updateTransports() {
    const prov = document.getElementById("inst-provider").value
    const sel = document.getElementById("inst-transport")
    sel.innerHTML = ""
    transports[prov].forEach(t => {
        const opt = document.createElement("option")
        opt.value = t
        opt.textContent = t
        sel.appendChild(opt)
    })
}

function fillRoomURL() {
    const host = document.getElementById("jitsi-host").value
    if (host) {
        const room = Math.random().toString(36).slice(2, 10)
        document.getElementById("inst-room").value = `https://${host}/${room}`
    }
}

async function doInstall() {
    const host = document.getElementById("ssh-host").value
    const port = document.getElementById("ssh-port").value || "22"
    const user = document.getElementById("ssh-user").value
    const pass = document.getElementById("ssh-pass").value

    if (!host || !user || !pass) {
        document.getElementById("install-status").textContent = "fill all fields"
        return
    }

    document.getElementById("install-status").textContent = "installing..."
    document.getElementById("install-btn").disabled = true

    let result = "failed"
    try {
        result = await installAPI(host, parseInt(port), user, pass)
    } catch (e) {
        result = String(e)
    }

    if (result === "ok" || result.includes("ok")) {
        document.getElementById("install-status").textContent = "installation ok"
        setTimeout(showMain, 1000)
    } else {
        document.getElementById("install-status").textContent = result || "failed"
        document.getElementById("install-btn").disabled = false
    }
}

async function testRoomAPI() {
    const provider = document.getElementById("inst-provider").value
    const transport = document.getElementById("inst-transport").value
    const room = document.getElementById("inst-room").value
    const el = document.getElementById("test-result")

    if (!room) {
        el.textContent = "enter room url"
        return
    }

    el.textContent = "testing..."
    try {
        const result = await testRoom(provider, transport, room)
        el.textContent = result
    } catch (e) {
        el.textContent = String(e)
    }
}

async function createInstance() {
    const inst = {
        name: document.getElementById("inst-name").value,
        provider: document.getElementById("inst-provider").value,
        transport: document.getElementById("inst-transport").value,
        room_id: document.getElementById("inst-room").value,
        limits: {
            traffic_limit: parseSize(document.getElementById("inst-traffic").value),
            speed_limit: parseSpeed(document.getElementById("inst-speed").value),
        },
    }
    try { await createInstanceAPI(inst) } catch (e) {}
    showMain()
}

function parseSize(s) {
    if (!s) return 0
    const m = s.match(/^(\d+)\s*(gb|mb|tb|kb)?/i)
    if (!m) return 0
    const n = parseInt(m[1])
    const u = (m[2] || "").toLowerCase()
    const mult = {kb: 1e3, mb: 1e6, gb: 1e9, tb: 1e12}
    return n * (mult[u] || 1)
}

function parseSpeed(s) {
    if (!s) return 0
    const m = s.match(/^(\d+)\s*(gb|mb|tb|kb)?\/s?/i)
    if (!m) return 0
    const n = parseInt(m[1])
    const u = (m[2] || "").toLowerCase()
    const mult = {kb: 1e3, mb: 1e6, gb: 1e9, tb: 1e12}
    return n * (mult[u] || 1)
}

window.showMain = showMain
window.showCreate = showCreate
window.updateTransports = updateTransports
window.fillRoomURL = fillRoomURL
window.doInstall = doInstall
window.testRoom = testRoomAPI
window.createInstance = createInstance

init()
