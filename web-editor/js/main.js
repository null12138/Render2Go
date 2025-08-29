// 主应用程序 - main.js

// 应用程序主类
class App extends EventEmitter {
    constructor() {
        super();
        this.scene = null;
        this.objectManager = null;
        this.timeline = null;
        this.propertiesPanel = null;
        this.exportManager = null;
        this.undoManager = null;

        this.isLoaded = false;
        this.projectModified = false;
        this.currentProjectFile = null;

        this.init();
    }

    init() {
        // 等待DOM加载完成
        if (document.readyState === 'loading') {
            document.addEventListener('DOMContentLoaded', () => this.startup());
        } else {
            this.startup();
        }
    }

    startup() {
        try {
            // 初始化核心组件
            this.initializeComponents();

            // 设置事件监听器
            this.setupEventListeners();

            // 初始化UI
            this.initializeUI();

            // 标记为已加载
            this.isLoaded = true;
            this.emit('ready');

            console.log('Render2Go Web Editor initialized successfully');

        } catch (error) {
            console.error('Failed to initialize application:', error);
            this.showError('初始化失败', error.message);
        }
    }

    initializeComponents() {
        // 撤销管理器
        this.undoManager = new UndoRedoManager();

        // 对象管理器
        this.objectManager = new ObjectManager();

        // 场景管理器
        this.scene = new Scene();

        // 时间轴
        this.timeline = new Timeline();

        // 属性面板
        this.propertiesPanel = new PropertiesPanel();

        // 导出管理器
        this.exportManager = new ExportManager();

        // 设置全局引用
        window.app = this;

        // 重新绑定需要app引用的组件事件
        setTimeout(() => {
            if (this.propertiesPanel && this.propertiesPanel.bindObjectManagerEvents) {
                this.propertiesPanel.bindObjectManagerEvents();
            }
        }, 10);
    }

    setupEventListeners() {
        // 工具栏事件
        this.setupToolbarEvents();

        // 菜单事件
        this.setupMenuEvents();

        // 文件操作事件
        this.setupFileEvents();

        // 窗口事件
        this.setupWindowEvents();

        // 键盘快捷键
        this.setupKeyboardShortcuts();

        // 组件间通信
        this.setupComponentEvents();
    }

    setupToolbarEvents() {
        // 工具按钮
        document.querySelectorAll('.tool-btn').forEach(btn => {
            btn.addEventListener('click', () => {
                const tool = btn.dataset.tool;
                if (tool) {
                    this.setActiveTool(tool);
                    this.scene.setTool(tool);
                }
            });
        });

        // 播放控制
        const playBtn = document.getElementById('playBtn');
        const pauseBtn = document.getElementById('pauseBtn');
        const stopBtn = document.getElementById('stopBtn');

        if (playBtn) playBtn.addEventListener('click', () => this.timeline.play());
        if (pauseBtn) pauseBtn.addEventListener('click', () => this.timeline.pause());
        if (stopBtn) stopBtn.addEventListener('click', () => this.timeline.stop());

        // 缩放控制
        const zoomInBtn = document.getElementById('zoomInBtn');
        const zoomOutBtn = document.getElementById('zoomOutBtn');
        const fitScreenBtn = document.getElementById('fitScreenBtn');

        if (zoomInBtn) zoomInBtn.addEventListener('click', () => this.scene.zoomIn());
        if (zoomOutBtn) zoomOutBtn.addEventListener('click', () => this.scene.zoomOut());
        if (fitScreenBtn) fitScreenBtn.addEventListener('click', () => this.scene.fitToScreen());
    }

    setupMenuEvents() {
        // 文件菜单
        const newBtn = document.getElementById('newProject');
        const openBtn = document.getElementById('openProject');
        const saveBtn = document.getElementById('saveProject');
        const exportBtn = document.getElementById('exportScript');

        if (newBtn) newBtn.addEventListener('click', () => this.newProject());
        if (openBtn) openBtn.addEventListener('click', () => this.openProject());
        if (saveBtn) saveBtn.addEventListener('click', () => this.saveProject());
        if (exportBtn) exportBtn.addEventListener('click', () => this.exportManager.showExportSettings());

        // 编辑菜单
        const undoBtn = document.getElementById('undoBtn');
        const redoBtn = document.getElementById('redoBtn');
        const copyBtn = document.getElementById('copyBtn');
        const pasteBtn = document.getElementById('pasteBtn');
        const deleteBtn = document.getElementById('deleteBtn');

        if (undoBtn) undoBtn.addEventListener('click', () => this.undoManager.undo());
        if (redoBtn) redoBtn.addEventListener('click', () => this.undoManager.redo());
        if (copyBtn) copyBtn.addEventListener('click', () => this.objectManager.copy());
        if (pasteBtn) pasteBtn.addEventListener('click', () => this.objectManager.paste());
        if (deleteBtn) deleteBtn.addEventListener('click', () => this.objectManager.deleteSelected());
    }

    setupFileEvents() {
        // 文件拖放
        this.setupDragAndDrop();

        // 文件输入
        const fileInput = document.getElementById('fileInput');
        if (fileInput) {
            fileInput.addEventListener('change', (e) => {
                if (e.target.files.length > 0) {
                    this.loadProjectFile(e.target.files[0]);
                }
            });
        }
    }

    setupDragAndDrop() {
        const dropZone = document.body;

        dropZone.addEventListener('dragover', (e) => {
            e.preventDefault();
            e.dataTransfer.dropEffect = 'copy';
            dropZone.classList.add('drag-over');
        });

        dropZone.addEventListener('dragleave', (e) => {
            if (!dropZone.contains(e.relatedTarget)) {
                dropZone.classList.remove('drag-over');
            }
        });

        dropZone.addEventListener('drop', (e) => {
            e.preventDefault();
            dropZone.classList.remove('drag-over');

            const files = Array.from(e.dataTransfer.files);
            const projectFile = files.find(file =>
                file.name.endsWith('.r2gp') || file.name.endsWith('.json'));

            if (projectFile) {
                this.loadProjectFile(projectFile);
            }
        });
    }

    setupWindowEvents() {
        // 窗口大小变化
        window.addEventListener('resize', Utils.debounce(() => {
            this.handleWindowResize();
        }, 250));

        // 页面卸载前确认
        window.addEventListener('beforeunload', (e) => {
            if (this.projectModified) {
                e.preventDefault();
                e.returnValue = '项目有未保存的更改，确定要离开吗？';
                return e.returnValue;
            }
        });

        // 全屏变化
        document.addEventListener('fullscreenchange', () => {
            this.updateFullscreenButton();
        });
    }

    setupKeyboardShortcuts() {
        document.addEventListener('keydown', (e) => {
            // 忽略在输入框中的按键
            if (e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA') {
                return;
            }

            // Ctrl/Cmd + 按键
            if (e.ctrlKey || e.metaKey) {
                switch (e.key.toLowerCase()) {
                    case 'n':
                        e.preventDefault();
                        this.newProject();
                        break;
                    case 'o':
                        e.preventDefault();
                        this.openProject();
                        break;
                    case 's':
                        e.preventDefault();
                        if (e.shiftKey) {
                            this.saveAsProject();
                        } else {
                            this.saveProject();
                        }
                        break;
                    case 'e':
                        e.preventDefault();
                        this.exportManager.showExportSettings();
                        break;
                    case 'z':
                        e.preventDefault();
                        if (e.shiftKey) {
                            this.undoManager.redo();
                        } else {
                            this.undoManager.undo();
                        }
                        break;
                    case 'y':
                        e.preventDefault();
                        this.undoManager.redo();
                        break;
                    case 'c':
                        e.preventDefault();
                        this.objectManager.copy();
                        break;
                    case 'v':
                        e.preventDefault();
                        this.objectManager.paste();
                        break;
                    case 'a':
                        e.preventDefault();
                        this.scene.selectAll();
                        break;
                    case 'd':
                        e.preventDefault();
                        this.objectManager.duplicate();
                        break;
                }
            }

            // 功能键
            switch (e.key) {
                case 'Delete':
                case 'Backspace':
                    this.objectManager.deleteSelected();
                    break;
                case 'F11':
                    e.preventDefault();
                    this.toggleFullscreen();
                    break;
                case ' ':
                    if (e.target === document.body) {
                        e.preventDefault();
                        this.timeline.isPlaying ? this.timeline.pause() : this.timeline.play();
                    }
                    break;
            }

            // 工具快捷键
            const toolKeys = {
                'v': 'select',
                'm': 'move',
                'r': 'rotate',
                's': 'scale',
                'c': 'circle',
                't': 'text',
                'l': 'line',
                'p': 'triangle'
            };

            if (!e.ctrlKey && !e.metaKey && toolKeys[e.key.toLowerCase()]) {
                const tool = toolKeys[e.key.toLowerCase()];
                this.setActiveTool(tool);
                this.scene.setTool(tool);
            }
        });
    }

    setupComponentEvents() {
        // 对象管理器事件
        this.objectManager.on('objectAdded', () => this.markModified());
        this.objectManager.on('objectRemoved', () => this.markModified());
        this.objectManager.on('objectModified', () => this.markModified());

        // 场景事件
        this.scene.on('sceneChanged', () => this.markModified());

        // 时间轴事件
        this.timeline.on('keyframeAdded', () => this.markModified());
        this.timeline.on('keyframeRemoved', () => this.markModified());

        // 撤销管理器事件
        this.undoManager.on('executed', () => this.updateUndoRedoButtons());
        this.undoManager.on('undone', () => this.updateUndoRedoButtons());
        this.undoManager.on('redone', () => this.updateUndoRedoButtons());
    }

    initializeUI() {
        // 初始化对象库
        this.initializeObjectLibrary();

        // 初始化面板大小
        this.initializePanelSizes();

        // 初始化工具提示
        this.initializeTooltips();

        // 初始化状态栏
        this.initializeStatusBar();

        // 初始化主题
        this.initializeTheme();

        // 设置默认值
        this.updateUI();
    }

    initializePanelSizes() {
        // 设置面板分割器
        this.setupResizablePanels();
    }

    setupResizablePanels() {
        // 左侧面板
        const leftSplitter = document.querySelector('.left-splitter');
        if (leftSplitter) {
            this.setupSplitter(leftSplitter, 'left');
        }

        // 右侧面板
        const rightSplitter = document.querySelector('.right-splitter');
        if (rightSplitter) {
            this.setupSplitter(rightSplitter, 'right');
        }

        // 底部面板
        const bottomSplitter = document.querySelector('.bottom-splitter');
        if (bottomSplitter) {
            this.setupSplitter(bottomSplitter, 'bottom');
        }
    }

    setupSplitter(splitter, direction) {
        let isResizing = false;
        let startPos = 0;
        let startSize = 0;

        splitter.addEventListener('mousedown', (e) => {
            isResizing = true;
            startPos = direction === 'bottom' ? e.clientY : e.clientX;

            const panel = splitter.previousElementSibling || splitter.nextElementSibling;
            if (panel) {
                const rect = panel.getBoundingClientRect();
                startSize = direction === 'bottom' ? rect.height : rect.width;
            }

            document.addEventListener('mousemove', onMouseMove);
            document.addEventListener('mouseup', onMouseUp);
            document.body.style.cursor = direction === 'bottom' ? 'ns-resize' : 'ew-resize';
        });

        const onMouseMove = (e) => {
            if (!isResizing) return;

            const currentPos = direction === 'bottom' ? e.clientY : e.clientX;
            const delta = currentPos - startPos;
            const newSize = startSize + (direction === 'right' ? -delta : delta);

            const panel = splitter.previousElementSibling || splitter.nextElementSibling;
            if (panel) {
                const minSize = 200;
                const maxSize = window.innerWidth * 0.8;
                const clampedSize = Math.max(minSize, Math.min(maxSize, newSize));

                if (direction === 'bottom') {
                    panel.style.height = clampedSize + 'px';
                } else {
                    panel.style.width = clampedSize + 'px';
                }
            }
        };

        const onMouseUp = () => {
            isResizing = false;
            document.removeEventListener('mousemove', onMouseMove);
            document.removeEventListener('mouseup', onMouseUp);
            document.body.style.cursor = '';
        };
    }

    initializeTooltips() {
        // 为所有有title属性的元素添加工具提示
        document.querySelectorAll('[title]').forEach(element => {
            this.addTooltip(element);
        });
    }

    addTooltip(element) {
        const title = element.getAttribute('title');
        if (!title) return;

        element.removeAttribute('title'); // 移除默认title

        let tooltip = null;

        element.addEventListener('mouseenter', () => {
            tooltip = Utils.createElement('div', 'tooltip');
            tooltip.textContent = title;
            document.body.appendChild(tooltip);

            const rect = element.getBoundingClientRect();
            tooltip.style.left = rect.left + rect.width / 2 + 'px';
            tooltip.style.top = rect.bottom + 5 + 'px';
        });

        element.addEventListener('mouseleave', () => {
            if (tooltip) {
                tooltip.remove();
                tooltip = null;
            }
        });
    }

    initializeStatusBar() {
        this.updateStatusBar();
    }

    initializeTheme() {
        const savedTheme = localStorage.getItem('render2go-theme') || 'dark';
        this.setTheme(savedTheme);
    }

    initializeObjectLibrary() {
        console.log('Initializing object library...');
        const objectsGrid = document.getElementById('objectsGrid');
        if (!objectsGrid) {
            console.error('objectsGrid element not found!');
            return;
        }

        // 获取对象类型
        const objectTypes = ObjectFactory.getObjectTypes();
        console.log('Object types:', objectTypes);

        // 清空现有内容
        objectsGrid.innerHTML = '';

        // 为每种对象类型创建按钮
        objectTypes.forEach(type => {
            console.log(`Creating button for ${type.type}`);
            const objectItem = document.createElement('div');
            objectItem.className = 'object-item';
            objectItem.dataset.type = type.type;
            objectItem.innerHTML = `
                <i class="${type.icon}"></i>
                <span>${type.name}</span>
            `;

            // 添加点击事件
            objectItem.addEventListener('click', (e) => {
                console.log(`Clicked on ${type.type} button`);
                e.preventDefault();
                e.stopPropagation();
                this.createObject(type.type);
            });

            // 添加拖拽事件
            objectItem.draggable = true;
            objectItem.addEventListener('dragstart', (e) => {
                console.log(`Dragging ${type.type}`);
                e.dataTransfer.setData('text/plain', type.type);
                e.dataTransfer.effectAllowed = 'copy';
            });

            objectsGrid.appendChild(objectItem);
        });

        console.log('Object library initialized successfully');
    }

    createObject(type) {
        try {
            console.log(`Starting to create ${type} object...`);
            const object = ObjectFactory.createObject(type);
            console.log('Object created:', object);

            // 设置初始位置（画布中心）
            const canvasRect = this.scene.canvas.getBoundingClientRect();
            const centerX = this.scene.width / 2 - object.width / 2;
            const centerY = this.scene.height / 2 - object.height / 2;

            console.log(`Setting position to: ${centerX}, ${centerY}`);
            object.setPosition(centerX, centerY);

            // 添加到场景
            console.log('Adding object to scene...');
            this.objectManager.addObject(object);

            // 选择新创建的对象
            console.log('Selecting object...');
            this.objectManager.clearSelection();
            this.objectManager.selectObject(object.id);

            console.log(`Successfully created ${type} object:`, object);
        } catch (error) {
            console.error(`Failed to create ${type} object:`, error);
            this.showError('创建对象失败', error.message);
        }
    }

    setActiveTool(tool) {
        // 更新工具按钮状态
        document.querySelectorAll('.tool-btn').forEach(btn => {
            btn.classList.remove('active');
            if (btn.dataset.tool === tool) {
                btn.classList.add('active');
            }
        });

        // 更新光标样式
        const canvas = document.getElementById('sceneCanvas');
        if (canvas) {
            canvas.className = canvas.className.replace(/tool-\w+/g, '');
            canvas.classList.add(`tool-${tool}`);
        }
    }

    showError(title, message) {
        console.error(`${title}: ${message}`);
        // 简单的错误显示，后续可以改为更好的UI
        alert(`${title}\n${message}`);
    }

    showSuccess(message) {
        console.log(message);
        // 简单的成功提示，后续可以改为更好的UI
    }

    // 项目管理
    newProject() {
        if (this.projectModified) {
            if (!confirm('当前项目有未保存的更改，确定要创建新项目吗？')) {
                return;
            }
        }

        this.scene.clear();
        this.timeline.clear();
        this.currentProjectFile = null;
        this.projectModified = false;

        this.updateTitle();
        this.updateUI();

        this.emit('projectNew');
    }

    openProject() {
        const input = document.getElementById('fileInput');
        if (input) {
            input.click();
        }
    }

    async loadProjectFile(file) {
        try {
            const projectData = await this.exportManager.loadProject(file);
            this.currentProjectFile = file.name;
            this.projectModified = false;

            this.updateTitle();
            this.updateUI();

            this.showNotification('项目加载成功', 'success');

        } catch (error) {
            this.showError('加载项目失败', error.message);
        }
    }

    saveProject() {
        if (this.currentProjectFile) {
            this.exportManager.saveProject(this.currentProjectFile);
        } else {
            this.saveAsProject();
        }
    }

    saveAsProject() {
        const filename = prompt('请输入项目名称:', this.scene.name || 'my_project');
        if (filename) {
            this.exportManager.saveProject(filename);
            this.currentProjectFile = filename;
            this.projectModified = false;
            this.updateTitle();
        }
    }

    // UI更新
    updateUI() {
        this.updateUndoRedoButtons();
        this.updateStatusBar();
        this.updateTitle();
    }

    updateUndoRedoButtons() {
        const undoBtn = document.getElementById('undoBtn');
        const redoBtn = document.getElementById('redoBtn');

        if (undoBtn) {
            undoBtn.disabled = !this.undoManager.canUndo();
        }

        if (redoBtn) {
            redoBtn.disabled = !this.undoManager.canRedo();
        }
    }

    updateStatusBar() {
        const statusBar = document.querySelector('.status-bar');
        if (!statusBar) return;

        const selectedCount = this.objectManager.getSelectedObjects().length;
        const totalCount = this.objectManager.getAllObjects().length;

        statusBar.textContent = `对象: ${totalCount} | 选中: ${selectedCount} | 时间: ${this.timeline.currentTime.toFixed(2)}s`;
    }

    updateTitle() {
        const title = this.currentProjectFile || '未命名项目';
        const modified = this.projectModified ? ' *' : '';
        document.title = `${title}${modified} - Render2Go Web Editor`;
    }

    markModified() {
        if (!this.projectModified) {
            this.projectModified = true;
            this.updateTitle();
        }
    }

    // 主题管理
    setTheme(theme) {
        document.body.className = document.body.className.replace(/theme-\w+/g, '');
        document.body.classList.add(`theme-${theme}`);
        localStorage.setItem('render2go-theme', theme);
    }

    // 全屏
    toggleFullscreen() {
        if (document.fullscreenElement) {
            document.exitFullscreen();
        } else {
            document.documentElement.requestFullscreen();
        }
    }

    updateFullscreenButton() {
        const btn = document.getElementById('fullscreenBtn');
        if (btn) {
            btn.classList.toggle('active', document.fullscreenElement !== null);
        }
    }

    // 响应式处理
    handleWindowResize() {
        // 更新画布大小
        if (this.scene && this.scene.canvas) {
            this.scene.updateCanvasSize();
        }

        // 更新面板大小
        this.updateStatusBar();
    }

    // 通知系统
    showNotification(message, type = 'info', duration = 3000) {
        const notification = Utils.createElement('div', `notification notification-${type}`);
        notification.textContent = message;

        const container = document.querySelector('.notifications') || this.createNotificationContainer();
        container.appendChild(notification);

        // 自动消失
        setTimeout(() => {
            notification.classList.add('fade-out');
            setTimeout(() => notification.remove(), 300);
        }, duration);
    }

    createNotificationContainer() {
        const container = Utils.createElement('div', 'notifications');
        document.body.appendChild(container);
        return container;
    }

    showError(title, message) {
        alert(`${title}\n\n${message}`);
    }

    // 帮助和关于
    showHelp() {
        window.open('https://github.com/render2go/web-editor/wiki', '_blank');
    }

    showAbout() {
        alert('Render2Go Web Editor\n\n一个可视化的动画脚本编辑器\n\n版本: 1.0.0');
    }
}

// 初始化应用程序
window.addEventListener('DOMContentLoaded', () => {
    window.app = new App();
});

// 导出
window.App = App;
