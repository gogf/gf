// Internationalization (i18n) Module
const i18n = {
    currentLang: 'en',
    
    translations: {
        en: {
            title: 'Go Package Dependencies',
            pageTitle: 'Package Dependency Viewer',
            currentModule: 'Current:',
            remoteModuleLabel: 'Analyze Module:',
            remoteModulePlaceholder: 'e.g. github.com/gogf/gf/v2',
            selectVersion: 'Version',
            analyze: 'Analyze',
            resetLocal: 'Reset',
            fetchingVersions: 'Fetching...',
            analyzing: 'Analyzing...',
            viewLabel: 'View:',
            viewGraph: 'Graph',
            viewTree: 'Tree',
            viewList: 'List',
            depthLabel: 'Depth:',
            depthUnlimited: 'Unlimited',
            reverseLabel: 'Reverse (show who uses)',
            groupLabel: 'Group by directory',
            layoutLabel: 'Layout:',
            layoutTD: 'Top-Down',
            layoutLR: 'Left-Right',
            showAll: 'Show All',
            packagesTitle: 'Packages',
            searchPlaceholder: 'Search packages...',
            flatMode: 'Flat list',
            treeMode: 'Tree view',
            statsPackages: 'Packages',
            statsDependencies: 'Dependencies',
            dependencies: 'Dependencies',
            usedBy: 'Used by',
            noPackages: 'No packages found',
            renderError: 'Unable to render graph. Try reducing depth or selecting a specific package.',
            packageNotFound: 'Package not found',
            zoomIn: 'Zoom in',
            zoomOut: 'Zoom out',
            fitToScreen: 'Fit to screen',
            resetZoom: 'Reset zoom',
            dragToMove: 'Drag to move, scroll to zoom',
            noVersions: 'No versions found',
            latestVersion: 'Latest',
            errorFetchVersions: 'Failed to fetch versions',
            errorAnalyze: 'Failed to analyze module',
            packageDetails: 'Package Details'
        },
        zh: {
            title: 'Go 包依赖分析',
            pageTitle: '包依赖查看器',
            currentModule: '当前模块:',
            remoteModuleLabel: '分析模块:',
            remoteModulePlaceholder: '例如 github.com/gogf/gf/v2',
            selectVersion: '版本',
            analyze: '分析',
            resetLocal: '重置',
            fetchingVersions: '获取中...',
            analyzing: '分析中...',
            viewLabel: '视图:',
            viewGraph: '图表',
            viewTree: '树形',
            viewList: '列表',
            depthLabel: '深度:',
            depthUnlimited: '无限制',
            reverseLabel: '反向依赖 (谁引用了它)',
            groupLabel: '按目录分组',
            layoutLabel: '布局:',
            layoutTD: '从上到下',
            layoutLR: '从左到右',
            showAll: '显示全部',
            packagesTitle: '包列表',
            searchPlaceholder: '搜索包...',
            flatMode: '平铺列表',
            treeMode: '目录树',
            statsPackages: '包数量',
            statsDependencies: '依赖数',
            dependencies: '依赖',
            usedBy: '被引用',
            noPackages: '未找到包',
            renderError: '无法渲染图表。请尝试减少深度或选择特定的包。',
            packageNotFound: '未找到包',
            zoomIn: '放大',
            zoomOut: '缩小',
            fitToScreen: '适应屏幕',
            resetZoom: '重置缩放',
            dragToMove: '拖拽移动，滚轮缩放',
            noVersions: '未找到版本',
            latestVersion: '最新版本',
            errorFetchVersions: '获取版本失败',
            errorAnalyze: '分析模块失败',
            packageDetails: '包详情'
        }
    },
    
    init() {
        // Load saved language preference
        const savedLang = localStorage.getItem('dep-viewer-lang');
        if (savedLang && this.translations[savedLang]) {
            this.currentLang = savedLang;
        } else {
            // Detect browser language
            const browserLang = navigator.language.toLowerCase();
            if (browserLang.startsWith('zh')) {
                this.currentLang = 'zh';
            }
        }
        
        // Set select value
        const langSelect = document.getElementById('langSelect');
        if (langSelect) {
            langSelect.value = this.currentLang;
            langSelect.addEventListener('change', (e) => {
                this.setLanguage(e.target.value);
            });
        }
        
        this.applyTranslations();
    },
    
    setLanguage(lang) {
        if (this.translations[lang]) {
            this.currentLang = lang;
            localStorage.setItem('dep-viewer-lang', lang);
            this.applyTranslations();
        }
    },
    
    t(key) {
        return this.translations[this.currentLang][key] || this.translations['en'][key] || key;
    },
    
    applyTranslations() {
        // Update elements with data-i18n attribute
        document.querySelectorAll('[data-i18n]').forEach(el => {
            const key = el.getAttribute('data-i18n');
            el.textContent = this.t(key);
        });
        
        // Update placeholders
        document.querySelectorAll('[data-i18n-placeholder]').forEach(el => {
            const key = el.getAttribute('data-i18n-placeholder');
            el.placeholder = this.t(key);
        });
        
        // Update page title
        document.title = this.t('title');
    }
};

// Initialize i18n when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    i18n.init();
});
