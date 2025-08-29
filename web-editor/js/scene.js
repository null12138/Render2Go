// 场景管理系统 - scene.js

// 场景类
class Scene extends EventEmitter {
    constructor() {
        super();
        this.width = 800;
        this.height = 600;
        this.name = 'my_animation';
        this.frameRate = 30;
        this.duration = 5; // 秒
        this.backgroundColor = '#ffffff';
        this.currentTool = 'select';
        this.zoom = 1.0;
        this.grid = {
            visible: true,
            size: 20,
            color: '#e2e8f0'
        };

        // DOM元素
        this.canvas = null;
        this.container = null;
        this.viewport = null;

        // 交互状态
        this.isDragging = false;
        this.dragStartPos = { x: 0, y: 0 };
        this.dragStartObjectPos = { x: 0, y: 0 };
        this.draggedObject = null;
        this.selectionBox = null;

        this.init();
    }

    init() {
        this.canvas = document.getElementById('sceneCanvas');
        this.container = document.getElementById('objectsContainer');
        this.viewport = document.getElementById('sceneViewport');

        if (this.canvas) {
            this.setupEventListeners();
            this.updateCanvasSize();
            this.updateGrid();
        }
    }

    setupEventListeners() {
        // 鼠标事件
        this.canvas.addEventListener('mousedown', this.onMouseDown.bind(this));
        this.canvas.addEventListener('mousemove', this.onMouseMove.bind(this));
        this.canvas.addEventListener('mouseup', this.onMouseUp.bind(this));
        this.canvas.addEventListener('click', this.onClick.bind(this));
        this.canvas.addEventListener('contextmenu', this.onContextMenu.bind(this));

        // 键盘事件
        document.addEventListener('keydown', this.onKeyDown.bind(this));

        // 缩放事件
        this.canvas.addEventListener('wheel', this.onWheel.bind(this));

        // 拖放事件
        this.canvas.addEventListener('dragover', this.onDragOver.bind(this));
        this.canvas.addEventListener('drop', this.onDrop.bind(this));
    } onMouseDown(e) {
        e.preventDefault();
        const pos = this.screenToCanvasCoords(e.clientX, e.clientY);
        this.dragStartPos = pos;
        console.log('Mouse down at:', pos);

        // 检查是否点击了对象
        const hitObject = this.getObjectAtPosition(pos.x, pos.y);
        console.log('Hit object:', hitObject);

        if (hitObject && !hitObject.locked) {
            if (this.currentTool === 'select') {
                // 选择模式
                console.log('Starting drag for object:', hitObject.id);
                if (!e.ctrlKey && !e.shiftKey) {
                    app.objectManager.clearSelection();
                }
                app.objectManager.selectObject(hitObject.id, e.ctrlKey || e.shiftKey);

                this.isDragging = true;
                this.draggedObject = hitObject;
                this.dragStartObjectPos = { x: hitObject.x, y: hitObject.y };
                this.canvas.style.cursor = 'move';
            }
        } else {
            // 没有点击对象，开始框选
            console.log('No object hit, starting selection box');
            if (this.currentTool === 'select' && !e.ctrlKey) {
                app.objectManager.clearSelection();
                this.startSelectionBox(pos);
            }
        }
    }

    onMouseMove(e) {
        const pos = this.screenToCanvasCoords(e.clientX, e.clientY);

        if (this.isDragging && this.draggedObject) {
            // 拖拽对象
            console.log('Dragging object:', this.draggedObject.id, 'to position:', pos);
            const dx = pos.x - this.dragStartPos.x;
            const dy = pos.y - this.dragStartPos.y;

            const newX = this.dragStartObjectPos.x + dx;
            const newY = this.dragStartObjectPos.y + dy;

            // 网格吸附
            const snappedPos = this.snapToGrid(newX, newY);
            this.draggedObject.setPosition(snappedPos.x, snappedPos.y);

            // 更新其他选中对象
            const selectedObjects = app.objectManager.getSelectedObjects();
            selectedObjects.forEach(obj => {
                if (obj !== this.draggedObject) {
                    obj.setPosition(obj.x + dx, obj.y + dy);
                }
            });

        } else if (this.selectionBox) {
            // 更新选择框
            this.updateSelectionBox(pos);
        } else {
            // 更新鼠标样式
            const hitObject = this.getObjectAtPosition(pos.x, pos.y);
            this.canvas.style.cursor = hitObject && !hitObject.locked ? 'pointer' : 'default';
        }
    }

    onMouseUp(e) {
        if (this.isDragging) {
            this.isDragging = false;
            this.draggedObject = null;
            this.canvas.style.cursor = 'default';

            // 发送对象移动事件
            this.emit('objectMoved');
        }

        if (this.selectionBox) {
            this.endSelectionBox();
        }
    }

    onClick(e) {
        const pos = this.screenToCanvasCoords(e.clientX, e.clientY);

        // 根据当前工具处理点击
        switch (this.currentTool) {
            case 'circle':
                this.createObjectAt('circle', pos);
                break;
            case 'rectangle':
                this.createObjectAt('rectangle', pos);
                break;
            case 'text':
                this.createObjectAt('text', pos);
                break;
            case 'line':
                this.createObjectAt('line', pos);
                break;
            case 'triangle':
                this.createObjectAt('triangle', pos);
                break;
        }
    }

    onContextMenu(e) {
        e.preventDefault();
        const pos = this.screenToCanvasCoords(e.clientX, e.clientY);
        const hitObject = this.getObjectAtPosition(pos.x, pos.y);

        // 显示上下文菜单
        this.showContextMenu(e.clientX, e.clientY, hitObject);
    }

    onKeyDown(e) {
        // 只在编辑器聚焦时处理
        if (!document.activeElement || document.activeElement.tagName === 'INPUT') {
            return;
        }

        switch (e.key) {
            case 'Delete':
            case 'Backspace':
                app.objectManager.deleteSelected();
                break;
            case 'c':
                if (e.ctrlKey) {
                    app.objectManager.copy();
                }
                break;
            case 'v':
                if (e.ctrlKey) {
                    app.objectManager.paste();
                }
                break;
            case 'a':
                if (e.ctrlKey) {
                    e.preventDefault();
                    this.selectAll();
                }
                break;
            case 'z':
                if (e.ctrlKey && !e.shiftKey) {
                    app.undoManager.undo();
                } else if (e.ctrlKey && e.shiftKey) {
                    app.undoManager.redo();
                }
                break;
            case 'Escape':
                this.setTool('select');
                app.objectManager.clearSelection();
                break;
        }
    }

    onWheel(e) {
        e.preventDefault();
        const delta = e.deltaY > 0 ? 0.9 : 1.1;
        this.setZoom(this.zoom * delta);
    }

    // 坐标转换
    screenToCanvasCoords(screenX, screenY) {
        const rect = this.canvas.getBoundingClientRect();
        return {
            x: (screenX - rect.left) / this.zoom,
            y: (screenY - rect.top) / this.zoom
        };
    }

    // 工具方法
    setTool(tool) {
        this.currentTool = tool;
        this.canvas.className = `scene-canvas tool-${tool}`;
        this.emit('toolChanged', tool);

        // 更新工具按钮状态
        document.querySelectorAll('.tool-btn').forEach(btn => {
            btn.classList.toggle('active', btn.dataset.tool === tool);
        });
    }

    // 对象创建
    createObjectAt(type, pos) {
        const object = ObjectFactory.createObject(type);
        object.setPosition(pos.x - object.width / 2, pos.y - object.height / 2);

        this.container.appendChild(object.element);
        app.objectManager.addObject(object);
        app.objectManager.clearSelection();
        app.objectManager.selectObject(object.id);

        // 切换回选择工具
        this.setTool('select');

        return object;
    }

    // 位置检测
    getObjectAtPosition(x, y) {
        const objects = app.objectManager.getAllObjects();
        // 从后往前检查（后添加的在上层）
        for (let i = objects.length - 1; i >= 0; i--) {
            const object = objects[i];
            if (object.visible && object.hitTest(x, y)) {
                return object;
            }
        }
        return null;
    }

    // 选择框
    startSelectionBox(pos) {
        this.selectionBox = {
            startX: pos.x,
            startY: pos.y,
            element: Utils.createElement('div', 'selection-box')
        };
        this.container.appendChild(this.selectionBox.element);
    }

    updateSelectionBox(pos) {
        if (!this.selectionBox) return;

        const startX = this.selectionBox.startX;
        const startY = this.selectionBox.startY;
        const x = Math.min(startX, pos.x);
        const y = Math.min(startY, pos.y);
        const width = Math.abs(pos.x - startX);
        const height = Math.abs(pos.y - startY);

        this.selectionBox.element.style.cssText = `
            left: ${x}px;
            top: ${y}px;
            width: ${width}px;
            height: ${height}px;
        `;
    }

    endSelectionBox() {
        if (!this.selectionBox) return;

        const rect = this.selectionBox.element.getBoundingClientRect();
        const canvasRect = this.canvas.getBoundingClientRect();

        const selectionRect = {
            x: rect.left - canvasRect.left,
            y: rect.top - canvasRect.top,
            width: rect.width,
            height: rect.height
        };

        // 选择框内的对象
        const objects = app.objectManager.getAllObjects();
        objects.forEach(object => {
            const bounds = object.getBounds();
            if (this.rectsIntersect(selectionRect, bounds)) {
                app.objectManager.selectObject(object.id, true);
            }
        });

        this.container.removeChild(this.selectionBox.element);
        this.selectionBox = null;
    }

    rectsIntersect(rect1, rect2) {
        return !(rect1.x > rect2.x + rect2.width ||
            rect1.x + rect1.width < rect2.x ||
            rect1.y > rect2.y + rect2.height ||
            rect1.y + rect1.height < rect2.y);
    }

    // 全选
    selectAll() {
        app.objectManager.clearSelection();
        app.objectManager.getAllObjects().forEach(object => {
            app.objectManager.selectObject(object.id, true);
        });
    }

    // 网格吸附
    snapToGrid(x, y) {
        if (this.grid.visible) {
            const gridSize = this.grid.size;
            return {
                x: Math.round(x / gridSize) * gridSize,
                y: Math.round(y / gridSize) * gridSize
            };
        }
        return { x, y };
    }

    // 缩放
    setZoom(zoom) {
        this.zoom = Utils.clamp(zoom, 0.1, 5.0);
        this.canvas.style.transform = `scale(${this.zoom})`;
        this.canvas.style.transformOrigin = 'center center';

        // 更新缩放显示
        const zoomDisplay = document.getElementById('zoomLevel');
        if (zoomDisplay) {
            zoomDisplay.textContent = Math.round(this.zoom * 100) + '%';
        }

        this.emit('zoomChanged', this.zoom);
    }

    zoomIn() {
        this.setZoom(this.zoom * 1.2);
    }

    zoomOut() {
        this.setZoom(this.zoom / 1.2);
    }

    fitToScreen() {
        this.setZoom(1.0);
    }

    // 场景属性
    setSize(width, height) {
        this.width = width;
        this.height = height;
        this.updateCanvasSize();
        this.emit('sizeChanged', { width, height });
    }

    updateCanvasSize() {
        if (this.canvas) {
            this.canvas.style.width = this.width + 'px';
            this.canvas.style.height = this.height + 'px';
        }
        this.updateGrid();
    }

    setName(name) {
        this.name = name;
        this.emit('nameChanged', name);
    }

    setBackgroundColor(color) {
        this.backgroundColor = color;
        if (this.canvas) {
            this.canvas.style.backgroundColor = color;
        }
        this.emit('backgroundColorChanged', color);
    }

    // 网格
    updateGrid() {
        const gridOverlay = document.getElementById('gridOverlay');
        if (gridOverlay && this.grid.visible) {
            gridOverlay.style.backgroundSize = `${this.grid.size}px ${this.grid.size}px`;
            gridOverlay.style.opacity = '0.1';
        }
    }

    setGridVisible(visible) {
        this.grid.visible = visible;
        const gridOverlay = document.getElementById('gridOverlay');
        if (gridOverlay) {
            gridOverlay.style.display = visible ? 'block' : 'none';
        }
    }

    // 上下文菜单
    showContextMenu(x, y, object) {
        // 移除已存在的菜单
        const existingMenu = document.querySelector('.context-menu');
        if (existingMenu) {
            existingMenu.remove();
        }

        const menu = Utils.createElement('div', 'context-menu');
        menu.style.left = x + 'px';
        menu.style.top = y + 'px';

        const items = [];

        if (object) {
            items.push(
                { text: '复制', action: () => app.objectManager.copy() },
                { text: '删除', action: () => app.objectManager.removeObject(object.id) },
                { separator: true },
                { text: '置于顶层', action: () => app.objectManager.bringToFront(object.id) },
                { text: '置于底层', action: () => app.objectManager.sendToBack(object.id) }
            );
        } else {
            items.push(
                { text: '粘贴', action: () => app.objectManager.paste(), disabled: app.objectManager.clipboard.length === 0 },
                { separator: true },
                { text: '全选', action: () => this.selectAll() }
            );
        }

        items.forEach(item => {
            if (item.separator) {
                menu.appendChild(Utils.createElement('div', 'context-menu-separator'));
            } else {
                const menuItem = Utils.createElement('div',
                    `context-menu-item ${item.disabled ? 'disabled' : ''}`);
                menuItem.textContent = item.text;
                if (!item.disabled) {
                    menuItem.addEventListener('click', () => {
                        item.action();
                        menu.remove();
                    });
                }
                menu.appendChild(menuItem);
            }
        });

        document.body.appendChild(menu);

        // 点击外部关闭菜单
        const closeMenu = (e) => {
            if (!menu.contains(e.target)) {
                menu.remove();
                document.removeEventListener('click', closeMenu);
            }
        };
        setTimeout(() => document.addEventListener('click', closeMenu), 0);
    }

    // 序列化
    toJSON() {
        return {
            width: this.width,
            height: this.height,
            name: this.name,
            frameRate: this.frameRate,
            duration: this.duration,
            backgroundColor: this.backgroundColor,
            grid: this.grid,
            objects: app.objectManager.toJSON()
        };
    }

    // 从JSON恢复
    fromJSON(data) {
        this.width = data.width || 800;
        this.height = data.height || 600;
        this.name = data.name || 'my_animation';
        this.frameRate = data.frameRate || 30;
        this.duration = data.duration || 5;
        this.backgroundColor = data.backgroundColor || '#ffffff';
        this.grid = { ...this.grid, ...data.grid };

        this.updateCanvasSize();
        this.setBackgroundColor(this.backgroundColor);

        if (data.objects) {
            app.objectManager.fromJSON(data.objects);
            // 将对象添加到DOM
            app.objectManager.getAllObjects().forEach(object => {
                this.container.appendChild(object.element);
            });
        }
    }

    // 清空场景
    clear() {
        app.objectManager.clear();
        this.setZoom(1.0);
        this.setTool('select');
    }

    // 缩放控制
    zoomIn() {
        this.setZoom(Math.min(this.zoom * 1.2, 5.0));
    }

    zoomOut() {
        this.setZoom(Math.max(this.zoom / 1.2, 0.1));
    }

    fitToScreen() {
        const viewport = this.viewport.getBoundingClientRect();
        const canvasWidth = this.width;
        const canvasHeight = this.height;
        const scaleX = (viewport.width - 40) / canvasWidth;
        const scaleY = (viewport.height - 40) / canvasHeight;
        this.setZoom(Math.min(scaleX, scaleY, 1.0));
    }

    setZoom(zoom) {
        this.zoom = Utils.clamp(zoom, 0.1, 5.0);

        if (this.canvas) {
            this.canvas.style.transform = `scale(${this.zoom})`;
        }

        // 更新缩放显示
        const zoomLevel = document.getElementById('zoomLevel');
        if (zoomLevel) {
            zoomLevel.textContent = Math.round(this.zoom * 100) + '%';
        }

        this.emit('zoomChanged', this.zoom);
    }

    // 拖放事件处理
    onDragOver(e) {
        e.preventDefault();
        e.dataTransfer.dropEffect = 'copy';
    }

    onDrop(e) {
        e.preventDefault();
        const objectType = e.dataTransfer.getData('text/plain');

        if (objectType) {
            const pos = this.screenToCanvasCoords(e.clientX, e.clientY);

            try {
                const object = ObjectFactory.createObject(objectType);
                object.setPosition(pos.x - object.width / 2, pos.y - object.height / 2);

                app.objectManager.addObject(object);
                app.objectManager.clearSelection();
                app.objectManager.selectObject(object.id);

                console.log(`Dropped ${objectType} object at:`, pos);
            } catch (error) {
                console.error(`Failed to create dropped object:`, error);
            }
        }
    }
}

// 导出
window.Scene = Scene;
