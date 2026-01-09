// Main Application Module
let currentZoom = 1;
let currentView = 'graph';
let currentLayout = 'TD';
let selectedPackage = null;
let allPackages = [];
let isRemoteMode = false;
let currentRemoteModule = '';
let packageListMode = 'tree'; // 'flat' or 'tree'
let expandedNodes = new Set(); // Track expanded tree nodes

// Pan and Zoom state
let panX = 0;
let panY = 0;
let isPanning = false;
let startPanX = 0;
let startPanY = 0;
let zoomIndicatorTimeout = null;

const MIN_ZOOM = 0.1;
const MAX_ZOOM = 10;
const ZOOM_STEP = 0.1;

// Theme Management
const theme = {
    current: 'light',
    
    init() {
        const savedTheme = localStorage.getItem('dep-viewer-theme');
        if (savedTheme) {
            this.current = savedTheme;
        } else {
            if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
                this.current = 'dark';
            }
        }
        this.apply();
        
        window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
            if (!localStorage.getItem('dep-viewer-theme')) {
                this.current = e.matches ? 'dark' : 'light';
                this.apply();
            }
        });
        
        const toggleBtn = document.getElementById('themeToggle');
        if (toggleBtn) {
            toggleBtn.addEventListener('click', () => this.toggle());
        }
    },
    
    toggle() {
        this.current = this.current === 'dark' ? 'light' : 'dark';
        localStorage.setItem('dep-viewer-theme', this.current);
        this.apply();
    },
    
    apply() {
        if (this.current === 'dark') {
            document.body.setAttribute('data-theme', 'dark');
        } else {
            document.body.removeAttribute('data-theme');
        }
        
        mermaid.initialize({
            startOnLoad: false,
            theme: this.current === 'dark' ? 'dark' : 'default',
            maxEdges: 2000,  // Increase edge limit for large dependency graphs
            flowchart: {
                useMaxWidth: false,
                htmlLabels: true,
                curve: 'basis'
            }
        });
        
        if (currentView === 'graph') {
            refresh();
        }
    }
};

// Initialize mermaid
mermaid.initialize({
    startOnLoad: false,
    theme: 'default',
    maxEdges: 2000,  // Increase edge limit for large dependency graphs
    flowchart: {
        useMaxWidth: false,
        htmlLabels: true,
        curve: 'basis'
    }
});

// Initialize application
async function init() {
    theme.init();
    initPanZoom();
    initRemoteModuleInput();
    const hasLocalModule = await loadModuleName();
    if (hasLocalModule) {
        await loadPackages();
        await refresh();
    }
}

// Initialize remote module input
function initRemoteModuleInput() {
    const input = document.getElementById('remoteModuleInput');
    if (input) {
        // Fetch versions when input loses focus or Enter is pressed
        input.addEventListener('blur', fetchVersions);
        input.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                fetchVersions();
            }
        });
    }
}

// Fetch versions for a module from Go proxy
async function fetchVersions() {
    const input = document.getElementById('remoteModuleInput');
    const versionSelect = document.getElementById('versionSelect');
    const spinner = document.getElementById('loadingSpinner');
    
    let modulePath = input.value.trim();
    // Remove http:// or https:// prefix if present
    modulePath = modulePath.replace(/^https?:\/\//, '');
    input.value = modulePath;
    
    if (!modulePath) {
        versionSelect.disabled = true;
        versionSelect.innerHTML = `<option value="">${i18n.t('selectVersion')}</option>`;
        return;
    }
    
    spinner.classList.remove('hidden');
    versionSelect.disabled = true;
    
    try {
        const response = await fetch('/api/versions?module=' + encodeURIComponent(modulePath));
        const data = await response.json();
        
        if (data.error) {
            versionSelect.innerHTML = `<option value="">${i18n.t('errorFetchVersions')}</option>`;
        } else if (data.versions && data.versions.length > 0) {
            versionSelect.innerHTML = data.versions.map((v, i) => {
                const label = i === 0 ? `${v} (${i18n.t('latestVersion')})` : v;
                return `<option value="${v}">${label}</option>`;
            }).join('');
            versionSelect.disabled = false;
        } else {
            versionSelect.innerHTML = `<option value="">${i18n.t('noVersions')}</option>`;
        }
    } catch (e) {
        console.error('Failed to fetch versions:', e);
        versionSelect.innerHTML = `<option value="">${i18n.t('errorFetchVersions')}</option>`;
    } finally {
        spinner.classList.add('hidden');
    }
}

// Analyze remote module
async function analyzeRemoteModule() {
    const input = document.getElementById('remoteModuleInput');
    const versionSelect = document.getElementById('versionSelect');
    const spinner = document.getElementById('loadingSpinner');
    const analyzeBtn = document.getElementById('analyzeBtn');
    
    let modulePath = input.value.trim();
    // Remove http:// or https:// prefix if present
    modulePath = modulePath.replace(/^https?:\/\//, '');
    input.value = modulePath;
    
    const version = versionSelect.value;
    
    if (!modulePath) {
        alert('Please enter a module path');
        return;
    }
    
    spinner.classList.remove('hidden');
    analyzeBtn.disabled = true;
    analyzeBtn.textContent = i18n.t('analyzing');
    
    try {
        const url = `/api/analyze?module=${encodeURIComponent(modulePath)}${version ? '&version=' + encodeURIComponent(version) : ''}`;
        const response = await fetch(url);
        const data = await response.json();
        
        if (data.error) {
            alert(data.error);
            return;
        }
        
        // Switch to remote mode
        isRemoteMode = true;
        currentRemoteModule = modulePath + (version ? '@' + version : '');
        document.getElementById('moduleName').textContent = currentRemoteModule;
        
        // Clear selection and reload
        selectedPackage = null;
        await loadPackages();
        await refresh();
    } catch (e) {
        console.error('Failed to analyze module:', e);
        alert(i18n.t('errorAnalyze'));
    } finally {
        spinner.classList.add('hidden');
        analyzeBtn.disabled = false;
        analyzeBtn.textContent = i18n.t('analyze');
    }
}

// Reset to local module
async function resetToLocal() {
    const spinner = document.getElementById('loadingSpinner');
    spinner.classList.remove('hidden');
    
    try {
        await fetch('/api/reset');
        
        isRemoteMode = false;
        currentRemoteModule = '';
        document.getElementById('remoteModuleInput').value = '';
        document.getElementById('versionSelect').innerHTML = `<option value="">${i18n.t('selectVersion')}</option>`;
        document.getElementById('versionSelect').disabled = true;
        
        selectedPackage = null;
        await loadModuleName();
        await loadPackages();
        await refresh();
    } catch (e) {
        console.error('Failed to reset:', e);
    } finally {
        spinner.classList.add('hidden');
    }
}

// Initialize pan and zoom functionality
function initPanZoom() {
    const viewport = document.getElementById('graphView');
    if (!viewport) return;

    // Mouse wheel zoom
    viewport.addEventListener('wheel', (e) => {
        if (currentView !== 'graph') return;
        e.preventDefault();
        
        const rect = viewport.getBoundingClientRect();
        const mouseX = e.clientX - rect.left;
        const mouseY = e.clientY - rect.top;
        
        // Calculate zoom
        const delta = e.deltaY > 0 ? -ZOOM_STEP : ZOOM_STEP;
        const newZoom = Math.min(MAX_ZOOM, Math.max(MIN_ZOOM, currentZoom + delta * currentZoom));
        
        if (newZoom !== currentZoom) {
            // Zoom towards mouse position
            const scale = newZoom / currentZoom;
            panX = mouseX - (mouseX - panX) * scale;
            panY = mouseY - (mouseY - panY) * scale;
            currentZoom = newZoom;
            applyTransform();
            showZoomIndicator();
        }
    }, { passive: false });

    // Pan with mouse drag
    viewport.addEventListener('mousedown', (e) => {
        if (currentView !== 'graph') return;
        if (e.button !== 0) return; // Only left click
        
        isPanning = true;
        startPanX = e.clientX - panX;
        startPanY = e.clientY - panY;
        viewport.style.cursor = 'grabbing';
    });

    document.addEventListener('mousemove', (e) => {
        if (!isPanning) return;
        
        panX = e.clientX - startPanX;
        panY = e.clientY - startPanY;
        applyTransform();
    });

    document.addEventListener('mouseup', () => {
        if (isPanning) {
            isPanning = false;
            const viewport = document.getElementById('graphViewport');
            if (viewport) viewport.style.cursor = 'grab';
        }
    });

    // Touch support for mobile
    let lastTouchDistance = 0;
    let lastTouchCenter = { x: 0, y: 0 };

    viewport.addEventListener('touchstart', (e) => {
        if (currentView !== 'graph') return;
        
        if (e.touches.length === 1) {
            isPanning = true;
            startPanX = e.touches[0].clientX - panX;
            startPanY = e.touches[0].clientY - panY;
        } else if (e.touches.length === 2) {
            isPanning = false;
            lastTouchDistance = getTouchDistance(e.touches);
            lastTouchCenter = getTouchCenter(e.touches);
        }
    }, { passive: true });

    viewport.addEventListener('touchmove', (e) => {
        if (currentView !== 'graph') return;
        e.preventDefault();
        
        if (e.touches.length === 1 && isPanning) {
            panX = e.touches[0].clientX - startPanX;
            panY = e.touches[0].clientY - startPanY;
            applyTransform();
        } else if (e.touches.length === 2) {
            const distance = getTouchDistance(e.touches);
            const center = getTouchCenter(e.touches);
            
            if (lastTouchDistance > 0) {
                const scale = distance / lastTouchDistance;
                const newZoom = Math.min(MAX_ZOOM, Math.max(MIN_ZOOM, currentZoom * scale));
                
                if (newZoom !== currentZoom) {
                    const rect = viewport.getBoundingClientRect();
                    const centerX = center.x - rect.left;
                    const centerY = center.y - rect.top;
                    
                    const zoomScale = newZoom / currentZoom;
                    panX = centerX - (centerX - panX) * zoomScale;
                    panY = centerY - (centerY - panY) * zoomScale;
                    currentZoom = newZoom;
                    applyTransform();
                    showZoomIndicator();
                }
            }
            
            lastTouchDistance = distance;
            lastTouchCenter = center;
        }
    }, { passive: false });

    viewport.addEventListener('touchend', () => {
        isPanning = false;
        lastTouchDistance = 0;
    });
}

function getTouchDistance(touches) {
    const dx = touches[0].clientX - touches[1].clientX;
    const dy = touches[0].clientY - touches[1].clientY;
    return Math.sqrt(dx * dx + dy * dy);
}

function getTouchCenter(touches) {
    return {
        x: (touches[0].clientX + touches[1].clientX) / 2,
        y: (touches[0].clientY + touches[1].clientY) / 2
    };
}

function applyTransform() {
    const container = document.getElementById('mermaidContainer');
    if (container) {
        container.style.transform = `translate(${panX}px, ${panY}px) scale(${currentZoom})`;
    }
}

function showZoomIndicator() {
    const indicator = document.getElementById('zoomIndicator');
    if (indicator) {
        indicator.textContent = `${Math.round(currentZoom * 100)}%`;
        indicator.classList.add('visible');
        
        if (zoomIndicatorTimeout) {
            clearTimeout(zoomIndicatorTimeout);
        }
        zoomIndicatorTimeout = setTimeout(() => {
            indicator.classList.remove('visible');
        }, 1500);
    }
}

// Load module name from server
// Returns true if local module exists, false otherwise
async function loadModuleName() {
    try {
        const response = await fetch('/api/module');
        const data = await response.json();
        document.getElementById('moduleName').textContent = data.name || '';
        
        // If no local module, set default value and auto-analyze
        if (!data.name) {
            const input = document.getElementById('remoteModuleInput');
            if (input && !input.value) {
                input.value = 'github.com/gogf/gf/v2';
                // Fetch versions first, then auto-analyze
                await fetchVersions();
                analyzeRemoteModule();
            }
            return false;
        }
        return true;
    } catch (e) {
        console.error('Failed to load module name:', e);
        return false;
    }
}

// Load packages list
async function loadPackages() {
    try {
        const internal = document.getElementById('internal').checked;
        const external = document.getElementById('external') ? document.getElementById('external').checked : false;
        const moduleLevel = document.getElementById('moduleLevel') ? document.getElementById('moduleLevel').checked : false;
        const directOnly = document.getElementById('directOnly') ? document.getElementById('directOnly').checked : false;
        const response = await fetch(`/api/packages?internal=${internal}&external=${external}&module=${moduleLevel}&direct=${directOnly}`);
        const data = await response.json();
        
        // Handle new API response format with packages and statistics
        if (data.packages && Array.isArray(data.packages)) {
            allPackages = data.packages;
            // Update statistics display if available
            if (data.statistics) {
                updateStatisticsDisplay(data.statistics);
            }
        } else if (Array.isArray(data)) {
            // Fallback for old format
            allPackages = data;
        } else {
            console.error('Unexpected API response format:', data);
            allPackages = [];
        }
        
        document.getElementById('packageCount').textContent = allPackages.length;
        renderPackageList(allPackages);
    } catch (e) {
        console.error('Failed to load packages:', e);
    }
}

// Update statistics display
function updateStatisticsDisplay(statistics) {
    if (statistics) {
        document.getElementById('internalCount').textContent = statistics.internal || 0;
        document.getElementById('externalCount').textContent = statistics.external || 0;
        document.getElementById('stdlibCount').textContent = statistics.stdlib || 0;
        
        // Update total count
        const totalCount = statistics.total || (statistics.internal + statistics.external + statistics.stdlib);
        document.getElementById('nodeCount').textContent = totalCount;
    }
}

// Get package name from package object or string
function getPkgName(pkg) {
    return typeof pkg === 'object' ? pkg.name : pkg;
}

// Set package list display mode
function setPackageListMode(mode) {
    packageListMode = mode;
    document.getElementById('modeFlat').classList.toggle('active', mode === 'flat');
    document.getElementById('modeTree').classList.toggle('active', mode === 'tree');
    
    const query = document.getElementById('searchInput').value.toLowerCase();
    const filtered = query ? allPackages.filter(pkg => getPkgName(pkg).toLowerCase().includes(query)) : allPackages;
    renderPackageList(filtered);
}

// Render package list in sidebar
function renderPackageList(packages) {
    const list = document.getElementById('packageList');
    if (packages.length === 0) {
        list.innerHTML = '<div class="loading">' + i18n.t('noPackages') + '</div>';
        return;
    }
    
    if (packageListMode === 'tree') {
        renderPackageTree(packages, list);
    } else {
        renderPackageFlat(packages, list);
    }
}

// Render flat package list
function renderPackageFlat(packages, container) {
    container.innerHTML = packages.map(pkg => {
        const name = getPkgName(pkg);
        const isActive = name === selectedPackage ? ' active' : '';
        const escaped = name.replace(/'/g, "\\'");
        
        // Build stats display
        let statsHtml = '';
        if (typeof pkg === 'object') {
            statsHtml = `<span class="pkg-stats-inline">
                <span class="dep-count" title="${i18n.t('dependencies')}">→${pkg.depCount}</span>
                <span class="used-count" title="${i18n.t('usedBy')}">←${pkg.usedByCount}</span>
            </span>`;
        }
        
        return `<div class="package-item${isActive}" onclick="selectPackage('${escaped}')" title="${name}">
            <span class="pkg-name-text">${name}</span>${statsHtml}
        </div>`;
    }).join('');
}

// Build tree structure from package paths
function buildPackageTree(packages) {
    const root = { children: {}, packages: [] };
    
    packages.forEach(pkg => {
        const name = getPkgName(pkg);
        const parts = name.split('/');
        let current = root;
        
        parts.forEach((part, index) => {
            if (!current.children[part]) {
                current.children[part] = { children: {}, packages: [], path: parts.slice(0, index + 1).join('/') };
            }
            current = current.children[part];
        });
        current.isPackage = true;
        current.fullPath = name;
        // Store stats if available
        if (typeof pkg === 'object') {
            current.depCount = pkg.depCount;
            current.usedByCount = pkg.usedByCount;
        }
    });
    
    return root;
}

// Render package tree
function renderPackageTree(packages, container) {
    const tree = buildPackageTree(packages);
    container.innerHTML = renderTreeNode(tree, '');
}

// Render a tree node recursively
function renderTreeNode(node, path) {
    const children = Object.keys(node.children).sort();
    if (children.length === 0) return '';
    
    return children.map(name => {
        const child = node.children[name];
        const childPath = path ? `${path}/${name}` : name;
        const hasChildren = Object.keys(child.children).length > 0;
        const isExpanded = expandedNodes.has(childPath);
        const isActive = child.fullPath === selectedPackage;
        const isPackage = child.isPackage;
        
        const toggleClass = hasChildren ? (isExpanded ? 'expanded' : '') : 'empty';
        const activeClass = isActive ? ' active' : '';
        const nameClass = isPackage ? ' package' : '';
        const icon = isPackage ? '📦' : '📁';
        
        // Build stats for packages
        let statsHtml = '';
        if (isPackage && child.depCount !== undefined) {
            statsHtml = `<span class="pkg-stats-inline">
                <span class="dep-count" title="${i18n.t('dependencies')}">→${child.depCount}</span>
                <span class="used-count" title="${i18n.t('usedBy')}">←${child.usedByCount}</span>
            </span>`;
        }
        
        let html = `
            <div class="tree-node" data-path="${childPath}">
                <div class="tree-node-header${activeClass}">
                    <span class="tree-node-toggle ${toggleClass}" onclick="handleToggleClick(event, '${childPath.replace(/'/g, "\\'")}', ${hasChildren})">▶</span>
                    <span class="tree-node-icon">${icon}</span>
                    <span class="tree-node-name${nameClass}" onclick="handleNameClick(event, '${childPath.replace(/'/g, "\\'")}', ${isPackage}, ${hasChildren})">${name}</span>
                    ${statsHtml}
                </div>`;
        
        if (hasChildren) {
            const childrenHtml = renderTreeNode(child, childPath);
            html += `<div class="tree-node-children${isExpanded ? ' expanded' : ''}">${childrenHtml}</div>`;
        }
        
        html += '</div>';
        return html;
    }).join('');
}

// Handle toggle arrow click - expand/collapse
function handleToggleClick(event, path, hasChildren) {
    event.stopPropagation();
    if (hasChildren) {
        toggleTreeNode(path);
    }
}

// Handle name click - select package or toggle if folder
function handleNameClick(event, path, isPackage, hasChildren) {
    event.stopPropagation();
    if (isPackage) {
        selectPackage(path);
    } else if (hasChildren) {
        toggleTreeNode(path);
    }
}

// Toggle tree node expansion
function toggleTreeNode(path) {
    if (expandedNodes.has(path)) {
        expandedNodes.delete(path);
    } else {
        expandedNodes.add(path);
    }
    
    // Re-render with current filter
    const query = document.getElementById('searchInput').value.toLowerCase();
    const filtered = query ? allPackages.filter(pkg => getPkgName(pkg).toLowerCase().includes(query)) : allPackages;
    renderPackageList(filtered);
}

// Filter packages by search query
function filterPackages() {
    const query = document.getElementById('searchInput').value.toLowerCase();
    const filtered = allPackages.filter(pkg => getPkgName(pkg).toLowerCase().includes(query));
    renderPackageList(filtered);
}

// Select a package
async function selectPackage(pkg) {
    selectedPackage = pkg;
    
    // Update flat list items
    document.querySelectorAll('.package-item').forEach(el => {
        el.classList.toggle('active', el.textContent === pkg);
    });
    
    // Update tree items
    document.querySelectorAll('.tree-node-header').forEach(el => {
        const node = el.closest('.tree-node');
        el.classList.toggle('active', node && node.dataset.path === pkg);
    });
    
    await refresh();
}

// Clear package selection
function clearSelection() {
    selectedPackage = null;
    document.querySelectorAll('.package-item').forEach(el => el.classList.remove('active'));
    document.querySelectorAll('.tree-node-header').forEach(el => el.classList.remove('active'));
    closePackageInfo();
    refresh();
}

// Close package info sidebar
function closePackageInfo() {
    document.getElementById('packageInfo').classList.remove('visible');
}

// Set layout direction
function setLayout(layout) {
    currentLayout = layout;
    document.getElementById('layoutTD').classList.toggle('active', layout === 'TD');
    document.getElementById('layoutLR').classList.toggle('active', layout === 'LR');
    refresh();
}

// Set view mode
function setView(view) {
    currentView = view;
    document.getElementById('btnGraph').classList.toggle('active', view === 'graph');
    document.getElementById('btnTree').classList.toggle('active', view === 'tree');
    document.getElementById('btnList').classList.toggle('active', view === 'list');
    document.getElementById('zoomControls').classList.toggle('hidden', view !== 'graph');
    document.getElementById('layoutGroup').classList.toggle('hidden', view !== 'graph');
    refresh();
}

// Main refresh function
async function refresh() {
    const depth = document.getElementById('depth').value;
    const group = document.getElementById('group').checked;
    const reverse = document.getElementById('reverse').checked;
    const internal = document.getElementById('internal').checked;
    const external = document.getElementById('external') ? document.getElementById('external').checked : false;
    const moduleLevel = document.getElementById('moduleLevel') ? document.getElementById('moduleLevel').checked : false;
    const directOnly = document.getElementById('directOnly') ? document.getElementById('directOnly').checked : false;

    if (selectedPackage) {
        await showPackageInfo(selectedPackage);
    } else {
        closePackageInfo();
    }

    if (currentView === 'graph') {
        await refreshGraph(depth, group, reverse, internal, external, moduleLevel, directOnly);
    } else if (currentView === 'tree') {
        await refreshTree(depth, internal, external, moduleLevel, directOnly);
    } else {
        await refreshList(internal, external, moduleLevel, directOnly);
    }
}

// Show package info panel
async function showPackageInfo(pkg) {
    try {
        const response = await fetch('/api/package?name=' + encodeURIComponent(pkg));
        const info = await response.json();
        
        const infoDiv = document.getElementById('packageInfo');
        const contentDiv = document.getElementById('packageInfoContent');
        
        infoDiv.classList.add('visible');
        contentDiv.innerHTML = `
            <span class="pkg-name" title="${info.name}">${info.name}</span>
            <div class="pkg-stats">
                <div class="pkg-stat">
                    <span class="pkg-stat-label">${i18n.t('dependencies')}:</span>
                    <span class="pkg-stat-value">${info.dependencies.length}</span>
                </div>
                <div class="pkg-stat">
                    <span class="pkg-stat-label">${i18n.t('usedBy')}:</span>
                    <span class="pkg-stat-value">${info.usedBy.length}</span>
                </div>
            </div>
        `;
    } catch (e) {
        console.error('Failed to load package info:', e);
    }
}

// Refresh graph view
async function refreshGraph(depth, group, reverse, internal, external, moduleLevel, directOnly) {
    document.getElementById('graphView').classList.remove('hidden');
    document.getElementById('textView').classList.add('hidden');

    // Reset pan/zoom on refresh
    currentZoom = 1;
    panX = 0;
    panY = 0;
    applyTransform();

    let url = `/api/graph?depth=${depth}&group=${group}&reverse=${reverse}&internal=${internal}`;
    if (external !== undefined) {
        url += `&external=${external}`;
    }
    if (moduleLevel !== undefined) {
        url += `&module=${moduleLevel}`;
    }
    if (directOnly !== undefined) {
        url += `&direct=${directOnly}`;
    }
    if (selectedPackage) {
        url += '&package=' + encodeURIComponent(selectedPackage);
    }

    try {
        const response = await fetch(url);
        const data = await response.json();

        document.getElementById('nodeCount').textContent = data.nodes.length;
        document.getElementById('edgeCount').textContent = data.edges.length;

        // Generate mermaid code
        // Build node id to label map
        const nodeLabels = {};
        data.nodes.forEach(node => {
            nodeLabels[node.id] = node.label;
        });

        // Use current layout direction
        let mermaidCode = `graph ${currentLayout}\n`;
        // First define all nodes with labels
        data.nodes.forEach(node => {
            mermaidCode += `    ${node.id}["${node.label}"]\n`;
        });
        // Then add edges
        data.edges.forEach(edge => {
            mermaidCode += `    ${edge.from} --> ${edge.to}\n`;
        });

        const container = document.getElementById('mermaidGraph');
        container.innerHTML = mermaidCode;
        container.removeAttribute('data-processed');
        
        await mermaid.run({ nodes: [container] });
        
        // Auto-fit if graph is large
        setTimeout(() => {
            autoFitGraph();
        }, 100);
    } catch (e) {
        console.error('Failed to render graph:', e);
        document.getElementById('mermaidGraph').innerHTML = 
            `<div class="loading">${i18n.t('renderError')}</div>`;
    }
}

// Auto-fit graph to viewport
function autoFitGraph() {
    const viewport = document.getElementById('graphView');
    const svg = document.querySelector('.mermaid svg');
    
    if (!viewport || !svg) return;
    
    const viewportRect = viewport.getBoundingClientRect();
    const svgRect = svg.getBoundingClientRect();
    
    if (svgRect.width === 0 || svgRect.height === 0) return;
    
    // Calculate scale to fit
    const scaleX = (viewportRect.width - 40) / svgRect.width;
    const scaleY = (viewportRect.height - 40) / svgRect.height;
    const scale = Math.min(scaleX, scaleY, 1); // Don't zoom in beyond 100%
    
    if (scale < 1) {
        currentZoom = scale;
        // Center the graph
        panX = (viewportRect.width - svgRect.width * scale) / 2;
        panY = 20;
        applyTransform();
        showZoomIndicator();
    }
}

// Refresh tree view
async function refreshTree(depth, internal, external, moduleLevel, directOnly) {
    document.getElementById('graphView').classList.add('hidden');
    document.getElementById('textView').classList.remove('hidden');

    let url = `/api/tree?depth=${depth}&internal=${internal}`;
    if (external !== undefined) {
        url += `&external=${external}`;
    }
    if (moduleLevel !== undefined) {
        url += `&module=${moduleLevel}`;
    }
    if (directOnly !== undefined) {
        url += `&direct=${directOnly}`;
    }
    if (selectedPackage) {
        url += '&package=' + encodeURIComponent(selectedPackage);
    }

    try {
        const response = await fetch(url);
        const text = await response.text();
        document.getElementById('textView').textContent = text;

        const lines = text.split('\n').filter(l => l.trim());
        document.getElementById('nodeCount').textContent = lines.length;
        document.getElementById('edgeCount').textContent = '-';
    } catch (e) {
        console.error('Failed to load tree:', e);
    }
}

// Refresh list view
async function refreshList(internal, external, moduleLevel, directOnly) {
    document.getElementById('graphView').classList.add('hidden');
    document.getElementById('textView').classList.remove('hidden');

    let url = `/api/list?internal=${internal}`;
    if (external !== undefined) {
        url += `&external=${external}`;
    }
    if (moduleLevel !== undefined) {
        url += `&module=${moduleLevel}`;
    }
    if (directOnly !== undefined) {
        url += `&direct=${directOnly}`;
    }
    if (selectedPackage) {
        url += '&package=' + encodeURIComponent(selectedPackage);
    }

    try {
        const response = await fetch(url);
        const text = await response.text();
        document.getElementById('textView').textContent = text;

        const lines = text.split('\n').filter(l => l.trim());
        document.getElementById('nodeCount').textContent = lines.length;
        document.getElementById('edgeCount').textContent = '-';
    } catch (e) {
        console.error('Failed to load list:', e);
    }
}

// Zoom functions
function zoomIn() {
    const viewport = document.getElementById('graphView');
    if (!viewport) return;
    
    const rect = viewport.getBoundingClientRect();
    const centerX = rect.width / 2;
    const centerY = rect.height / 2;
    
    const newZoom = Math.min(MAX_ZOOM, currentZoom * 1.2);
    const scale = newZoom / currentZoom;
    
    panX = centerX - (centerX - panX) * scale;
    panY = centerY - (centerY - panY) * scale;
    currentZoom = newZoom;
    
    applyTransform();
    showZoomIndicator();
}

function zoomOut() {
    const viewport = document.getElementById('graphView');
    if (!viewport) return;
    
    const rect = viewport.getBoundingClientRect();
    const centerX = rect.width / 2;
    const centerY = rect.height / 2;
    
    const newZoom = Math.max(MIN_ZOOM, currentZoom / 1.2);
    const scale = newZoom / currentZoom;
    
    panX = centerX - (centerX - panX) * scale;
    panY = centerY - (centerY - panY) * scale;
    currentZoom = newZoom;
    
    applyTransform();
    showZoomIndicator();
}

function resetZoom() {
    currentZoom = 1;
    panX = 0;
    panY = 0;
    applyTransform();
    showZoomIndicator();
}

function fitToScreen() {
    autoFitGraph();
}

// Initialize when DOM is ready
document.addEventListener('DOMContentLoaded', init);
