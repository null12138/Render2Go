// 对象管理系统 - objects.js

// 基础场景对象类
class SceneObject {
    constructor(type, id = null) {
        this.id = id || Utils.generateId();
        this.type = type;
        this.name = `${type}_${this.id.split('_')[1]}`;
        this.x = 0;
        this.y = 0;
        this.width = 100;
        this.height = 100;
        this.rotation = 0;
        this.scaleX = 1;
        this.scaleY = 1;
        this.opacity = 1;
        this.color = '#3b82f6';
        this.strokeColor = '#1e293b';
        this.strokeWidth = 2;
        this.visible = true;
        this.locked = false;
        this.selected = false;

        // 动画属性
        this.animations = [];
        this.currentTime = 0;

        // DOM元素
        this.element = null;
        this.createElement();
    }

    createElement() {
        this.element = Utils.createElement('div', `scene-object object-${this.type}`);
        this.element.dataset.objectId = this.id;

        // 添加事件监听器
        this.element.addEventListener('mousedown', (e) => {
            // 不阻止事件传播，让Scene类也能处理
            if (window.app && window.app.objectManager) {
                if (!e.ctrlKey && !e.shiftKey) {
                    window.app.objectManager.clearSelection();
                }
                window.app.objectManager.selectObject(this.id, e.ctrlKey || e.shiftKey);
            }
        });

        this.element.addEventListener('contextmenu', (e) => {
            e.preventDefault();
            e.stopPropagation();
            if (window.app && window.app.scene) {
                window.app.scene.showContextMenu(e.clientX, e.clientY, this);
            }
        });

        this.updateElement();
    }

    updateElement() {
        if (!this.element) return;

        const transform = `translate(${this.x}px, ${this.y}px) 
                          rotate(${this.rotation}deg) 
                          scale(${this.scaleX}, ${this.scaleY})`;

        this.element.style.cssText = `
            position: absolute;
            transform: ${transform};
            opacity: ${this.opacity};
            visibility: ${this.visible ? 'visible' : 'hidden'};
            z-index: ${this.locked ? 1 : 10};
        `;

        if (this.selected) {
            this.element.classList.add('selected');
        } else {
            this.element.classList.remove('selected');
        }

        this.updateTypeSpecificStyles();

        // 更新变形控制点位置
        if (this.resizeHandles) {
            this.updateResizeHandles();
        }
    }

    updateTypeSpecificStyles() {
        // 子类重写此方法
    }

    // 属性设置
    setPosition(x, y) {
        this.x = x;
        this.y = y;
        this.updateElement();
    }

    setSize(width, height) {
        this.width = width;
        this.height = height;
        this.updateElement();
    }

    setRotation(rotation) {
        this.rotation = rotation;
        this.updateElement();
    }

    setScale(scaleX, scaleY = scaleX) {
        this.scaleX = scaleX;
        this.scaleY = scaleY;
        this.updateElement();
    }

    setOpacity(opacity) {
        this.opacity = Utils.clamp(opacity, 0, 1);
        this.updateElement();
    }

    setColor(color) {
        this.color = color;
        this.updateElement();
    }

    setVisible(visible) {
        this.visible = visible;
        this.updateElement();
    }

    setLocked(locked) {
        this.locked = locked;
        this.updateElement();
    }

    setSelected(selected) {
        this.selected = selected;
        this.updateElement();

        // 添加或移除变形控制点
        if (selected) {
            this.addResizeHandles();
        } else {
            this.removeResizeHandles();
        }
    }

    // 添加变形控制点
    addResizeHandles() {
        if (this.resizeHandles) return; // 已存在则不重复添加

        this.resizeHandles = [];
        const handlePositions = [
            { pos: 'nw', cursor: 'nw-resize' },
            { pos: 'n', cursor: 'n-resize' },
            { pos: 'ne', cursor: 'ne-resize' },
            { pos: 'e', cursor: 'e-resize' },
            { pos: 'se', cursor: 'se-resize' },
            { pos: 's', cursor: 's-resize' },
            { pos: 'sw', cursor: 'sw-resize' },
            { pos: 'w', cursor: 'w-resize' }
        ];

        handlePositions.forEach(({ pos, cursor }) => {
            const handle = Utils.createElement('div', 'resize-handle');
            handle.className = `resize-handle resize-handle-${pos}`;
            handle.style.cssText = `
                position: absolute;
                width: 8px;
                height: 8px;
                background: #3b82f6;
                border: 1px solid #ffffff;
                border-radius: 50%;
                cursor: ${cursor};
                z-index: 1001;
                user-select: none;
                pointer-events: auto;
            `;

            // 绑定拖拽事件
            handle.addEventListener('mousedown', (e) => this.onResizeStart(e, pos));

            this.element.appendChild(handle);
            this.resizeHandles.push({ element: handle, position: pos });
        });

        this.updateResizeHandles();
    }

    // 移除变形控制点
    removeResizeHandles() {
        if (this.resizeHandles) {
            this.resizeHandles.forEach(handle => {
                if (handle.element.parentNode) {
                    handle.element.parentNode.removeChild(handle.element);
                }
            });
            this.resizeHandles = null;
        }
    }

    // 更新控制点位置
    updateResizeHandles() {
        if (!this.resizeHandles) return;

        this.resizeHandles.forEach(handle => {
            const pos = handle.position;
            let x, y;

            switch (pos) {
                case 'nw': x = -4; y = -4; break;
                case 'n': x = this.width / 2 - 4; y = -4; break;
                case 'ne': x = this.width - 4; y = -4; break;
                case 'e': x = this.width - 4; y = this.height / 2 - 4; break;
                case 'se': x = this.width - 4; y = this.height - 4; break;
                case 's': x = this.width / 2 - 4; y = this.height - 4; break;
                case 'sw': x = -4; y = this.height - 4; break;
                case 'w': x = -4; y = this.height / 2 - 4; break;
            }

            handle.element.style.left = x + 'px';
            handle.element.style.top = y + 'px';
        });
    }

    // 开始变形拖拽
    onResizeStart(e, position) {
        e.preventDefault();
        e.stopPropagation();

        this.resizing = true;
        this.resizePosition = position;
        this.resizeStartBounds = {
            x: this.x,
            y: this.y,
            width: this.width,
            height: this.height
        };

        const rect = app.scene.canvas.getBoundingClientRect();
        this.resizeStartMouse = {
            x: e.clientX - rect.left,
            y: e.clientY - rect.top
        };

        // 绑定全局事件
        this.onResizeMove = this.onResizeMove.bind(this);
        this.onResizeEnd = this.onResizeEnd.bind(this);
        document.addEventListener('mousemove', this.onResizeMove);
        document.addEventListener('mouseup', this.onResizeEnd);

        console.log('Started resizing:', position);
    }

    // 变形拖拽中
    onResizeMove(e) {
        if (!this.resizing) return;

        const rect = app.scene.canvas.getBoundingClientRect();
        const currentMouse = {
            x: e.clientX - rect.left,
            y: e.clientY - rect.top
        };

        const dx = currentMouse.x - this.resizeStartMouse.x;
        const dy = currentMouse.y - this.resizeStartMouse.y;

        const { x, y, width, height } = this.resizeStartBounds;
        let newX = x, newY = y, newWidth = width, newHeight = height;

        // 根据拖拽位置计算新的尺寸和位置
        switch (this.resizePosition) {
            case 'nw':
                newX = x + dx;
                newY = y + dy;
                newWidth = width - dx;
                newHeight = height - dy;
                break;
            case 'n':
                newY = y + dy;
                newHeight = height - dy;
                break;
            case 'ne':
                newY = y + dy;
                newWidth = width + dx;
                newHeight = height - dy;
                break;
            case 'e':
                newWidth = width + dx;
                break;
            case 'se':
                newWidth = width + dx;
                newHeight = height + dy;
                break;
            case 's':
                newHeight = height + dy;
                break;
            case 'sw':
                newX = x + dx;
                newWidth = width - dx;
                newHeight = height + dy;
                break;
            case 'w':
                newX = x + dx;
                newWidth = width - dx;
                break;
        }

        // 限制最小尺寸
        const minSize = 10;
        if (newWidth < minSize) {
            if (this.resizePosition.includes('w')) {
                newX = x + width - minSize;
            }
            newWidth = minSize;
        }
        if (newHeight < minSize) {
            if (this.resizePosition.includes('n')) {
                newY = y + height - minSize;
            }
            newHeight = minSize;
        }

        // 应用新的尺寸和位置
        this.x = newX;
        this.y = newY;
        this.width = newWidth;
        this.height = newHeight;

        // 特殊处理圆形对象
        if (this.type === 'circle') {
            this.radius = Math.min(newWidth, newHeight) / 2;
            this.width = this.radius * 2;
            this.height = this.radius * 2;
        }

        this.updateElement();
        this.updateResizeHandles();
    }

    // 结束变形拖拽
    onResizeEnd(e) {
        if (!this.resizing) return;

        this.resizing = false;
        this.resizePosition = null;
        this.resizeStartBounds = null;
        this.resizeStartMouse = null;

        // 移除全局事件
        document.removeEventListener('mousemove', this.onResizeMove);
        document.removeEventListener('mouseup', this.onResizeEnd);

        console.log('Finished resizing');

        // 通知属性面板更新
        if (window.app && window.app.propertiesPanel) {
            window.app.propertiesPanel.updateProperties();
        }
    }

    // 边界检测
    getBounds() {
        return {
            x: this.x,
            y: this.y,
            width: this.width * this.scaleX,
            height: this.height * this.scaleY
        };
    }

    hitTest(x, y) {
        const bounds = this.getBounds();
        return Utils.pointInRect(x, y, bounds.x, bounds.y, bounds.width, bounds.height);
    }

    // 序列化
    toJSON() {
        return {
            id: this.id,
            type: this.type,
            name: this.name,
            x: this.x,
            y: this.y,
            width: this.width,
            height: this.height,
            rotation: this.rotation,
            scaleX: this.scaleX,
            scaleY: this.scaleY,
            opacity: this.opacity,
            color: this.color,
            strokeColor: this.strokeColor,
            strokeWidth: this.strokeWidth,
            visible: this.visible,
            locked: this.locked,
            animations: this.animations
        };
    }

    // 从JSON恢复
    fromJSON(data) {
        Object.assign(this, data);
        this.updateElement();
    }

    // 克隆对象
    clone() {
        const cloned = new this.constructor();
        cloned.fromJSON(this.toJSON());
        cloned.id = Utils.generateId();
        cloned.name = `${this.type}_${cloned.id.split('_')[1]}`;
        cloned.createElement();
        return cloned;
    }

    // 销毁对象
    destroy() {
        // 清理变形控制点
        this.removeResizeHandles();

        // 移除DOM元素
        if (this.element && this.element.parentNode) {
            this.element.parentNode.removeChild(this.element);
        }
    }
}

// 圆形对象
class CircleObject extends SceneObject {
    constructor() {
        super('circle');
        this.radius = 50;
        this.width = this.radius * 2;
        this.height = this.radius * 2;
    }

    updateTypeSpecificStyles() {
        this.element.style.cssText += `
            width: ${this.width}px;
            height: ${this.height}px;
            background: ${this.color};
            border: ${this.strokeWidth}px solid ${this.strokeColor};
            border-radius: 50%;
            box-shadow: ${this.selected ? '0 0 0 2px #3b82f6' : 'none'};
        `;
    }

    setRadius(radius) {
        this.radius = radius;
        this.width = radius * 2;
        this.height = radius * 2;
        this.updateElement();
    }

    hitTest(x, y) {
        const centerX = this.x + this.radius;
        const centerY = this.y + this.radius;
        return Utils.pointInCircle(x, y, centerX, centerY, this.radius);
    }
}

// 矩形对象
class RectangleObject extends SceneObject {
    constructor() {
        super('rectangle');
        this.width = 100;
        this.height = 60;
    }

    updateTypeSpecificStyles() {
        this.element.style.cssText += `
            width: ${this.width}px;
            height: ${this.height}px;
            background: ${this.color};
            border: ${this.strokeWidth}px solid ${this.strokeColor};
            box-shadow: ${this.selected ? '0 0 0 2px #3b82f6' : 'none'};
        `;
    }
}

// 文本对象
class TextObject extends SceneObject {
    constructor() {
        super('text');
        this.text = 'Hello World';
        this.fontSize = 24;
        this.fontFamily = 'Arial, sans-serif';
        this.fontWeight = 'normal';
        this.textAlign = 'center';
        this.width = 200;
        this.height = 30;
    }

    updateTypeSpecificStyles() {
        this.element.style.cssText += `
            width: ${this.width}px;
            height: ${this.height}px;
            color: ${this.color};
            font-size: ${this.fontSize}px;
            font-family: ${this.fontFamily};
            font-weight: ${this.fontWeight};
            text-align: ${this.textAlign};
            line-height: ${this.height}px;
            border: ${this.selected ? '2px dashed #3b82f6' : '1px solid transparent'};
            background: ${this.selected ? 'rgba(59, 130, 246, 0.1)' : 'rgba(255, 255, 255, 0.1)'};
            user-select: none;
            display: flex;
            align-items: center;
            justify-content: center;
        `;
        this.element.textContent = this.text;
    }

    setText(text) {
        this.text = text;
        this.updateElement();
    }

    setFontSize(fontSize) {
        this.fontSize = fontSize;
        this.updateElement();
    }

    setFontFamily(fontFamily) {
        this.fontFamily = fontFamily;
        this.updateElement();
    }

    hitTest(x, y) {
        const hit = Utils.pointInRect(x, y, this.x, this.y, this.width, this.height);
        console.log(`TextObject hitTest: point(${x},${y}) in rect(${this.x},${this.y},${this.width},${this.height}) = ${hit}`);
        return hit;
    }
}

// 线条对象
class LineObject extends SceneObject {
    constructor() {
        super('line');
        this.startX = 0;
        this.startY = 0;
        this.endX = 100;
        this.endY = 100;
        this.updateFromPoints();
    }

    updateFromPoints() {
        const dx = this.endX - this.startX;
        const dy = this.endY - this.startY;
        this.width = Math.abs(dx);
        this.height = Math.abs(dy);
        this.x = Math.min(this.startX, this.endX);
        this.y = Math.min(this.startY, this.endY);

        // 计算旋转角度
        this.rotation = Math.atan2(dy, dx) * 180 / Math.PI;
    }

    updateTypeSpecificStyles() {
        const length = Utils.distance(this.startX, this.startY, this.endX, this.endY);
        this.element.style.cssText += `
            width: ${length}px;
            height: ${this.strokeWidth}px;
            background: ${this.strokeColor};
            transform-origin: 0 50%;
        `;
    }

    setPoints(startX, startY, endX, endY) {
        this.startX = startX;
        this.startY = startY;
        this.endX = endX;
        this.endY = endY;
        this.updateFromPoints();
        this.updateElement();
    }

    hitTest(x, y) {
        // 对于线条，使用更宽松的点击检测
        const tolerance = Math.max(this.strokeWidth, 5);
        return Utils.pointInRect(x, y, this.x - tolerance, this.y - tolerance,
            this.width + tolerance * 2, this.height + tolerance * 2);
    }
}

// 三角形对象
class TriangleObject extends SceneObject {
    constructor() {
        super('triangle');
        this.width = 100;
        this.height = 87;
    }

    updateTypeSpecificStyles() {
        this.element.style.cssText += `
            width: 0;
            height: 0;
            border-left: ${this.width / 2}px solid transparent;
            border-right: ${this.width / 2}px solid transparent;
            border-bottom: ${this.height}px solid ${this.color};
            background: transparent;
        `;
    }

    hitTest(x, y) {
        return Utils.pointInRect(x, y, this.x, this.y, this.width, this.height);
    }
}

// 对象工厂
class ObjectFactory {
    static createObject(type) {
        switch (type) {
            case 'circle':
                return new CircleObject();
            case 'rectangle':
                return new RectangleObject();
            case 'text':
                return new TextObject();
            case 'line':
                return new LineObject();
            case 'triangle':
                return new TriangleObject();
            default:
                throw new Error(`Unknown object type: ${type}`);
        }
    }

    static getObjectTypes() {
        return [
            { type: 'circle', name: '圆形', icon: 'fas fa-circle' },
            { type: 'rectangle', name: '矩形', icon: 'fas fa-square' },
            { type: 'triangle', name: '三角形', icon: 'fas fa-play' },
            { type: 'line', name: '线条', icon: 'fas fa-minus' },
            { type: 'text', name: '文本', icon: 'fas fa-font' }
        ];
    }
}

// 对象管理器
class ObjectManager extends EventEmitter {
    constructor() {
        super();
        this.objects = new Map();
        this.selectedObjects = new Set();
        this.clipboard = [];
    }

    addObject(object) {
        this.objects.set(object.id, object);

        // 将对象元素添加到DOM
        const container = document.getElementById('objectsContainer');
        if (container && object.element) {
            container.appendChild(object.element);
        }

        this.emit('objectAdded', object);
        return object;
    }

    removeObject(id) {
        const object = this.objects.get(id);
        if (object) {
            // 从DOM中移除元素
            if (object.element && object.element.parentNode) {
                object.element.parentNode.removeChild(object.element);
            }

            object.destroy();
            this.objects.delete(id);
            this.selectedObjects.delete(id);
            this.emit('objectRemoved', object);
        }
    }

    getObject(id) {
        return this.objects.get(id);
    }

    getAllObjects() {
        return Array.from(this.objects.values());
    }

    selectObject(id, addToSelection = false) {
        if (!addToSelection) {
            this.clearSelection();
        }

        const object = this.objects.get(id);
        if (object && !object.locked) {
            this.selectedObjects.add(id);
            object.setSelected(true);
            this.emit('selectionChanged', Array.from(this.selectedObjects));
        }
    }

    deselectObject(id) {
        const object = this.objects.get(id);
        if (object) {
            this.selectedObjects.delete(id);
            object.setSelected(false);
            this.emit('selectionChanged', Array.from(this.selectedObjects));
        }
    }

    clearSelection() {
        this.selectedObjects.forEach(id => {
            const object = this.objects.get(id);
            if (object) {
                object.setSelected(false);
            }
        });
        this.selectedObjects.clear();
        this.emit('selectionChanged', []);
    }

    getSelectedObjects() {
        return Array.from(this.selectedObjects).map(id => this.objects.get(id));
    }

    selectAll() {
        this.clearSelection();
        this.objects.forEach((object, id) => {
            if (!object.locked) {
                this.selectedObjects.add(id);
                object.setSelected(true);
            }
        });
        this.emit('selectionChanged', Array.from(this.selectedObjects));
    }

    // 复制粘贴
    copy() {
        this.clipboard = this.getSelectedObjects().map(obj => obj.clone());
        Utils.showNotification(`已复制 ${this.clipboard.length} 个对象`, 'success');
    }

    paste() {
        if (this.clipboard.length > 0) {
            this.clearSelection();
            this.clipboard.forEach(obj => {
                const cloned = obj.clone();
                cloned.setPosition(cloned.x + 20, cloned.y + 20);
                this.addObject(cloned);
                this.selectObject(cloned.id, true);
            });
            Utils.showNotification(`已粘贴 ${this.clipboard.length} 个对象`, 'success');
        }
    }

    // 删除选中对象
    deleteSelected() {
        const selectedIds = Array.from(this.selectedObjects);
        selectedIds.forEach(id => this.removeObject(id));
        if (selectedIds.length > 0) {
            Utils.showNotification(`已删除 ${selectedIds.length} 个对象`, 'success');
        }
    }

    // 图层操作
    bringToFront(id) {
        const object = this.objects.get(id);
        if (object && object.element) {
            object.element.style.zIndex = '1000';
            this.emit('layerChanged', object);
        }
    }

    sendToBack(id) {
        const object = this.objects.get(id);
        if (object && object.element) {
            object.element.style.zIndex = '1';
            this.emit('layerChanged', object);
        }
    }

    // 对齐功能
    alignLeft() {
        const objects = this.getSelectedObjects();
        if (objects.length > 1) {
            const minX = Math.min(...objects.map(obj => obj.x));
            objects.forEach(obj => obj.setPosition(minX, obj.y));
        }
    }

    alignCenter() {
        const objects = this.getSelectedObjects();
        if (objects.length > 1) {
            const avgX = objects.reduce((sum, obj) => sum + obj.x + obj.width / 2, 0) / objects.length;
            objects.forEach(obj => obj.setPosition(avgX - obj.width / 2, obj.y));
        }
    }

    alignRight() {
        const objects = this.getSelectedObjects();
        if (objects.length > 1) {
            const maxX = Math.max(...objects.map(obj => obj.x + obj.width));
            objects.forEach(obj => obj.setPosition(maxX - obj.width, obj.y));
        }
    }

    // 分布功能
    distributeHorizontally() {
        const objects = this.getSelectedObjects().sort((a, b) => a.x - b.x);
        if (objects.length > 2) {
            const first = objects[0];
            const last = objects[objects.length - 1];
            const totalWidth = last.x - first.x;
            const spacing = totalWidth / (objects.length - 1);

            objects.forEach((obj, index) => {
                if (index > 0 && index < objects.length - 1) {
                    obj.setPosition(first.x + spacing * index, obj.y);
                }
            });
        }
    }

    // 清空所有对象
    clear() {
        this.objects.forEach(object => object.destroy());
        this.objects.clear();
        this.selectedObjects.clear();
        this.emit('cleared');
    }

    // 序列化
    toJSON() {
        return {
            objects: Array.from(this.objects.values()).map(obj => obj.toJSON())
        };
    }

    // 从JSON恢复
    fromJSON(data) {
        this.clear();
        if (data.objects) {
            data.objects.forEach(objData => {
                const object = ObjectFactory.createObject(objData.type);
                object.fromJSON(objData);
                this.addObject(object);
            });
        }
    }
}

// 导出对象
window.SceneObject = SceneObject;
window.CircleObject = CircleObject;
window.RectangleObject = RectangleObject;
window.TextObject = TextObject;
window.LineObject = LineObject;
window.TriangleObject = TriangleObject;
window.ObjectFactory = ObjectFactory;
window.ObjectManager = ObjectManager;
