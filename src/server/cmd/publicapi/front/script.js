const API_BASE = 'http://127.0.0.1:8080/api/v1';
let pollInterval = null;

document.addEventListener('DOMContentLoaded', () => {
    const currReq = getCookie('curr_req');
    if (currReq) startPolling(currReq);

    document.getElementById('upload-form').addEventListener('submit', async (e) => {
        e.preventDefault();

        const files = document.getElementById('file-input').files;
        // const groupCount = parseInt(document.getElementById('group-count').value, 10);

        // if (files.length < groupCount) {
        //     setStatus('Error: Number of files must be at least equal to group count.', true);
        //     return;
        // }

        const formData = new FormData();
        for (let file of files) {
            formData.append('file', file);
        }
        formData.append('group_count', 1);

        setStatus('Загрузка файлов...', false);
        document.getElementById('upload-btn').disabled = true;

        try {
            const res = await fetch(`${API_BASE}/uploadFiles`, {
                method: 'POST',
                body: formData
            });
            const data = await res.json();
            if (!res.ok) {
                const errMsg = data.error || res.statusText;
                throw new Error(errMsg);
            }
            const { uuid } = data;
            setCookie('curr_req', uuid, 1);
            startPolling(uuid);
        } catch (err) {
            setStatus('Error: ' + err.message, true);
            document.getElementById('upload-btn').disabled = false;
        }
    });

    document.getElementById('new-request').addEventListener('click', () => {
        deleteCookie('curr_req');
        location.reload();
    });
});

function startPolling(uuid) {
    document.getElementById('upload-form').style.display = 'none';
    document.getElementById('new-request').style.display = 'none';
    setStatus('Processing... Please wait.', false);
    pollInterval = setInterval(async () => {
        try {
            const res = await fetch(`${API_BASE}/getClusterizations/${uuid}`);
            const data = await res.json();
            if (res.status === 404 || data.error === 'NOT_READY' || data.error === 'NO_RESULT') return;

            if (data.error === 'REQUEST_FAILED') {
                clearInterval(pollInterval);
                setStatus('Обработка не удалась. Повторите попытку.', true);
                document.getElementById('new-request').style.display = 'inline-block';
                return;
            }

            if (!res.ok) throw new Error(data.error || res.statusText);
            clearInterval(pollInterval);
            renderGraph(data.result);
            setStatus('Отчет готов', false);
            document.getElementById('new-request').style.display = 'inline-block';
        } catch (err) {
            console.error('Polling error:', err);
            clearInterval(pollInterval);
            setStatus('Ошибка при запросе. Повторите позже.', true);
            document.getElementById('new-request').style.display = 'inline-block';
        }
    }, 1000);
}
function setStatus(msg, isError) {
    const statusEl = document.getElementById('status');
    statusEl.textContent = msg;
    statusEl.style.color = isError ? 'red' : 'black';
}

function setCookie(name, value, days) {
    const d = new Date();
    d.setTime(d.getTime() + (days*24*60*60*1000));
    document.cookie = `${name}=${value};expires=${d.toUTCString()};path=/`;
}

function getCookie(name) {
    const v = document.cookie.match('(^|;) ?' + name + '=([^;]*)(;|$)');
    return v ? v[2] : null;
}

function deleteCookie(name) {
    document.cookie = `${name}=; Max-Age=0; path=/`;
}


function renderGraph(graphData) {
    const container = document.getElementById('graph-container');
    container.innerHTML = '';
    const { nodes, links } = graphData;
    const width = container.clientWidth;
    const height = container.clientHeight;
    const svgNS = 'http://www.w3.org/2000/svg';

    const svg = document.createElementNS(svgNS, 'svg');
    svg.setAttribute('width', width);
    svg.setAttribute('height', height);
    svg.style.position = 'absolute';
    svg.style.top = 0;
    svg.style.left = 0;
    container.appendChild(svg);

    const clusterIds = [...new Set(nodes.map(n => n.cluster))];
    const clusterColors = {};
    clusterIds.forEach(cid => {
        clusterColors[cid] = '#' + Math.floor(Math.random() * 0xffffff).toString(16).padStart(6, '0');
    });

    const cx = width / 2, cy = height / 2, R = Math.min(cx, cy) - 100;
    const positions = {};
    nodes.forEach((node, i) => {
        const angle = (2 * Math.PI * i) / nodes.length;
        positions[node.id] = {
            x: cx + R * Math.cos(angle),
            y: cy + R * Math.sin(angle),
        };
    });

    const lineElems = {};
    links.forEach(link => {
        const p1 = positions[link.source];
        const p2 = positions[link.target];
        const line = document.createElementNS(svgNS, 'line');
        line.setAttribute('x1', p1.x);
        line.setAttribute('y1', p1.y);
        line.setAttribute('x2', p2.x);
        line.setAttribute('y2', p2.y);
        line.setAttribute('stroke', '#999');
        line.setAttribute('stroke-width', Math.max(1, link.weight * 3));
        line.setAttribute('stroke-opacity', link.weight);
        svg.appendChild(line);
        lineElems[`${link.source}-${link.target}`] = line;
    });

    const circleElems = {}, textElems = {};
    nodes.forEach(node => {
        const { x, y } = positions[node.id];
        const circle = document.createElementNS(svgNS, 'circle');
        circle.setAttribute('cx', x);
        circle.setAttribute('cy', y);
        circle.setAttribute('r', 28);
        circle.setAttribute('fill', clusterColors[node.cluster]);
        circle.setAttribute('data-id', node.id);
        circle.style.cursor = 'move';
        svg.appendChild(circle);
        circleElems[node.id] = circle;

        // … внутри renderGraph, вместо bisherigen label-генерации:
        const label = document.createElement('a');
        label.href = `/api/v1/downloadFile/${node.id}`;
// если хотите явно скачать, а не открывать в новой вкладке:
// label.setAttribute('download', '');
        label.target = '_blank';
        label.textContent = node.title;
        label.className = 'node-label';

// позиционирование
        label.style.left = `${x}px`;
        label.style.top  = `${y}px`;
        label.style.position = 'absolute';  // если не задали в CSS

        container.appendChild(label);

        textElems[node.id] = label;
    });

    let dragNode = null;
    let offset = { x: 0, y: 0 };
    svg.addEventListener('mousedown', e => {
        if (e.target.tagName === 'circle') {
            dragNode = e.target.getAttribute('data-id');
            offset.x = positions[dragNode].x - e.clientX;
            offset.y = positions[dragNode].y - e.clientY;
        }
    });
    document.addEventListener('mousemove', e => {
        if (!dragNode) return;
        const x = e.clientX + offset.x;
        const y = e.clientY + offset.y;
        positions[dragNode] = { x, y };
        circleElems[dragNode].setAttribute('cx', x);
        circleElems[dragNode].setAttribute('cy', y);
        textElems[dragNode].style.left = `${x}px`;
        textElems[dragNode].style.top = `${y}px`;
        links.forEach(link => {
            const key = `${link.source}-${link.target}`;
            const line = lineElems[key];
            if (link.source === dragNode) {
                line.setAttribute('x1', x);
                line.setAttribute('y1', y);
            }
            if (link.target === dragNode) {
                line.setAttribute('x2', x);
                line.setAttribute('y2', y);
            }
        });
    });
    document.addEventListener('mouseup', () => {
        dragNode = null;
    });
}

